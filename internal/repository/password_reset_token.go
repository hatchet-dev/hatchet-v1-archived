package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// PasswordResetTokenRepository represents the set of queries on the PasswordResetToken model
type PasswordResetTokenRepository interface {
	CreatePasswordResetToken(pwt *models.PasswordResetToken) (*models.PasswordResetToken, RepositoryError)
	ReadPasswordResetTokenByEmailAndTokenID(email, tokID string) (*models.PasswordResetToken, RepositoryError)
	UpdatePasswordResetToken(pwt *models.PasswordResetToken) (*models.PasswordResetToken, RepositoryError)
}
