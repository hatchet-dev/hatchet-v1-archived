package models

import "time"

type NotificationInbox struct {
	Base

	TeamID string `gorm:"unique"`

	Notifications []Notification
}

type Notification struct {
	Base

	TeamID              string
	NotificationInboxID string
	NotificationID      string
	Title               string
	Message             string

	LastNotified *time.Time

	// Whether this has been resolved - can be manually resolved
	Resolved bool

	// Whether this has been auto-resolved - for example, if a monitor eventually succeeds
	AutoResolved bool

	// A list of runs for failed module operations
	Runs []ModuleRun `gorm:"many2many:notification_to_runs;"`

	// A list of monitor results for monitor failures
	MonitorResults []ModuleMonitorResult `gorm:"many2many:notification_to_monitor_results;"`
}
