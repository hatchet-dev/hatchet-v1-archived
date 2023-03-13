package types

const (
	URLParamModuleID        URLParam = "module_id"
	URLParamModuleRunID     URLParam = "module_run_id"
	URLParamModuleEnvVarsID URLParam = "module_env_vars_id"
	URLParamModuleValuesID  URLParam = "module_values_id"
)

type ModuleLockKind string

const (
	ModuleLockKindGithubBranch ModuleLockKind = "github_branch"
	ModuleLockKindManual       ModuleLockKind = "manual"
)

// swagger:model
type Module struct {
	*APIResourceMeta

	// the name for the module
	// example: eks
	Name string `json:"name"`

	DeploymentConfig ModuleDeploymentConfig `json:"deployment"`

	LockID   string         `json:"lock_id"`
	LockKind ModuleLockKind `json:"lock_kind"`

	CurrentValuesVersionID  string `json:"current_values_version_id"`
	CurrentEnvVarsVersionID string `json:"current_env_vars_version_id"`
}

type ModuleRunStatus string

const (
	ModuleRunStatusCompleted ModuleRunStatus = "completed"
	ModuleRunStatusFailed    ModuleRunStatus = "failed"
)

type ModuleRunKind string

const (
	ModuleRunKindPlan    ModuleRunKind = "plan"
	ModuleRunKindApply   ModuleRunKind = "apply"
	ModuleRunKindDestroy ModuleRunKind = "destroy"
	ModuleRunKindInit    ModuleRunKind = "init"
)

// swagger:model
type ModuleRunOverview struct {
	*APIResourceMeta

	Status            ModuleRunStatus `json:"status"`
	StatusDescription string          `json:"status_description"`
	Kind              ModuleRunKind   `json:"kind"`
}

// swagger:model
type ModuleRun struct {
	*ModuleRunOverview

	ModuleRunConfig *ModuleRunConfig `json:"config,omitempty"`

	ModuleRunPullRequest *GithubPullRequest `json:"github_pull_request,omitempty"`

	Monitors []ModuleMonitor `json:"monitors,omitempty"`

	MonitorResults []ModuleMonitorResult `json:"monitor_results,omitempty"`
}

type ModuleRunTriggerKind string

const (
	ModuleRunTriggerKindGithub ModuleRunTriggerKind = "github"
	ModuleRunTriggerKindManual ModuleRunTriggerKind = "manual"
)

// swagger:model
type ModuleRunConfig struct {
	TriggerKind     ModuleRunTriggerKind `json:"trigger_kind"`
	GithubCommitSHA string               `json:"github_commit_sha"`
	EnvVarVersionID string               `json:"env_var_version_id"`
	ValuesVersionID string               `json:"values_version_id"`
}

// swagger:model
type GithubPullRequest struct {
	GithubRepositoryOwner       string `json:"github_repository_owner"`
	GithubRepositoryName        string `json:"github_repository_name"`
	GithubPullRequestID         int64  `json:"github_pull_request_id"`
	GithubPullRequestTitle      string `json:"github_pull_request_title"`
	GithubPullRequestNumber     int64  `json:"github_pull_request_number"`
	GithubPullRequestHeadBranch string `json:"github_pull_request_head_branch"`
	GithubPullRequestBaseBranch string `json:"github_pull_request_base_branch"`
	GithubPullRequestState      string `json:"github_pull_request_state"`
}

// swagger:model
type ModuleDeploymentConfig struct {
	Path                    string `json:"path"`
	GithubRepoName          string `json:"github_repo_name"`
	GithubRepoOwner         string `json:"github_repo_owner"`
	GithubRepoBranch        string `json:"github_repo_branch"`
	GithubAppInstallationID string `json:"github_app_installation_id"`
}

// swagger:model
type CreateModuleRequest struct {
	Name string `json:"name" form:"required,max=255"`

	EnvVars map[string]string `json:"env_vars"`

	ValuesRaw map[string]interface{} `json:"values_raw" form:"required_without=ValuesGithub,omitempty"`

	ValuesGithub *CreateModuleValuesRequestGithub `json:"values_github,omitempty" form:"required_without=ValuesRaw,omitempty"`

	DeploymentGithub *CreateModuleRequestGithub `json:"github,omitempty" form:"omitempty"`

	DeploymentLocal *CreateModuleRequestLocal `json:"local,omitempty" form:"omitempty"`
}

type CreateModuleRequestGithub struct {
	// path to the module in the github repository
	// required: true
	// example: ./staging/eks
	Path string `json:"path" form:"required"`

	// this refers to the Hatchet app installation id, **not** the installation id stored on Github
	// required: true
	// example: bb214807-246e-43a5-a25d-41761d1cff9e
	GithubAppInstallationID string `json:"github_app_installation_id" form:"required,uuid"`

	// the repository owner on Github
	// required: true
	// example: hatchet-dev
	GithubRepositoryOwner string `json:"github_repository_owner" form:"required"`

	// the repository name on Github
	// required: true
	// example: infra
	GithubRepositoryName string `json:"github_repository_name" form:"required"`

	// the repository branch on Github
	// required: true
	// example: main
	GithubRepositoryBranch string `json:"github_repository_branch" form:"required"`
}

type CreateModuleRequestLocal struct {
	// the local path to the module
	LocalPath string `json:"local_path" form:"required"`
}

type CreateModuleValuesRequestGithub struct {
	// path to the module values in the github repository (including file name)
	// required: true
	// example: ./staging/eks
	Path string `json:"path" form:"required"`

	// this refers to the Hatchet app installation id, **not** the installation id stored on Github
	// required: true
	// example: bb214807-246e-43a5-a25d-41761d1cff9e
	GithubAppInstallationID string `json:"github_app_installation_id" form:"required,uuid"`

	// the repository owner on Github
	// required: true
	// example: hatchet-dev
	GithubRepositoryOwner string `json:"github_repository_owner" form:"required"`

	// the repository name on Github
	// required: true
	// example: infra
	GithubRepositoryName string `json:"github_repository_name" form:"required"`

	// the repository branch on Github
	// required: true
	// example: main
	GithubRepositoryBranch string `json:"github_repository_branch" form:"required"`
}

// swagger:model
type CreateModuleResponse struct {
	Module
}

// swagger:model
type GetModuleResponse Module

// swagger:model
type UpdateModuleRequest struct {
	Name string `json:"name" form:"max=255"`

	EnvVars map[string]string `json:"env_vars,omitempty" form:"omitempty"`

	ValuesRaw map[string]interface{} `json:"values_raw,omitempty" form:"omitempty"`

	ValuesGithub *CreateModuleValuesRequestGithub `json:"values_github,omitempty" form:"omitempty"`

	DeploymentGithub *CreateModuleRequestGithub `json:"github,omitempty" form:"omitempty"`
}

// swagger:model
type UpdateModuleResponse Module

// swagger:model
type ForceUnlockModuleResponse Module

// swagger:model
type DeleteModuleResponse Module

// swagger:parameters listModules
type ListModulesRequest struct {
	*PaginationRequest
}

// swagger:model
type ListModulesResponse struct {
	Pagination *PaginationResponse `json:"pagination"`
	Rows       []*Module           `json:"rows"`
}

// swagger:parameters listModuleRuns
type ListModuleRunsRequest struct {
	*PaginationRequest

	// an optional list of statuses for the module run
	// in: query
	Status []ModuleRunStatus `schema:"status" json:"status,omitempty"`

	// an optional list of kinds for the module run
	// in: query
	Kind []ModuleRunKind `schema:"kind" json:"kind,omitempty"`
}

// swagger:model
type ListModuleRunsResponse struct {
	Pagination *PaginationResponse  `json:"pagination"`
	Rows       []*ModuleRunOverview `json:"rows"`
}

// swagger:model
type CreateModuleRunRequest struct {
	Kind ModuleRunKind `json:"kind" form:"required,oneof=plan apply init destroy"`

	Hostname string `json:"hostname"`
}

// swagger:model
type CreateModuleRunResponse struct {
	ModuleRun
}

// swagger:model
type CreateTerraformStateRequest struct {
	ID string `schema:"ID"`
}

// swagger:model
type LockTerraformStateRequest struct {
	ID        string `json:"ID"`
	Operation string `json:"Operation"`
	Info      string `json:"Info"`
	Who       string `json:"Who"`
	Version   string `json:"Version"`
	Created   string `json:"Created"`
	Path      string `json:"Path"`
}

// swagger:model
type LockTerraformStateResponse struct {
	*TerraformLock
}

type TerraformLock struct {
	ID        string `json:"ID"`
	Operation string `json:"Operation"`
	Info      string `json:"Info"`
	Who       string `json:"Who"`
	Version   string `json:"Version"`
	Created   string `json:"Created"`
	Path      string `json:"Path"`
}

// swagger:model
type GetModuleTarballURLResponse struct {
	URL string `json:"url"`
}

// swagger:model
type CreateTerraformPlanRequest struct {
	// the prettified contents of the plan
	// required: true
	PlanPretty string `json:"plan_pretty"`

	// the JSON contents of the plan
	// required: true
	PlanJSON string `json:"plan_json"`
}

type ModuleRunReportKind string

const (
	// ModuleRunReportKindCore reports for the core run -- plan, apply, destroy
	ModuleRunReportKindCore = "core"

	// ModuleRunReportKindMonitor reports for monitor runs
	ModuleRunReportKindMonitor = "monitor"
)

// swagger:model
type FinalizeModuleRunRequest struct {
	// the status of the module run
	// required: true
	Status ModuleRunStatus `json:"status" form:"required"`

	// the description for the module run status
	// required: true
	Description string `json:"description"`

	// the report kind for the finalizer request
	// required: true
	ReportKind ModuleRunReportKind `json:"report_kind" form:"required,oneof=core monitor"`
}

// swagger:model
type FinalizeModuleRunResponse ModuleRunOverview

// swagger:model
type GetModuleRunResponse ModuleRun

// swagger:parameters getModuleTarballURL
type GetModuleTarballURLRequest struct {
	// the SHA to get the tarball from
	// name: github_sha
	// in: query
	GithubSHA string `schema:"github_sha" json:"github_sha"`
}

// swagger:parameters getCurrentModuleValues
type GetModuleValuesRequest struct {
	// the SHA to get the module values file from
	// name: github_sha
	// in: query
	GithubSHA string `schema:"github_sha" json:"github_sha"`
}

// swagger:model
type GetModuleValuesCurrentResponse map[string]interface{}

// swagger:model
type ModulePlanSummary []ModulePlannedChangeSummary

// swagger:model
type ModulePlannedChangeSummary struct {
	Address string   `json:"address"`
	Actions []string `json:"actions"`
}

// swagger:model
type GetModulePlanSummaryResponse []ModulePlannedChangeSummary

// swagger:model
type ModuleEnvVarsVersion struct {
	*APIResourceMeta

	Version uint `json:"version"`

	EnvVars []ModuleEnvVar `json:"env_vars"`
}

type ModuleEnvVar struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

// swagger:model
type GetModuleEnvVarsVersionResponse ModuleEnvVarsVersion

// swagger:model
type ModuleValues struct {
	*APIResourceMeta

	Version uint `json:"version"`

	// Github-based values
	Github *ModuleValuesGithubConfig `json:"github,omitempty"`

	// Raw values (may be omitted)
	Values map[string]interface{} `json:"raw_values,omitempty"`
}

// swagger:model
type ModuleValuesGithubConfig struct {
	Path                    string `json:"path"`
	GithubRepoName          string `json:"github_repo_name"`
	GithubRepoOwner         string `json:"github_repo_owner"`
	GithubRepoBranch        string `json:"github_repo_branch"`
	GithubAppInstallationID string `json:"github_app_installation_id"`
}

// swagger:model
type GetModuleValuesResponse ModuleValues

// swagger:model
type GetModuleRunTokenResponse struct {
	Token string `json:"token"`
}
