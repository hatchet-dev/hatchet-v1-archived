package notifier

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hatchet-dev/hatchet/internal/config/worker"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/notifier"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"go.temporal.io/sdk/workflow"
)

type NotifierInput struct {
}

type Notifier struct {
	Config *worker.BackgroundConfig
}

func NewNotifier(config *worker.BackgroundConfig) *Notifier {
	return &Notifier{config}
}

func (n *Notifier) NotifyWorkflow(ctx workflow.Context, input NotifierInput) (string, error) {
	repo := n.Config.DB.Repository

	// iterate through all notifs
	notifs, _, err := repo.Notification().ListNotifications(&repository.ListNotificationOpts{
		Resolved: boolPointer(false),
	}, repository.WithLimit(100))

	if err != nil {
		return "", err
	}

	// TODO: paginate
	for _, notif := range notifs {
		// case on last alerted time
		shouldAlert := notif.LastNotified == nil || elapsedHours(notif.LastNotified) > 24

		// call notifier
		if shouldAlert {
			// get a list of users to alert
			members, _, err := repo.Team().ListTeamMembersByTeamID(notif.TeamID, false)

			if err != nil {
				err = multierror.Append(contextualize(err, notif))
				continue
			}

			users := make([]*models.User, 0)

			for _, member := range members {
				users = append(users, &member.OrgMember.User)
			}

			err = n.Config.IncidentNotifier.SendIncidentNotification(&notifier.SendIncidentNotificationOpts{
				Users:        users,
				URL:          fmt.Sprintf("%s/teams/%s/notifications", n.Config.ServerURL, notif.TeamID),
				ModuleName:   notif.Module.Name,
				Title:        notif.Title,
				Message:      notif.Message,
				Notification: notif,
			})

			if err != nil {
				err = multierror.Append(contextualize(err, notif))
				continue
			}

			// update notif with last alerted time
			now := time.Now()
			notif.LastNotified = &now
			notif, err = repo.Notification().UpdateNotification(notif)

			if err != nil {
				err = multierror.Append(contextualize(err, notif))
				continue
			}
		}
	}

	if err != nil {
		return "", err
	}

	return "success", nil
}

func boolPointer(b bool) *bool {
	return &b
}

func elapsedHours(t *time.Time) uint {
	if t == nil {
		return 0
	}

	elapsedTime := time.Now().Sub(*t)
	elapsedHours := elapsedTime.Truncate(time.Hour).Hours()

	return uint(elapsedHours)
}

func contextualize(err error, notif *models.Notification) error {
	if notif == nil {
		return fmt.Errorf("notifier error. Error is: %s", err.Error())
	} else if notif.Module.ID == "" {
		return fmt.Errorf("notifier error: team: [%s], notification:[%s]. Error is: %s", notif.TeamID, notif.NotificationID, err.Error())
	}
	return fmt.Errorf("notifier error: team: [%s], module:[%s %s] notification:[%s]. Error is: %s", notif.TeamID, notif.Module.ID, notif.Module.Name, notif.NotificationID, err.Error())
}
