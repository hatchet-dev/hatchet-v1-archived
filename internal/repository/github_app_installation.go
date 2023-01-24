package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// GithubAppInstallationRepository represents the set of queries on the GithubAppInstallation model
type GithubAppInstallationRepository interface {
	CreateGithubAppInstallation(gai *models.GithubAppInstallation) (*models.GithubAppInstallation, RepositoryError)
	ReadGithubAppInstallationByInstallationAndAccountID(installationID, accountID int64) (*models.GithubAppInstallation, RepositoryError)
	ListGithubAppInstallationsByUserID(userID string, opts ...QueryOption) ([]*models.GithubAppInstallation, *PaginatedResult, RepositoryError)
	UpdateGithubAppInstallation(gai *models.GithubAppInstallation) (*models.GithubAppInstallation, RepositoryError)
	DeleteGithubAppInstallation(gai *models.GithubAppInstallation) (*models.GithubAppInstallation, RepositoryError)
}
