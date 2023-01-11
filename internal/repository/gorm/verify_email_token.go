package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

// VerifyEmailTokenRepository uses gorm.DB for querying the database
type VerifyEmailTokenRepository struct {
	db *gorm.DB
}

// NewVerifyEmailTokenRepository returns a VerifyEmailTokenRepository which uses
// gorm.DB for querying the database
func NewVerifyEmailTokenRepository(db *gorm.DB) repository.VerifyEmailTokenRepository {
	return &VerifyEmailTokenRepository{db}
}

// CreateVerifyEmailToken creates a new PRT in the database
func (repo *VerifyEmailTokenRepository) CreateVerifyEmailToken(vet *models.VerifyEmailToken) (*models.VerifyEmailToken, repository.RepositoryError) {
	if err := repo.db.Create(vet).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return vet, nil
}

func (repo *VerifyEmailTokenRepository) ReadVerifyEmailTokenByEmailAndTokenID(email, tokID string) (*models.VerifyEmailToken, repository.RepositoryError) {
	vet := &models.VerifyEmailToken{}

	if err := repo.db.Where("email = ? AND id = ?", email, tokID).First(&vet).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return vet, nil

}

func (repo *VerifyEmailTokenRepository) UpdateVerifyEmailToken(vet *models.VerifyEmailToken) (*models.VerifyEmailToken, repository.RepositoryError) {
	if err := repo.db.Save(vet).Error; err != nil {
		return nil, err
	}

	return vet, nil
}
