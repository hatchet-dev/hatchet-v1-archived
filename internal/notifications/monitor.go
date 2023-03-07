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
		title, err := getNotificationTitleFromMonitor(config, teamID, result)

		if err != nil {
			return err
		}

		notif := &models.Notification{
			TeamID:              teamID,
			NotificationInboxID: inbox.ID,
			NotificationID:      notificationID,
			Title:               title,
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
	return fmt.Sprintf("monitor-%s-%s", result.ModuleMonitorID, ToSnake(result.Title))
}

func getNotificationTitleFromMonitor(config *server.Config, teamID string, result *models.ModuleMonitorResult) (string, error) {
	// get monitor name
	monitor, err := config.DB.Repository.ModuleMonitor().ReadModuleMonitorByID(teamID, result.ModuleMonitorID)

	if err != nil {
		return "", err
	}

	// get module name
	module := &result.Module

	if module.ID == "" {
		module, err = config.DB.Repository.Module().ReadModuleByID(teamID, result.ModuleID)

		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("Monitor %s (%s) for module %s failed", monitor.DisplayName, getMonitorDescription(monitor), module.Name), nil
}

func getMonitorDescription(monitor *models.ModuleMonitor) string {
	switch monitor.Kind {
	case models.MonitorKindState:
		return "state check"
	case models.MonitorKindPlan:
		return "plan check"
	case models.MonitorKindBeforePlan:
		return "run before plan"
	case models.MonitorKindAfterPlan:
		return "run after plan"
	case models.MonitorKindBeforeApply:
		return "run before apply"
	case models.MonitorKindAfterApply:
		return "run after apply"
	}

	return ""
}
