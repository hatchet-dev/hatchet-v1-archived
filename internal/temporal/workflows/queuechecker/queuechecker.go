package queuechecker

import (
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type QueueChecker struct {
	repo repository.Repository
	mqc  *modulequeuechecker.ModuleQueueChecker
}

func NewQueueChecker(repo repository.Repository, mqc *modulequeuechecker.ModuleQueueChecker) *QueueChecker {
	return &QueueChecker{repo, mqc}
}

type GlobalQueueInput struct{}

func (qc *QueueChecker) CheckQueues(ctx workflow.Context, input GlobalQueueInput) (string, error) {
	repo := qc.repo

	// list all modules with at least one run in the queue
	modules, _, err := repo.ModuleRunQueue().ListModulesWithQueueItems()

	if err != nil {
		return "", err
	}

	var allErrs error

	// TODO: handle pagination
	for _, module := range modules {
		cwo := workflow.ChildWorkflowOptions{
			WorkflowExecutionTimeout: 1 * time.Minute,
			WorkflowTaskTimeout:      time.Minute,
			ParentClosePolicy:        enums.PARENT_CLOSE_POLICY_ABANDON,
		}

		childCtx := workflow.WithChildOptions(ctx, cwo)

		childWorkflowFuture := workflow.ExecuteChildWorkflow(childCtx, qc.mqc.ScheduleFromQueue, modulequeuechecker.CheckQueueInput{
			TeamID:   module.TeamID,
			ModuleID: module.ID,
		})

		var childWE workflow.Execution

		if err := childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE); err != nil {
			err = multierror.Append(allErrs, err)
		}
	}

	if allErrs != nil {
		return "", allErrs
	}

	return "success", nil
}
