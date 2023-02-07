package modules

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type ModuleRunFinalizeHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleRunFinalizeHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleRunFinalizeHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleRunFinalizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	req := &types.FinalizeModuleRunRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	switch req.Status {
	case types.ModuleRunStatusFailed:
		run.Status = models.ModuleRunStatusFailed
	case types.ModuleRunStatusCompleted:
		run.Status = models.ModuleRunStatusCompleted
	}

	if req.Description == "" {
		desc, err := generateRunDescription(m.Config(), module, run, run.Status)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		run.StatusDescription = desc
	} else {
		run.StatusDescription = req.Description
	}

	run, err := m.Repo().Module().UpdateModuleRun(run)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// if this is a successful apply, clear the lock from the module
	if run.Kind == models.ModuleRunKindApply && run.Status == models.ModuleRunStatusCompleted {
		module.LockID = ""
		module.LockKind = models.ModuleLockKind("")

		module, err = m.Repo().Module().UpdateModule(module)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}
	}

	m.WriteResult(w, r, run.ToAPIType())
}

func generateRunDescription(config *server.Config, module *models.Module, run *models.ModuleRun, status models.ModuleRunStatus) (string, error) {
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
