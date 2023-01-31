package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

// GithubWebhookRepository uses gorm.DB for querying the database
type GithubWebhookRepository struct {
	db  *gorm.DB
	key *[32]byte
}

// NewGithubWebhookRepository returns a DefaultGithubWebhookRepository which uses
// gorm.DB for querying the database
func NewGithubWebhookRepository(db *gorm.DB, key *[32]byte) repository.GithubWebhookRepository {
	return &GithubWebhookRepository{db, key}
}

func (repo *GithubWebhookRepository) CreateGithubWebhook(gw *models.GithubWebhook) (*models.GithubWebhook, repository.RepositoryError) {
	err := gw.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Create(gw).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	// return the oauth credential unencrypted
	err = gw.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return gw, nil
}

func (repo *GithubWebhookRepository) ReadGithubWebhookByID(teamID, id string) (*models.GithubWebhook, repository.RepositoryError) {
	gw := &models.GithubWebhook{}

	if err := repo.db.Where("team_id = ? AND id = ?", teamID, id).First(&gw).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	err := gw.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return gw, nil
}

func (repo *GithubWebhookRepository) ReadGithubWebhookByTeamID(teamID, repoOwner, repoName string) (*models.GithubWebhook, repository.RepositoryError) {
	gw := &models.GithubWebhook{}

	if err := repo.db.Where("team_id = ? AND github_repository_owner = ? AND github_repository_name = ?", teamID, repoOwner, repoName).First(&gw).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	err := gw.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return gw, nil
}

func (repo *GithubWebhookRepository) UpdateGithubWebhook(
	gw *models.GithubWebhook,
) (*models.GithubWebhook, repository.RepositoryError) {
	err := gw.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Save(gw).Error; err != nil {
		return nil, err
	}

	// return the oauth credential unencrypted
	err = gw.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return gw, nil
}

func (repo *GithubWebhookRepository) DeleteGithubWebhook(gw *models.GithubWebhook) (*models.GithubWebhook, repository.RepositoryError) {
	// encrypt the PAT
	err := gw.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Delete(gw).Error; err != nil {
		return nil, err
	}

	// return the PAT decrypted
	err = gw.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return gw, nil
}

func (repo *GithubWebhookRepository) AppendGithubAppInstallation(gw *models.GithubWebhook, gai *models.GithubAppInstallation) (*models.GithubWebhook, repository.RepositoryError) {
	panic("unimplemented")
}

func (repo *GithubWebhookRepository) RemoveGithubAppInstallation(gw *models.GithubWebhook, gai *models.GithubAppInstallation) (*models.GithubWebhook, repository.RepositoryError) {
	panic("unimplemented")
}
