package types

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
