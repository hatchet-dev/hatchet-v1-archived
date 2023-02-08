package types

const (
	URLParamModuleID    URLParam = "module_id"
	URLParamModuleRunID URLParam = "module_run_id"
)

// swagger:model
type Module struct {
	*APIResourceMeta

	// the name for the module
	// example: eks
	Name string `json:"name"`

	DeploymentConfig ModuleDeploymentConfig `json:"deployment"`
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
)

// swagger:model
type ModuleRun struct {
	*APIResourceMeta

	Status            ModuleRunStatus `json:"status"`
	StatusDescription string          `json:"status_description"`
	Kind              ModuleRunKind   `json:"kind"`
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

	ValuesRaw map[string]interface{} `json:"values_raw"`

	ValuesGithub *CreateModuleValuesRequestGithub `json:"values_github,omitempty" form:"dive"`

	DeploymentGithub *CreateModuleRequestGithub `json:"github,omitempty" form:"dive"`
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
type CreateModuleResponse Module

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

	// the status of the module run
	// in: query
	Status ModuleRunStatus `schema:"status" json:"status"`
}

// swagger:model
type ListModuleRunsResponse struct {
	Pagination *PaginationResponse `json:"pagination"`
	Rows       []*ModuleRun        `json:"rows"`
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

// swagger:model
type FinalizeModuleRunRequest struct {
	// the status of the module run
	// required: true
	Status ModuleRunStatus `json:"status" form:"required"`

	// the description for the module run status
	// required: true
	Description string `json:"description"`
}

// swagger:model
type FinalizeModuleRunResponse ModuleRun

// swagger:parameters getModuleTarballURL
type GetModuleTarballURLRequest struct {
	// the SHA to get the tarball from
	// name: github_sha
	// in: query
	GithubSHA string `schema:"github_sha" json:"github_sha"`
}

// swagger:parameters getModuleValues
type GetModuleValuesRequest struct {
	// the SHA to get the module values file from
	// name: github_sha
	// in: query
	GithubSHA string `schema:"github_sha" json:"github_sha"`
}

// swagger:model
type GetModuleValuesResponse map[string]interface{}
