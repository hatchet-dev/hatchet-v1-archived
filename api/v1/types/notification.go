package types

import "time"

const (
	URLParamNotificationID URLParam = "notification_id"
)

// swagger:model
type NotificationMeta struct {
	*APIResourceMeta

	TeamID         string     `json:"team_id"`
	NotificationID string     `json:"notification_id"`
	Title          string     `json:"title"`
	Message        string     `json:"message"`
	LastNotified   *time.Time `json:"last_notified"`
	Resolved       bool       `json:"resolved"`
	ModuleID       string     `json:"module_id"`
}

// swagger:model
type Notification struct {
	*NotificationMeta

	Runs           []ModuleRun           `json:"runs"`
	MonitorResults []ModuleMonitorResult `json:"monitor_results"`
	Module         Module                `json:"module"`
}

// swagger:parameters listNotifications
type ListNotificationsRequest struct {
	*PaginationRequest

	TeamID string `json:"team_id" schema:"team_id"`
}

// swagger:model
type ListNotificationsResponse struct {
	Pagination *PaginationResponse `json:"pagination"`
	Rows       []*NotificationMeta `json:"rows"`
}

// swagger:model
type GetNotificationResponse Notification
