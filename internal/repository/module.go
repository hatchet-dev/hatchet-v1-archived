package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// ModuleRepository represents the set of queries on the Module model
type ModuleRepository interface {
	// --- Module queries ---
	//
	// CreateModule creates a new module in the database
	CreateModule(mod *models.Module) (*models.Module, RepositoryError)

	// ReadModuleByID reads the module by its unique UUID
	ReadModuleByID(id string) (*models.Module, RepositoryError)

	// ListModulesByTeamID lists all modules for a team
	ListModulesByTeamID(teamID string, opts ...QueryOption) ([]*models.Module, *PaginatedResult, RepositoryError)

	// UpdateModule updates any modified values for a module
	UpdateModule(module *models.Module) (*models.Module, RepositoryError)

	// DeleteModule soft-deletes a module
	DeleteModule(module *models.Module) (*models.Module, RepositoryError)
}
