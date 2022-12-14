package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// PersonalAccessTokenRepository represents the set of queries on the PersonalAccessToken model
type PersonalAccessTokenRepository interface {
	CreatePersonalAccessToken(pat *models.PersonalAccessToken) (*models.PersonalAccessToken, RepositoryError)
	ReadPersonalAccessToken(userID, tokenID string) (*models.PersonalAccessToken, RepositoryError)
	ListPersonalAccessTokensByUserID(userID string, opts ...QueryOption) ([]*models.PersonalAccessToken, *PaginatedResult, RepositoryError)
	UpdatePersonalAccessToken(pat *models.PersonalAccessToken) (*models.PersonalAccessToken, RepositoryError)
	DeletePersonalAccessToken(pat *models.PersonalAccessToken) (*models.PersonalAccessToken, RepositoryError)
}
