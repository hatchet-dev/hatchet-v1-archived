package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
)

// ModuleRunQueueRepository uses gorm.DB for querying the database
type ModuleRunQueueRepository struct {
	db *gorm.DB
}

// NewModuleRunQueueRepository returns a DefaultModuleRunQueueRepository which uses
// gorm.DB for querying the database
func NewModuleRunQueueRepository(db *gorm.DB) repository.ModuleRunQueueRepository {
	return &ModuleRunQueueRepository{db}
}

func (repo *ModuleRunQueueRepository) CreateModuleRunQueue(mod *models.Module, queue *models.ModuleRunQueue) (*models.ModuleRunQueue, repository.RepositoryError) {
	if err := repo.db.Create(queue).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return queue, nil
}

func (repo *ModuleRunQueueRepository) ReadModuleRunQueueByID(moduleID, moduleRunQueueID string) (*models.ModuleRunQueue, repository.RepositoryError) {
	queue := &models.ModuleRunQueue{}
	query := repo.db

	query = query.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("module_run_queue_items.priority DESC").Order("module_run_queue_items.created_at DESC")
	})

	if err := query.Where("module_id = ? AND id = ?", moduleID, moduleRunQueueID).First(&queue).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return queue, nil
}

func (repo *ModuleRunQueueRepository) ListModulesWithQueueItems(opts ...repository.QueryOption) ([]*models.Module, *repository.PaginatedResult, repository.RepositoryError) {
	var mods []*models.Module

	queryOpts := make([]repository.QueryOption, 0)

	queryOpts = append(queryOpts, opts...)

	// overwrite sort by for subquery
	queryOpts = append(queryOpts, repository.WithSortBy("modules.updated_at"))

	subQuery := repo.db.Model(&models.Module{}).
		Joins("INNER JOIN module_run_queues ON module_run_queues.module_id = modules.id").
		Joins("INNER JOIN module_run_queue_items ON module_run_queue_items.module_run_queue_id = module_run_queues.id AND module_run_queue_items.deleted_at IS NULL").
		Select("modules.id")

	// apply pagination to subquery as it's on the modules table as well
	subQuery = queryutils.ApplyOpts(subQuery, queryutils.ComputeQuery(queryOpts...))

	db := repo.db.Model(&models.Module{}).Where(`modules.id IN (?)`, subQuery)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(queryOpts, db, paginatedResult))

	if err := db.Find(&mods).Error; err != nil {
		return nil, nil, err
	}

	return mods, paginatedResult, nil
}

func (repo *ModuleRunQueueRepository) CreateModuleRunQueueItem(queue *models.ModuleRunQueue, item *models.ModuleRunQueueItem) (*models.ModuleRunQueueItem, repository.RepositoryError) {
	if err := repo.db.Create(item).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return item, nil
}

func (repo *ModuleRunQueueRepository) ReadModuleRunQueueItemByModuleRunID(moduleRunID string) (*models.ModuleRunQueueItem, repository.RepositoryError) {
	queueItem := &models.ModuleRunQueueItem{}

	if err := repo.db.Where("module_run_id = ?", moduleRunID).First(&queueItem).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return queueItem, nil
}

func (repo *ModuleRunQueueRepository) DeleteModuleRunQueueItem(item *models.ModuleRunQueueItem) (*models.ModuleRunQueueItem, repository.RepositoryError) {
	del := repo.db.Delete(&item)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return item, nil
}
