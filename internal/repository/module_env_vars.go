package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// ModuleEnvVarsRepository represents the set of queries on the ModuleEnvVarsVersion model
type ModuleEnvVarsRepository interface {
	// --- Module values version queries ---
	//
	// CreateModuleEnvVarsVersion creates a new module values version in the database
	CreateModuleEnvVarsVersion(mvv *models.ModuleEnvVarsVersion) (*models.ModuleEnvVarsVersion, RepositoryError)

	// ReadModuleEnvVarsVersionByID reads the module by its unique UUID
	ReadModuleEnvVarsVersionByID(moduleID, moduleValuesVersionID string) (*models.ModuleEnvVarsVersion, RepositoryError)

	// ListModuleValueVersionsByModuleID lists all module value versions for a module
	ListModuleValueVersionsByModuleID(moduleID string, opts ...QueryOption) ([]*models.ModuleEnvVarsVersion, *PaginatedResult, RepositoryError)

	// DeleteModuleEnvVarsVersion soft-deletes a module values version
	DeleteModuleEnvVarsVersion(mvv *models.ModuleEnvVarsVersion) (*models.ModuleEnvVarsVersion, RepositoryError)
}
