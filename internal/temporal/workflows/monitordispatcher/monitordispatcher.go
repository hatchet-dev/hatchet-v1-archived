package monitordispatcher

import (
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
	"github.com/hatchet-dev/hatchet/internal/provisioner/provisionerutils"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulerunner"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	hatchetenums "github.com/hatchet-dev/hatchet/internal/temporal/enums"
)

type MonitorDispatcher struct {
	logStore             logstorage.LogStorageBackend
	db                   database.Config
	tokenOpts            token.TokenOpts
	serverURL            string
	broadcastGRPCAddress string
}

func NewMonitorDispatcher(logStore logstorage.LogStorageBackend, db database.Config, tokenOpts token.TokenOpts, serverURL, broadcastGRPCAddress string) *MonitorDispatcher {
	return &MonitorDispatcher{logStore, db, tokenOpts, serverURL, broadcastGRPCAddress}
}

type MonitorDispatcherInput struct {
	TeamID, MonitorID string
}

func (md *MonitorDispatcher) DispatchMonitors(ctx workflow.Context, input MonitorDispatcherInput) (string, error) {
	repo := md.db.Repository

	team, err := repo.Team().ReadTeamByID(input.TeamID)

	if err != nil {
		return "", err
	}

	monitor, err := repo.ModuleMonitor().ReadModuleMonitorByID(team.ID, input.MonitorID)

	if err != nil {
		return "", err
	}

	var mods []*models.Module

	if len(monitor.Modules) > 0 {
		monitorModuleIDs := make([]string, 0)

		for _, monitorMod := range monitor.Modules {
			monitorModuleIDs = append(monitorModuleIDs, monitorMod.ID)
		}

		mods, _, err = repo.Module().ListModulesByIDs(team.ID, monitorModuleIDs)

		if err != nil {
			return "", err
		}
	} else {
		mods, _, err = repo.Module().ListModulesByTeamID(team.ID)

		if err != nil {
			return "", err
		}
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

	// TODO: pagination
	var resErr error

	for _, mod := range mods {
		// create a new module run of kind monitor (this is read-only and thus bypasses the queue)
		run := &models.ModuleRun{
			ModuleID:          mod.ID,
			Status:            models.ModuleRunStatusInProgress,
			StatusDescription: "New monitor triggered", // TODO: better description
			Kind:              models.ModuleRunKindMonitor,
			ModuleRunConfig: models.ModuleRunConfig{
				TriggerKind:            models.ModuleRunTriggerKindManual,
				ModuleValuesVersionID:  mod.CurrentModuleValuesVersionID,
				ModuleEnvVarsVersionID: mod.CurrentModuleEnvVarsVersionID,
			},
			LogLocation: md.logStore.GetID(),
		}

		run, err = repo.Module().CreateModuleRun(run)

		if err != nil {
			resErr = multierror.Append(err)
			continue
		}

		ctx = workflow.WithActivityOptions(ctx, options)

		// trigger child workflow
		envOpts, err := provisionerutils.GetProvisionerEnvOpts(team, mod, run, md.db, md.tokenOpts, md.serverURL, md.broadcastGRPCAddress)

		if err != nil {
			resErr = multierror.Append(err)
			continue
		}

		env, err := provisioner.GetHatchetRunnerEnv(envOpts, []string{})

		if err != nil {
			resErr = multierror.Append(err)
			continue
		}

		cwo := workflow.ChildWorkflowOptions{
			TaskQueue:                hatchetenums.ModuleRunQueueName,
			Namespace:                input.TeamID,
			WorkflowExecutionTimeout: 1 * time.Minute,
			WorkflowTaskTimeout:      time.Minute,
			ParentClosePolicy:        enums.PARENT_CLOSE_POLICY_ABANDON,
		}

		childCtx := workflow.WithChildOptions(ctx, cwo)

		childWorkflowFuture := workflow.ExecuteChildWorkflow(childCtx, hatchetenums.WorkflowTypeNameRunMonitor, modulerunner.MonitorInput{
			ModuleMonitorID: monitor.ID,
			Kind:            monitor.Kind,
			Opts: &provisioner.ProvisionOpts{
				Env: env,
			},
		})

		var childWE workflow.Execution

		if err := childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE); err != nil {
			var resErr error
			resErr = multierror.Append(resErr, err)

			run.Status = models.ModuleRunStatusFailed

			run, err = repo.Module().UpdateModuleRun(run)

			if err != nil {
				err = multierror.Append(resErr, err)
			}

			continue
		}
	}

	if resErr != nil {
		return "", resErr
	}

	return "triggered_workflow", nil
}
