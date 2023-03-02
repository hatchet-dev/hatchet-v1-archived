package notifications

import (
	"errors"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

func CreateNotificationFromModuleRun(config *server.Config, teamID string, run *models.ModuleRun) error {
	if run.Kind != models.ModuleRunKindApply || run.Status != models.ModuleRunStatusFailed {
		return nil
	}

	repo := config.DB.Repository

	// get the notification inbox
	inbox, err := repo.Notification().ReadNotificationInboxByTeamID(teamID)

	if err != nil {
		return err
	}

	notificationID := getNotificationIDFromRun(run)

	notif, err := repo.Notification().ReadNotificationByNotificationID(teamID, notificationID, &repository.ReadNotificationOpts{
		AutoResolved: repository.BoolPointer(false),
	})

	if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
		return err
	} else if errors.Is(err, repository.RepositoryErrorNotFound) {
		title, err := getNotificationTitleFromRun(config, teamID, run)

		if err != nil {
			return err
		}

		notif := &models.Notification{
			TeamID:              teamID,
			NotificationInboxID: inbox.ID,
			NotificationID:      notificationID,
			Title:               title,
			Message:             run.StatusDescription,
			ModuleID:            run.ModuleID,
		}

		notif, err = repo.Notification().CreateNotification(notif)

		if err != nil {
			return err
		}
	}

	notif, err = repo.Notification().AppendModuleRun(notif, run)

	return err
}

func getNotificationIDFromRun(run *models.ModuleRun) string {
	return fmt.Sprintf("run-%s-%s", run.ID, string(run.Kind))
}

func getNotificationTitleFromRun(config *server.Config, teamID string, run *models.ModuleRun) (string, error) {
	module, err := config.DB.Repository.Module().ReadModuleByID(teamID, run.ModuleID)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Running %s for module %s failed", string(run.Kind), module.Name), nil
}
