package models

import "github.com/hatchet-dev/hatchet/api/v1/types"

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

	TeamID string

	DisplayName  string
	Description  string
	Kind         ModuleMonitorKind
	CronSchedule string

	PresetPolicyName ModuleMonitorPresetPolicyName

	CurrentMonitorPolicyBytesVersionID string
	CurrentMonitorPolicyBytesVersion   MonitorPolicyBytesVersion `gorm:"foreignKey:CurrentMonitorPolicyBytesVersionID"`

	// A list of modules to target. If left empty, targets all modules.
	Modules []Module `gorm:"many2many:monitors_to_modules;"`

	// IsDefault controls whether this is a default monitor for all modules. If this is a default,
	// it cannot be configured from the dashboard.
	IsDefault bool

	// Whether the monitor is disabled
	Disabled bool

	MatchChildModules []byte
	MatchProviders    []byte
	MatchResources    []byte
}

func (m *ModuleMonitor) ToAPITypeMeta() *types.ModuleMonitorMeta {
	return &types.ModuleMonitorMeta{
		APIResourceMeta: m.ToAPITypeMetadata(),
		Name:            m.DisplayName,
		Description:     m.Description,
		Kind:            types.ModuleMonitorKind(m.Kind),
		CronSchedule:    m.CronSchedule,
		IsDefault:       m.IsDefault,
		Disabled:        m.Disabled,
	}
}

func (m *ModuleMonitor) ToAPIType() *types.ModuleMonitor {
	modules := make([]string, 0)

	for _, mod := range m.Modules {
		modules = append(modules, mod.ID)
	}

	return &types.ModuleMonitor{
		ModuleMonitorMeta: m.ToAPITypeMeta(),
		PolicyBytes:       string(m.CurrentMonitorPolicyBytesVersion.PolicyBytes),
		Modules:           modules,
	}
}

func (m *ModuleMonitor) IsCronKind() bool {
	return IsCronKind(m.Kind)
}

func IsCronKind(kind ModuleMonitorKind) bool {
	return kind == MonitorKindPlan || kind == MonitorKindState
}

func (m *ModuleMonitor) ShouldRunForModule(modID string) bool {
	if m.Disabled {
		return false
	}

	if m.Modules == nil || len(m.Modules) == 0 {
		return true
	}

	for _, mod := range m.Modules {
		if mod.ID == modID {
			return true
		}
	}

	return false
}

type MonitorPolicyBytesVersion struct {
	Base

	ModuleMonitorID string
	Version         uint

	PolicyBytes []byte
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

	TeamID string

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

func (m *ModuleMonitorResult) ToAPIType() *types.ModuleMonitorResult {
	return &types.ModuleMonitorResult{
		APIResourceMeta: m.ToAPITypeMetadata(),
		ModuleRunID:     m.ModuleRunID,
		ModuleID:        m.ModuleID,
		ModuleName:      m.Module.Name,
		ModuleMonitorID: m.ModuleMonitorID,
		Status:          types.MonitorResultStatus(m.Status),
		Title:           m.Title,
		Message:         m.Message,
		Severity:        types.MonitorResultSeverity(m.Severity),
	}
}
