package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// VerifyEmailTokenRepository represents the set of queries on the VerifyEmailToken model
type VerifyEmailTokenRepository interface {
	CreateVerifyEmailToken(vet *models.VerifyEmailToken) (*models.VerifyEmailToken, RepositoryError)
	ReadVerifyEmailTokenByEmailAndTokenID(email, tokID string) (*models.VerifyEmailToken, RepositoryError)
	UpdateVerifyEmailToken(vet *models.VerifyEmailToken) (*models.VerifyEmailToken, RepositoryError)
}
