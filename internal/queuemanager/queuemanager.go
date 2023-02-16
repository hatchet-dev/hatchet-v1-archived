package queuemanager

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

const MAX_QUEUE_LENGTH = 50

var MaxQueueLenError = fmt.Errorf("max queue length exceeded")

type ModuleRunQueueManager interface {
	Enqueue(module *models.Module, moduleRun *models.ModuleRun) error
	Len(module *models.Module) (int, error)
	Peek(module *models.Module) (*models.ModuleRun, error)
	Remove(module *models.Module, moduleRun *models.ModuleRun) error
}

type DefaultModuleRunQueueManager struct {
	repo repository.Repository
}

func NewDefaultModuleRunQueueManager(repo repository.Repository) *DefaultModuleRunQueueManager {
	return &DefaultModuleRunQueueManager{repo}
}

func (queue *DefaultModuleRunQueueManager) Enqueue(module *models.Module, moduleRun *models.ModuleRun) error {
	var runQueue *models.ModuleRunQueue
	var err error

	// create a queue for the module, if it does not exist
	if module.ModuleRunQueueID == "" {
		runQueue, err = queue.repo.ModuleRunQueue().CreateModuleRunQueue(module, &models.ModuleRunQueue{
			ModuleID: module.ID,
			Items:    []models.ModuleRunQueueItem{},
		})

		if err != nil {
			return err
		}

		module.ModuleRunQueueID = runQueue.ID

		module, err = queue.repo.Module().UpdateModule(module)

		if err != nil {
			return err
		}
	} else {
		runQueue, err = queue.repo.ModuleRunQueue().ReadModuleRunQueueByID(module.ID, module.ModuleRunQueueID)

		if err != nil {
			return err
		}
	}

	if len(runQueue.Items) >= MAX_QUEUE_LENGTH {
		return MaxQueueLenError
	}

	runQueueItem := &models.ModuleRunQueueItem{
		ModuleRunID:      moduleRun.ID,
		ModuleRunQueueID: runQueue.ID,
	}

	switch moduleRun.Kind {
	case models.ModuleRunKindApply:
		runQueueItem.Priority = models.ModuleQueuePriorityApply
	case models.ModuleRunKindPlan:
		runQueueItem.Priority = models.ModuleQueuePriorityPlan
	case models.ModuleRunKindDestroy:
		runQueueItem.Priority = models.ModuleQueuePriorityDestroy
	}

	runQueueItem, err = queue.repo.ModuleRunQueue().CreateModuleRunQueueItem(runQueue, runQueueItem)

	if err != nil {
		return err
	}

	return nil
}

func (queue *DefaultModuleRunQueueManager) Len(module *models.Module) (int, error) {
	runQueue, err := queue.repo.ModuleRunQueue().ReadModuleRunQueueByID(module.ID, module.ModuleRunQueueID)

	if err != nil {
		return -1, err
	}

	return len(runQueue.Items), nil
}

func (queue *DefaultModuleRunQueueManager) Peek(module *models.Module) (*models.ModuleRun, error) {
	runQueue, err := queue.repo.ModuleRunQueue().ReadModuleRunQueueByID(module.ID, module.ModuleRunQueueID)

	if err != nil {
		return nil, err
	}

	if len(runQueue.Items) == 0 {
		return nil, nil
	}

	// the first item has highest priority
	moduleRunID := runQueue.Items[0].ModuleRunID

	return queue.repo.Module().ReadModuleRunByID(module.ID, moduleRunID)
}

func (queue *DefaultModuleRunQueueManager) Remove(module *models.Module, moduleRun *models.ModuleRun) error {
	runQueueItem, err := queue.repo.ModuleRunQueue().ReadModuleRunQueueItemByModuleRunID(moduleRun.ID)

	if err != nil {
		return err
	}

	runQueueItem, err = queue.repo.ModuleRunQueue().DeleteModuleRunQueueItem(runQueueItem)

	if err != nil {
		return err
	}

	return nil
}
