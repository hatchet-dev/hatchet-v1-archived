package types

const (
	URLParamGithubAppInstallationID URLParam = "github_app_installation_id"
	URLParamGithubRepoOwner         URLParam = "github_repo_owner"
	URLParamGithubRepoName          URLParam = "github_repo_name"
)

// swagger:model
type GithubAppInstallation struct {
	*APIResourceMeta

	InstallationID          int64  `json:"installation_id"`
	InstallationSettingsURL string `json:"installation_settings_url"`

	AccountName      string `json:"account_name"`
	AccountAvatarURL string `json:"account_avatar_url"`
}

// swagger:parameters listGithubAppInstallations
type ListGithubAppInstallationsRequest struct {
	*PaginationRequest
}

// swagger:model
type ListGithubAppInstallationsResponse struct {
	Pagination *PaginationResponse      `json:"pagination"`
	Rows       []*GithubAppInstallation `json:"rows"`
}

// swagger:model
type GithubRepo struct {
	RepoOwner string `json:"repo_owner"`
	RepoName  string `json:"repo_name"`
}

// swagger:model
type ListGithubReposResponse []GithubRepo

// swagger:model
type GithubBranch struct {
	BranchName string `json:"branch_name"`
	IsDefault  bool   `json:"is_default"`
}

// swagger:model
type ListGithubRepoBranchesResponse []GithubBranch
