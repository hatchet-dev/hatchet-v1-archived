package models

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

type ModuleMonitorPresetPolicyName string

const (
	ModuleMonitorPresetPolicyNameDrift ModuleMonitorPresetPolicyName = "drift"
)

type ModuleMonitor struct {
	Base

	Kind ModuleMonitorKind

	PresetPolicyName ModuleMonitorPresetPolicyName
	PolicyBytes      []byte

	MatchModules   []byte
	MatchResources []byte
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

type ModuleMonitorResult struct {
	Base

	ModuleID string
	Module   Module `gorm:"foreignKey:ModuleID"`

	// (optional) The module run id, if this result is attached to a specific module run id (for
	// before_plan, after_plan, etc)
	ModuleRunID string

	ModuleMonitorID string

	Status MonitorResultStatus

	Title    string
	Message  string
	Severity MonitorResultSeverity
}
