package repository

import "github.com/hatchet-dev/hatchet/internal/models"

type ReadNotificationOpts struct {
	AutoResolved *bool
}

type ListNotificationOpts struct {
	AutoResolved *bool
	Resolved     *bool
}

// NotificationRepository represents the set of queries on the NotificationInbox and Notification models
type NotificationRepository interface {
	// --- NotificationInbox queries ---
	//
	// CreateModule creates a new module in the database
	CreateNotificationInbox(inbox *models.NotificationInbox) (*models.NotificationInbox, RepositoryError)

	// ReadModuleByID reads the module by its unique team id
	ReadNotificationInboxByTeamID(teamID string) (*models.NotificationInbox, RepositoryError)

	// UpdateNotificationInbox updates any modified values for the notification inbox
	UpdateNotificationInbox(inbox *models.NotificationInbox) (*models.NotificationInbox, RepositoryError)

	// --- Notification queries ---
	//
	// CreateNotification creates a new notification in the database
	CreateNotification(notif *models.Notification) (*models.Notification, RepositoryError)

	// ReadNotificationByID reads the notification by its uuid
	ReadNotificationByID(teamID, id string) (*models.Notification, RepositoryError)

	// ReadNotificationByNotificationID reads the notification by its notification id
	// NOTE: this is NOT the UUID (use ReadNotificationByID for that)
	ReadNotificationByNotificationID(teamID, notificationID string, opts *ReadNotificationOpts) (*models.Notification, RepositoryError)

	// ListNotifications lists notifications (paginated)
	// NOTE: this should only be called by internal workflows. Any team-scoped handlers should use ListNotificationsByTeamID
	ListNotifications(filterOpts *ListNotificationOpts, opts ...QueryOption) ([]*models.Notification, *PaginatedResult, RepositoryError)

	// ListNotificationsByTeamIDs lists notifications (paginated)
	ListNotificationsByTeamIDs(teamIDs []string, filterOpts *ListNotificationOpts, opts ...QueryOption) ([]*models.Notification, *PaginatedResult, RepositoryError)

	// UpdateNotification updates any modified values for a notification
	UpdateNotification(notif *models.Notification) (*models.Notification, RepositoryError)

	// AppendModuleRun adds a single module run to the notification
	AppendModuleRun(notif *models.Notification, run *models.ModuleRun) (*models.Notification, RepositoryError)

	// AppendModuleRunMonitorResult adds a single monitor result to the notification
	AppendModuleRunMonitorResult(notif *models.Notification, result *models.ModuleMonitorResult) (*models.Notification, RepositoryError)

	// DeleteNotification soft-deletes a notif
	DeleteNotification(notif *models.Notification) (*models.Notification, RepositoryError)
}
