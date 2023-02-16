package modulequeuechecker

import (
	"time"

	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/queuemanager"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulerunner"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type ModuleQueueChecker struct {
	queueManager queuemanager.ModuleRunQueueManager
	repo         repository.Repository
}

func NewModuleQueueChecker(queueManager queuemanager.ModuleRunQueueManager, repo repository.Repository) *ModuleQueueChecker {
	return &ModuleQueueChecker{queueManager, repo}
}

type CheckQueueInput struct {
	TeamID, ModuleID string
}

func (mqc *ModuleQueueChecker) ScheduleFromQueue(ctx workflow.Context, input CheckQueueInput) (string, error) {
	queue := mqc.queueManager
	repo := mqc.repo

	module, err := repo.Module().ReadModuleByID(input.TeamID, input.ModuleID)

	if err != nil {
		return "", err
	}

	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    100 * time.Second,
		MaximumAttempts:    0,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retrypolicy,
	}

	currModuleRun, err := queue.Peek(module)

	if err != nil {
		return "", err
	}

	if currModuleRun == nil {
		return "no_queued_runs", nil
	}

	if currModuleRun.Status == models.ModuleRunStatusInProgress {
		return "run_in_progress", nil
	}

	if currModuleRun.Status == models.ModuleRunStatusFailed || currModuleRun.Status == models.ModuleRunStatusCompleted {
		// remove from queue and peek again
		err = queue.Remove(module, currModuleRun)

		if err != nil {
			return "", err
		}

		currModuleRun, err = queue.Peek(module)

		if err != nil {
			return "", err
		}

		if currModuleRun == nil {
			return "no_queued_runs", nil
		}
	}

	// if we've reached this point, we set the status to in progress and we trigger the run
	currModuleRun.Status = models.ModuleRunStatusInProgress

	currModuleRun, err = repo.Module().UpdateModuleRun(currModuleRun)

	if err != nil {
		return "", err
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var runOutput string

	err = workflow.ExecuteActivity(ctx, modulerunner.Run, modulerunner.RunInput{
		TeamID:      module.TeamID,
		ModuleID:    module.ID,
		ModuleRunID: currModuleRun.ID,
	}).Get(ctx, &runOutput)

	if err != nil {
		return "", err
	}

	return runOutput, nil
}
