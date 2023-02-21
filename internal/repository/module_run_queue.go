package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// ModuleRunQueueRepository represents the set of queries on the ModuleRunQueue model
type ModuleRunQueueRepository interface {
	// --- Module run queue queries ---
	//
	// CreateModuleRunQueue creates a new module run queue in the database, associating it
	// with the parent module
	CreateModuleRunQueue(mod *models.Module, queue *models.ModuleRunQueue) (*models.ModuleRunQueue, RepositoryError)

	// ReadModuleByID reads the module by its unique UUID
	ReadModuleRunQueueByID(moduleID, moduleRunQueueID, lockID string) (*models.ModuleRunQueue, RepositoryError)

	// ListModulesWithQueueItems lists all modules with at least one (non-deleted) queue item
	ListModulesWithQueueItems(opts ...QueryOption) ([]*models.Module, *PaginatedResult, RepositoryError)

	// --- Module run queue item queries ---
	//
	// CreateModuleRunQueue creates a new module run queue in the database, associating it
	// with the parent module
	CreateModuleRunQueueItem(queue *models.ModuleRunQueue, item *models.ModuleRunQueueItem) (*models.ModuleRunQueueItem, RepositoryError)

	// ReadModuleRunQueueItemByModuleRunID reads the first module run queue item corresponding to that module run id
	ReadModuleRunQueueItemByModuleRunID(moduleRunID string) (*models.ModuleRunQueueItem, RepositoryError)

	// DeleteModuleRunQueueItem soft-deletes a module run queue item in the database
	DeleteModuleRunQueueItem(item *models.ModuleRunQueueItem) (*models.ModuleRunQueueItem, RepositoryError)
}
