package types

const (
	URLParamMonitorID URLParam = "monitor_id"
)

type ModuleMonitorKind string

const (
	MonitorKindPlan          ModuleMonitorKind = "plan"
	MonitorKindState         ModuleMonitorKind = "state"
	MonitorKindBeforePlan    ModuleMonitorKind = "before_plan"
	MonitorKindAfterPlan     ModuleMonitorKind = "after_plan"
	MonitorKindBeforeApply   ModuleMonitorKind = "before_apply"
	MonitorKindAfterApply    ModuleMonitorKind = "after_apply"
	MonitorKindBeforeDestroy ModuleMonitorKind = "before_destroy"
	MonitorKindAfterDestroy  ModuleMonitorKind = "after_destroy"
)

// swagger:model
type ModuleMonitorMeta struct {
	*APIResourceMeta

	// the name for the monitor
	// example: drift
	Name string `json:"name"`

	// the description for the monitor
	// example: detects drift
	Description string `json:"description"`

	// the kind of monitor
	// example: plan
	Kind ModuleMonitorKind `json:"kind"`

	// the cron schedule for the monitor
	CronSchedule string `json:"cron_schedule"`
}

// swagger:model
type ModuleMonitor struct {
	*ModuleMonitorMeta

	// the policy bytes for the monitor
	PolicyBytes []byte `json:"policy_bytes"`
}

// swagger:model
type CreateMonitorResultRequest struct {
	MonitorID       string
	Status          string   `json:"status" mapstructure:"POLICY_STATUS" form:"omitempty,oneof=succeeded failed"`
	Severity        string   `json:"severity" mapstructure:"POLICY_SEVERITY" form:"omitempty,oneof=critical high low"`
	Title           string   `json:"title" mapstructure:"POLICY_TITLE" form:"required"`
	SuccessMessage  string   `json:"success_message" mapstructure:"POLICY_SUCCESS_MESSAGE"`
	FailureMessages []string `json:"failure_messages" mapstructure:"POLICY_FAILURE_MESSAGES"`
}

// swagger:model
type CreateMonitorRequest struct {
	Name         string `json:"name" form:"required"`
	CronSchedule string `json:"cron_schedule" form:"required"`
	PolicyBytes  []byte `json:"policy_bytes" form:"required"`
}

// swagger:parameters listMonitors
type ListMonitorsRequest struct {
	*PaginationRequest
}

// swagger:model
type ListMonitorsResponse struct {
	Pagination *PaginationResponse  `json:"pagination"`
	Rows       []*ModuleMonitorMeta `json:"rows"`
}
