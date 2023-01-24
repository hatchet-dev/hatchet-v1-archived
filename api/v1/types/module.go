package types

// swagger:model
type Module struct {
	*APIResourceMeta

	// the name for the module
	// example: eks
	Name string `json:"name"`

	DeploymentConfig ModuleDeploymentConfig `json:"deployment"`
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
