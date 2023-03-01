package notifications

import (
	"errors"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

func CreateNotificationFromMonitorResult(config *server.Config, teamID string, result *models.ModuleMonitorResult) error {
	if result.Status != models.MonitorResultStatusFailed || result.Severity == models.MonitorResultSeverityLow {
		return nil
	}

	repo := config.DB.Repository

	// get the notification inbox
	inbox, err := repo.Notification().ReadNotificationInboxByTeamID(teamID)

	if err != nil {
		return err
	}

	notificationID := getNotificationIDFromMonitor(result)

	notif, err := repo.Notification().ReadNotificationByNotificationID(teamID, notificationID, &repository.ReadNotificationOpts{
		AutoResolved: repository.BoolPointer(false),
	})

	if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
		return err
	} else if errors.Is(err, repository.RepositoryErrorNotFound) {
		notif := &models.Notification{
			TeamID:              teamID,
			NotificationInboxID: inbox.ID,
			NotificationID:      notificationID,
			Title:               result.Title,
			Message:             result.Message,
			ModuleID:            result.ModuleID,
		}

		notif, err = repo.Notification().CreateNotification(notif)

		if err != nil {
			return err
		}
	}

	notif, err = repo.Notification().AppendModuleRunMonitorResult(notif, result)

	return err
}

func getNotificationIDFromMonitor(result *models.ModuleMonitorResult) string {
	return ToSnake(fmt.Sprintf("monitor-%s-%s", result.ModuleMonitorID, result.Title))
}
