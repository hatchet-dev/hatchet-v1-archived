package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// ModuleValuesRepository represents the set of queries on the ModuleValuesVersion and ModuleValues models
type ModuleValuesRepository interface {
	// --- Module values version queries ---
	//
	// CreateModuleValuesVersion creates a new module values version in the database
	CreateModuleValuesVersion(mvv *models.ModuleValuesVersion) (*models.ModuleValuesVersion, RepositoryError)

	// ReadModuleValuesVersionByID reads the module by its unique UUID
	ReadModuleValuesVersionByID(moduleID, moduleValuesVersionID string) (*models.ModuleValuesVersion, RepositoryError)

	// ListModuleValueVersionsByModuleID lists all module value versions for a module
	ListModuleValueVersionsByModuleID(moduleID string, opts ...QueryOption) ([]*models.ModuleValuesVersion, *PaginatedResult, RepositoryError)

	// DeleteModuleValuesVersion soft-deletes a module values version
	DeleteModuleValuesVersion(mvv *models.ModuleValuesVersion) (*models.ModuleValuesVersion, RepositoryError)

	// --- Module values queries ---
	//
	// CreateModuleValues creates a new module values entry in the database
	CreateModuleValues(mv *models.ModuleValues) (*models.ModuleValues, RepositoryError)

	// ReadModuleValuesByID reads the module values by its unique UUID
	ReadModuleValuesByID(moduleValuesID string) (*models.ModuleValues, RepositoryError)

	// ReadModuleValuesByVersionID finds the first module values entry with the values version ID. There should
	// only be one entry per version since values are immutable.
	ReadModuleValuesByVersionID(moduleValuesVersionID string) (*models.ModuleValues, RepositoryError)

	// DeleteModuleValues soft-deletes a module values version
	DeleteModuleValues(mv *models.ModuleValues) (*models.ModuleValues, RepositoryError)
}
