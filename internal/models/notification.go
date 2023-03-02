package models

import (
	"time"

	"github.com/hatchet-dev/hatchet/api/v1/types"
)

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

	ModuleID string
	Module   Module `gorm:"foreignKey:ModuleID"`
}

func (n *Notification) ToAPITypeMeta() *types.NotificationMeta {
	return &types.NotificationMeta{
		APIResourceMeta: n.ToAPITypeMetadata(),
		TeamID:          n.TeamID,
		NotificationID:  n.NotificationID,
		Title:           n.Title,
		Message:         n.Message,
		LastNotified:    n.LastNotified,
		Resolved:        n.Resolved,
		ModuleID:        n.ModuleID,
	}
}

func (n *Notification) ToAPIType() *types.Notification {
	res := &types.Notification{
		NotificationMeta: n.ToAPITypeMeta(),
		Module:           *n.Module.ToAPIType(),
	}

	if n.Runs != nil {
		runs := make([]types.ModuleRun, 0)

		for _, run := range n.Runs {
			runs = append(runs, *run.ToAPIType(nil))
		}

		res.Runs = runs
	}

	if n.MonitorResults != nil {
		results := make([]types.ModuleMonitorResult, 0)

		for _, result := range n.MonitorResults {
			results = append(results, *result.ToAPIType())
		}

		res.MonitorResults = results
	}

	return res
}

// type NotificationMeta struct {
// 	*APIResourceMeta

// 	TeamID         string     `json:"team_id"`
// 	NotificationID string     `json:"notification_id"`
// 	Title          string     `json:"title"`
// 	Message        string     `json:"message"`
// 	LastNotified   *time.Time `json:"last_notified"`
// 	Resolved       bool       `json:"resolved"`
// 	ModuleID       string     `json:"module_id"`
// }

// // swagger:model
// type Notification struct {
// 	*NotificationMeta

// 	Runs           []ModuleRun           `json:"runs"`
// 	MonitorResults []ModuleMonitorResult `json:"monitor_results"`
// 	Module         Module                `json:"module"`
// }
