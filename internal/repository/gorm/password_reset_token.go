package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

// PasswordResetTokenRepository uses gorm.DB for querying the database
type PasswordResetTokenRepository struct {
	db *gorm.DB
}

// NewPasswordResetTokenRepository returns a PasswordResetTokenRepository which uses
// gorm.DB for querying the database
func NewPasswordResetTokenRepository(db *gorm.DB) repository.PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{db}
}

// CreatePasswordResetToken creates a new PRT in the database
func (repo *PasswordResetTokenRepository) CreatePasswordResetToken(prt *models.PasswordResetToken) (*models.PasswordResetToken, repository.RepositoryError) {
	if err := repo.db.Create(prt).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return prt, nil
}

func (repo *PasswordResetTokenRepository) ReadPasswordResetTokenByEmailAndTokenID(email, tokID string) (*models.PasswordResetToken, repository.RepositoryError) {
	prt := &models.PasswordResetToken{}

	if err := repo.db.Where("email = ? AND id = ?", email, tokID).First(&prt).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return prt, nil

}

func (repo *PasswordResetTokenRepository) UpdatePasswordResetToken(prt *models.PasswordResetToken) (*models.PasswordResetToken, repository.RepositoryError) {
	if err := repo.db.Save(prt).Error; err != nil {
		return nil, err
	}

	return prt, nil

}
