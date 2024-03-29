package runutils

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
	"github.com/hatchet-dev/hatchet/internal/models"
)

func GenerateRunDescription(
	config *server.Config,
	module *models.Module,
	run *models.ModuleRun,
	status models.ModuleRunStatus,
	failedMonitorResult *models.ModuleMonitorResult,
	pr vcs.VCSRepositoryPullRequest,
) (string, error) {
	switch run.Kind {
	case models.ModuleRunKindPlan:
		return generatePlanRunDescription(config, module, run, status, failedMonitorResult, pr)
	case models.ModuleRunKindApply:
		return generateApplyRunDescription(config, module, run, status, failedMonitorResult)
	case models.ModuleRunKindDestroy:
		return generateDestroyRunDescription(config, module, run, status, failedMonitorResult)
	case models.ModuleRunKindMonitor:
		return generateMonitorRunDescription(config, module, run, status, failedMonitorResult)
	case models.ModuleRunKindInit:
		return generateInitRunDescription(config, module, run, status)
	}

	return "", fmt.Errorf("unknown run kind %s", run.Kind)
}

func generatePlanRunDescription(
	config *server.Config,
	module *models.Module,
	run *models.ModuleRun,
	status models.ModuleRunStatus,
	failedMonitorResult *models.ModuleMonitorResult,
	pr vcs.VCSRepositoryPullRequest,
) (string, error) {
	prefix := "Plan"

	if run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindVCS && pr != nil {
		prefix = fmt.Sprintf("Plan for pull request %s/%s #%d", pr.GetRepoOwner(), pr.GetRepoName(), pr.GetPRNumber())
	}

	if failedMonitorResult != nil {
		return fmt.Sprintf("%s failed a monitor check with message: %s", prefix, failedMonitorResult.Message), nil
	}

	switch status {
	case models.ModuleRunStatusCompleted:
		return fmt.Sprintf("%s ran successfully", prefix), nil
	case models.ModuleRunStatusFailed:
		return fmt.Sprintf("%s failed", prefix), nil
	case models.ModuleRunStatusInProgress:
		return fmt.Sprintf("%s is in progress", prefix), nil
	case models.ModuleRunStatusQueued:
		return fmt.Sprintf("%s is queued", prefix), nil
	}

	return "", nil
}

func generateApplyRunDescription(
	config *server.Config,
	module *models.Module,
	run *models.ModuleRun,
	status models.ModuleRunStatus,
	failedMonitorResult *models.ModuleMonitorResult,
) (string, error) {
	prefix := "Apply"

	if run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindVCS {
		prefix = fmt.Sprintf("Apply for branch %s", module.DeploymentConfig.GitRepoBranch)
	}

	if failedMonitorResult != nil {
		return fmt.Sprintf("%s failed a monitor check with message: %s", prefix, failedMonitorResult.Message), nil
	}

	switch status {
	case models.ModuleRunStatusCompleted:
		return fmt.Sprintf("%s ran successfully", prefix), nil
	case models.ModuleRunStatusFailed:
		return fmt.Sprintf("%s failed", prefix), nil
	case models.ModuleRunStatusInProgress:
		return fmt.Sprintf("%s is in progress", prefix), nil
	case models.ModuleRunStatusQueued:
		return fmt.Sprintf("%s is queued", prefix), nil
	}

	return "", nil
}

func generateDestroyRunDescription(
	config *server.Config,
	module *models.Module,
	run *models.ModuleRun,
	status models.ModuleRunStatus,
	failedMonitorResult *models.ModuleMonitorResult,
) (string, error) {
	prefix := "Destroy"

	if run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindVCS {
		prefix = fmt.Sprintf("Destroy for branch %s", module.DeploymentConfig.GitRepoBranch)
	}

	if failedMonitorResult != nil {
		return fmt.Sprintf("%s failed a monitor check with message: %s", prefix, failedMonitorResult.Message), nil
	}

	switch status {
	case models.ModuleRunStatusCompleted:
		return fmt.Sprintf("%s ran successfully", prefix), nil
	case models.ModuleRunStatusFailed:
		return fmt.Sprintf("%s failed", prefix), nil
	case models.ModuleRunStatusInProgress:
		return fmt.Sprintf("%s is in progress", prefix), nil
	case models.ModuleRunStatusQueued:
		return fmt.Sprintf("%s is queued", prefix), nil
	}

	return "", nil
}

func generateInitRunDescription(
	config *server.Config,
	module *models.Module,
	run *models.ModuleRun,
	status models.ModuleRunStatus,
) (string, error) {
	prefix := "Init"

	switch status {
	case models.ModuleRunStatusCompleted:
		return fmt.Sprintf("%s ran successfully", prefix), nil
	case models.ModuleRunStatusFailed:
		return fmt.Sprintf("%s failed", prefix), nil
	case models.ModuleRunStatusInProgress:
		return fmt.Sprintf("%s is in progress", prefix), nil
	case models.ModuleRunStatusQueued:
		return fmt.Sprintf("%s is queued", prefix), nil
	}

	return "", nil
}

func generateMonitorRunDescription(
	config *server.Config,
	module *models.Module,
	run *models.ModuleRun,
	status models.ModuleRunStatus,
	failedMonitorResult *models.ModuleMonitorResult,
) (string, error) {
	prefix := "Monitor"

	if failedMonitorResult != nil {
		return fmt.Sprintf("%s failed a monitor check with message: %s", prefix, failedMonitorResult.Message), nil
	}

	switch status {
	case models.ModuleRunStatusCompleted:
		return fmt.Sprintf("%s ran successfully", prefix), nil
	case models.ModuleRunStatusFailed:
		return fmt.Sprintf("%s failed", prefix), nil
	case models.ModuleRunStatusInProgress:
		return fmt.Sprintf("%s is in progress", prefix), nil
	case models.ModuleRunStatusQueued:
		return fmt.Sprintf("%s is queued", prefix), nil
	}

	return "", nil
}
