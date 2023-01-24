package models

// GithubAppOAuth represents a user authenticated to a Github app via oauth
type GithubAppOAuth struct {
	Base
	*SharedOAuthFields

	GithubUserID int64
}
