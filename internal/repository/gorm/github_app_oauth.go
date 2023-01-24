package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

// GithubAppOAuthRepository uses gorm.DB for querying the database
type GithubAppOAuthRepository struct {
	db  *gorm.DB
	key *[32]byte
}

// NewGithubAppOAuthRepository returns a DefaultGithubAppOAuthRepository which uses
// gorm.DB for querying the database
func NewGithubAppOAuthRepository(db *gorm.DB, key *[32]byte) repository.GithubAppOAuthRepository {
	return &GithubAppOAuthRepository{db, key}
}

func (repo *GithubAppOAuthRepository) CreateGithubAppOAuth(gao *models.GithubAppOAuth) (*models.GithubAppOAuth, repository.RepositoryError) {
	err := gao.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Create(gao).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	// return the oauth credential unencrypted
	err = gao.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return gao, nil
}

func (repo *GithubAppOAuthRepository) ReadGithubAppOAuthByGithubUserID(githubUserID int64) (*models.GithubAppOAuth, repository.RepositoryError) {
	gao := &models.GithubAppOAuth{}

	if err := repo.db.Where("github_user_id = ?", githubUserID).First(&gao).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	err := gao.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return gao, nil
}

func (repo *GithubAppOAuthRepository) ReadGithubAppOAuthByUserID(userID string) (*models.GithubAppOAuth, repository.RepositoryError) {
	gao := &models.GithubAppOAuth{}

	if err := repo.db.Where("user_id = ?", userID).First(&gao).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	err := gao.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return gao, nil
}

func (repo *GithubAppOAuthRepository) UpdateGithubAppOAuth(
	gao *models.GithubAppOAuth,
) (*models.GithubAppOAuth, repository.RepositoryError) {
	err := gao.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Save(gao).Error; err != nil {
		return nil, err
	}

	// return the oauth credential unencrypted
	err = gao.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return gao, nil
}
