package repository

import "github.com/hatchet-dev/hatchet/internal/models"

type WorkerTokenRepository interface {
	// CreateWorkerToken creates a new module run token in the database
	CreateWorkerToken(wt *models.WorkerToken) (*models.WorkerToken, RepositoryError)

	// ReadWorkerToken reads the worker token by its token ID
	ReadWorkerToken(teamID, tokenID string) (*models.WorkerToken, RepositoryError)

	// UpdateWorkerToken updates a module run token
	UpdateWorkerToken(wt *models.WorkerToken) (*models.WorkerToken, RepositoryError)

	// DeleteWorkerToken soft-deletes a module run token in the DB
	DeleteWorkerToken(wt *models.WorkerToken) (*models.WorkerToken, RepositoryError)
}
