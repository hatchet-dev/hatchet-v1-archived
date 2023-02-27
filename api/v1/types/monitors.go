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
	PolicyBytes string `json:"policy_bytes"`
}

type MonitorResultSeverity string

const (
	MonitorResultSeverityCritical MonitorResultSeverity = "critical"
	MonitorResultSeverityHigh     MonitorResultSeverity = "high"
	MonitorResultSeverityLow      MonitorResultSeverity = "low"
)

type MonitorResultStatus string

const (
	MonitorResultStatusSucceeded MonitorResultStatus = "succeeded"
	MonitorResultStatusFailed    MonitorResultStatus = "failed"
)

// swagger:model
type ModuleMonitorResult struct {
	*APIResourceMeta

	ModuleID        string `json:"module_id"`
	ModuleMonitorID string `json:"module_monitor_id"`

	Status MonitorResultStatus `json:"status"`

	Title    string                `json:"title"`
	Message  string                `json:"message"`
	Severity MonitorResultSeverity `json:"severity"`
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
	Name         string            `json:"name" form:"required"`
	Description  string            `json:"description" form:"required"`
	Kind         ModuleMonitorKind `json:"kind" form:"required,oneof=plan state before_plan after_plan before_apply after_apply before_destroy after_destroy"`
	CronSchedule string            `json:"cron_schedule" form:"required,cron"`
	PolicyBytes  string            `json:"policy_bytes" form:"required"`

	// A list of module ids. If empty or omitted, this monitor targets all modules.
	Modules []string `json:"modules,omitempty"`
}

// swagger:model
type CreateMonitorResponse ModuleMonitor

// swagger:model
type GetMonitorResponse ModuleMonitor

// swagger:parameters listMonitors
type ListMonitorsRequest struct {
	*PaginationRequest
}

// swagger:model
type ListMonitorsResponse struct {
	Pagination *PaginationResponse  `json:"pagination"`
	Rows       []*ModuleMonitorMeta `json:"rows"`
}

// swagger:parameters listMonitorResults
type ListMonitorResultsRequest struct {
	*PaginationRequest

	// The monitor id to filter by
	// in: query
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	ModuleMonitorID string `schema:"module_monitor_id" json:"module_monitor_id" form:"omitempty,uuid"`

	// The module id to filter by
	// in: query
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	ModuleID string `schema:"module_id" json:"module_id" form:"omitempty,uuid"`
}

// swagger:model
type ListMonitorResultsResponse struct {
	Pagination *PaginationResponse    `json:"pagination"`
	Rows       []*ModuleMonitorResult `json:"rows"`
}
