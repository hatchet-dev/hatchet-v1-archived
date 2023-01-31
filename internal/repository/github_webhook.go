package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// GithubWebhookRepository represents the set of queries on the GithubWebhook model
type GithubWebhookRepository interface {
	CreateGithubWebhook(gw *models.GithubWebhook) (*models.GithubWebhook, RepositoryError)
	ReadGithubWebhookByID(teamID, id string) (*models.GithubWebhook, RepositoryError)
	ReadGithubWebhookByTeamID(teamID, repoOwner, repoName string) (*models.GithubWebhook, RepositoryError)
	UpdateGithubWebhook(gw *models.GithubWebhook) (*models.GithubWebhook, RepositoryError)
	DeleteGithubWebhook(gw *models.GithubWebhook) (*models.GithubWebhook, RepositoryError)

	AppendGithubAppInstallation(gw *models.GithubWebhook, gai *models.GithubAppInstallation) (*models.GithubWebhook, RepositoryError)
	RemoveGithubAppInstallation(gw *models.GithubWebhook, gai *models.GithubAppInstallation) (*models.GithubWebhook, RepositoryError)
}
