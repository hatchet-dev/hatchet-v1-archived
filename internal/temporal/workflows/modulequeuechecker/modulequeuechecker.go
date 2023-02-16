package modulequeuechecker

import (
	"time"

	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
	"github.com/hatchet-dev/hatchet/internal/provisioner/provisionerutils"
	"github.com/hatchet-dev/hatchet/internal/queuemanager"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulerunner"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	hatchetenums "github.com/hatchet-dev/hatchet/internal/temporal/enums"
)

type ModuleQueueChecker struct {
	queueManager queuemanager.ModuleRunQueueManager
	db           database.Config
	tokenOpts    token.TokenOpts
	serverURL    string
}

func NewModuleQueueChecker(queueManager queuemanager.ModuleRunQueueManager, db database.Config, tokenOpts token.TokenOpts, serverURL string) *ModuleQueueChecker {
	return &ModuleQueueChecker{queueManager, db, tokenOpts, serverURL}
}

type CheckQueueInput struct {
	TeamID, ModuleID string
}

func (mqc *ModuleQueueChecker) ScheduleFromQueue(ctx workflow.Context, input CheckQueueInput) (string, error) {
	queue := mqc.queueManager
	repo := mqc.db.Repository

	team, err := repo.Team().ReadTeamByID(input.TeamID)

	if err != nil {
		return "", err
	}

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

	envOpts, err := provisionerutils.GetProvisionerEnvOpts(team, module, currModuleRun, mqc.db, mqc.tokenOpts, mqc.serverURL)

	if err != nil {
		return "", err
	}

	env, err := provisioner.GetHatchetRunnerEnv(envOpts, []string{})

	if err != nil {
		return "", err
	}

	cwo := workflow.ChildWorkflowOptions{
		WorkflowExecutionTimeout: 1 * time.Minute,
		WorkflowTaskTimeout:      time.Minute,
		ParentClosePolicy:        enums.PARENT_CLOSE_POLICY_ABANDON,
	}

	childCtx := workflow.WithChildOptions(ctx, cwo)

	childWorkflowFuture := workflow.ExecuteChildWorkflow(childCtx, hatchetenums.WorkflowTypeNameProvision, modulerunner.RunInput{
		Kind: currModuleRun.Kind,
		Opts: &provisioner.ProvisionOpts{
			Env: env,
		},
	})

	var childWE workflow.Execution

	if err := childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE); err != nil {
		return "", err
	}

	return "triggered_workflow", nil
}
