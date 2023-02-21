package types

// swagger:model
type CreateMonitorResultRequest struct {
	MonitorID       string
	Status          string   `json:"status" mapstructure:"POLICY_STATUS" form:"omitempty,oneof=succeeded failed"`
	Severity        string   `json:"severity" mapstructure:"POLICY_SEVERITY" form:"omitempty,oneof=critical high low"`
	Title           string   `json:"title" mapstructure:"POLICY_TITLE" form:"required"`
	SuccessMessage  string   `json:"success_message" mapstructure:"POLICY_SUCCESS_MESSAGE"`
	FailureMessages []string `json:"failure_messages" mapstructure:"POLICY_FAILURE_MESSAGES"`
}
