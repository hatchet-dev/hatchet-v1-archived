package queuemanager

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

const MAX_QUEUE_LENGTH = 50

var MaxQueueLenError = fmt.Errorf("max queue length exceeded")

type ModuleRunQueueManager interface {
	FlushQueue(module *models.Module, lockOpts *LockOpts) error
	Enqueue(module *models.Module, moduleRun *models.ModuleRun, lockOpts *LockOpts) error
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

type LockOpts struct {
	LockID   string
	LockKind models.ModuleLockKind
}

func (queue *DefaultModuleRunQueueManager) Enqueue(module *models.Module, moduleRun *models.ModuleRun, lockOpts *LockOpts) error {
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
		runQueue, err = queue.repo.ModuleRunQueue().ReadModuleRunQueueByID(module.ID, module.ModuleRunQueueID, "")

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
		ModuleRunKind:    moduleRun.Kind,
	}

	runQueueItem.LockID = lockOpts.LockID
	runQueueItem.LockKind = lockOpts.LockKind

	if lockOpts.LockID != "" {
		runQueueItem.LockPriority = models.HasLockID
	} else {
		runQueueItem.LockPriority = models.NoLockID
	}

	switch moduleRun.Kind {
	case models.ModuleRunKindApply:
		runQueueItem.Priority = models.ModuleQueuePriorityApply
	case models.ModuleRunKindPlan:
		runQueueItem.Priority = models.ModuleQueuePriorityPlan
	case models.ModuleRunKindDestroy:
		runQueueItem.Priority = models.ModuleQueuePriorityDestroy
	}

	queueItemsToRemove := make([]models.ModuleRunQueueItem, 0)

	// if there are queue items that have the same lock id and kind, remove those queue items
	if lockOpts.LockID != "" {
		allQueueItems, err := queue.repo.ModuleRunQueue().ReadModuleRunQueueByID(module.ID, runQueue.ID, lockOpts.LockID)

		if err != nil {
			return err
		}

		for _, q := range allQueueItems.Items {
			if q.ModuleRunKind == runQueueItem.ModuleRunKind {
				queueItemsToRemove = append(queueItemsToRemove, q)
			}
		}
	}

	runQueueItem, err = queue.repo.ModuleRunQueue().CreateModuleRunQueueItem(runQueue, runQueueItem)

	if err != nil {
		return err
	}

	// go through queue items to remove
	for _, q := range queueItemsToRemove {
		_, err = queue.repo.ModuleRunQueue().DeleteModuleRunQueueItem(&q)

		if err != nil {
			return err
		}
	}

	return nil
}

func (queue *DefaultModuleRunQueueManager) FlushQueue(module *models.Module, lockOpts *LockOpts) error {
	runQueue, err := queue.repo.ModuleRunQueue().ReadModuleRunQueueByID(module.ID, module.ModuleRunQueueID, lockOpts.LockID)

	if err != nil {
		return err
	}

	for _, q := range runQueue.Items {
		_, err = queue.repo.ModuleRunQueue().DeleteModuleRunQueueItem(&q)

		if err != nil {
			return err
		}
	}

	// if the module's lock currently corresponds to that lock id, remove the lock on the module
	if module.LockID == lockOpts.LockID {
		module.LockID = ""
		module.LockKind = models.ModuleLockKind("")

		module, err = queue.repo.Module().UpdateModule(module)

		if err != nil {
			return err
		}
	}

	return nil
}

func (queue *DefaultModuleRunQueueManager) Len(module *models.Module) (int, error) {
	runQueue, err := queue.repo.ModuleRunQueue().ReadModuleRunQueueByID(module.ID, module.ModuleRunQueueID, "")

	if err != nil {
		return -1, err
	}

	return len(runQueue.Items), nil
}

// NOTE: Peek has the side effect of placing a lock on the module if the module run is currently in a queued state.
// Module locks are fully managed by the queue manager.
func (queue *DefaultModuleRunQueueManager) Peek(module *models.Module) (*models.ModuleRun, error) {
	runQueue, err := queue.repo.ModuleRunQueue().ReadModuleRunQueueByID(module.ID, module.ModuleRunQueueID, module.LockID)

	if err != nil {
		return nil, err
	}

	if len(runQueue.Items) == 0 {
		return nil, nil
	}

	// the first item has highest priority
	moduleRunID := runQueue.Items[0].ModuleRunID
	item := runQueue.Items[0]

	mr, err := queue.repo.Module().ReadModuleRunByID(module.ID, moduleRunID)

	if err != nil {
		return nil, err
	}

	if mr.Status == models.ModuleRunStatusQueued {
		module.LockID = item.LockID
		module.LockKind = item.LockKind
		module, err = queue.repo.Module().UpdateModule(module)

		if err != nil {
			return nil, err
		}
	}

	return mr, err
}

// Remove should remove the lock on the module if the module run corresponds to an apply or destroy.
func (queue *DefaultModuleRunQueueManager) Remove(module *models.Module, moduleRun *models.ModuleRun) error {
	runQueueItem, err := queue.repo.ModuleRunQueue().ReadModuleRunQueueItemByModuleRunID(moduleRun.ID)

	if err != nil {
		return err
	}

	if moduleRun.Kind == models.ModuleRunKindApply || moduleRun.Kind == models.ModuleRunKindDestroy {
		module.LockID = ""
		module.LockKind = models.ModuleLockKind("")

		module, err = queue.repo.Module().UpdateModule(module)

		if err != nil {
			return err
		}
	}

	runQueueItem, err = queue.repo.ModuleRunQueue().DeleteModuleRunQueueItem(runQueueItem)

	if err != nil {
		return err
	}

	return nil
}
