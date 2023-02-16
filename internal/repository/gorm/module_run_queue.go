package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
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

	if err := query.Where("module_id = ? AND module_run_queue_id = ?", moduleID, moduleRunQueueID).First(&queue).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return queue, nil
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
