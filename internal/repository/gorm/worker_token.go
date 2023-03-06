package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

// WorkerTokenRepository uses gorm.DB for querying the database
type WorkerTokenRepository struct {
	db *gorm.DB
}

// NewWorkerTokenRepository returns a DefaultWorkerTokenRepository which uses
// gorm.DB for querying the database
func NewWorkerTokenRepository(db *gorm.DB) repository.WorkerTokenRepository {
	return &WorkerTokenRepository{db}
}

func (repo *WorkerTokenRepository) CreateWorkerToken(wt *models.WorkerToken) (*models.WorkerToken, repository.RepositoryError) {
	if err := repo.db.Create(wt).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return wt, nil
}

func (repo *WorkerTokenRepository) ReadWorkerToken(teamID, tokenID string) (*models.WorkerToken, repository.RepositoryError) {
	wt := &models.WorkerToken{}

	if err := repo.db.Where("team_id = ? AND id = ?", teamID, tokenID).First(&wt).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return wt, nil
}

func (repo *WorkerTokenRepository) UpdateWorkerToken(wt *models.WorkerToken) (*models.WorkerToken, repository.RepositoryError) {
	if err := repo.db.Save(wt).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return wt, nil
}

func (repo *WorkerTokenRepository) DeleteWorkerToken(wt *models.WorkerToken) (*models.WorkerToken, repository.RepositoryError) {
	del := repo.db.Delete(&wt)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return wt, nil
}
