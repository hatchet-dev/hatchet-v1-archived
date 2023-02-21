package runutils

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

func GenerateRunDescription(config *server.Config, module *models.Module, run *models.ModuleRun, status models.ModuleRunStatus) (string, error) {
	switch run.Kind {
	case models.ModuleRunKindPlan:
		return generatePlanRunDescription(config, module, run, status)
	case models.ModuleRunKindApply:
		return generateApplyRunDescription(config, module, run, status)
	case models.ModuleRunKindDestroy:
		return generateDestroyRunDescription(config, module, run, status)
	}

	return "", fmt.Errorf("unknown run kind %s", run.Kind)
}

func generatePlanRunDescription(config *server.Config, module *models.Module, run *models.ModuleRun, status models.ModuleRunStatus) (string, error) {
	prefix := "Plan"

	if run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindGithub {
		pr, err := getPullRequestFromModuleRun(config, module, run)

		if err != nil {
			return "", err
		}

		prefix = fmt.Sprintf("Plan for pull request %s/%s #%d", pr.GithubRepositoryOwner, pr.GithubRepositoryName, pr.GithubPullRequestNumber)
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

func generateApplyRunDescription(config *server.Config, module *models.Module, run *models.ModuleRun, status models.ModuleRunStatus) (string, error) {
	prefix := "Apply"

	if run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindGithub {
		prefix = fmt.Sprintf("Apply for branch %s", module.DeploymentConfig.GithubRepoBranch)
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

func generateDestroyRunDescription(config *server.Config, module *models.Module, run *models.ModuleRun, status models.ModuleRunStatus) (string, error) {
	prefix := "Destroy"

	if run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindGithub {
		prefix = fmt.Sprintf("Destroy for branch %s", module.DeploymentConfig.GithubRepoBranch)
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

func getPullRequestFromModuleRun(config *server.Config, module *models.Module, run *models.ModuleRun) (*models.GithubPullRequest, error) {
	prComment, err := config.DB.Repository.GithubPullRequest().ReadGithubPullRequestCommentByGithubID(module.ID, run.ModuleRunConfig.GithubCommentID)

	if err != nil {
		return nil, err
	}

	return config.DB.Repository.GithubPullRequest().ReadGithubPullRequestByID(module.TeamID, prComment.GithubPullRequestID)
}
