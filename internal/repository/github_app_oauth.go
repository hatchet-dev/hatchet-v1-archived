package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// GithubAppOAuthRepository represents the set of queries on the GithubAppOAuth model
type GithubAppOAuthRepository interface {
	CreateGithubAppOAuth(gao *models.GithubAppOAuth) (*models.GithubAppOAuth, RepositoryError)
	ReadGithubAppOAuthByGithubUserID(githubUserID int64) (*models.GithubAppOAuth, RepositoryError)
	ReadGithubAppOAuthByUserID(userID string) (*models.GithubAppOAuth, RepositoryError)
	UpdateGithubAppOAuth(gao *models.GithubAppOAuth) (*models.GithubAppOAuth, RepositoryError)
}
