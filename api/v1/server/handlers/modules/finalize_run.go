package modules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
	"github.com/hatchet-dev/hatchet/internal/integrations/git/github"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/notifications"
	"github.com/hatchet-dev/hatchet/internal/runutils"

	githubsdk "github.com/google/go-github/v49/github"
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
		desc, err := runutils.GenerateRunDescription(m.Config(), module, run, run.Status)

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

	// if this is a successful destroy, delete the module
	if run.Kind == models.ModuleRunKindDestroy && run.Status == models.ModuleRunStatusCompleted {
		module, err = m.Repo().Module().DeleteModule(module)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}
	}

	m.WriteResult(w, r, run.ToAPITypeOverview())

	// write github comment
	if run.Kind == models.ModuleRunKindPlan && run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindGithub {
		client, err := github.GetGithubAppClientFromModule(m.Config(), module)

		if err != nil {
			m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

			return
		}

		var commentBody string

		// if the run was successful, write the prettified plan to github
		if run.Status == models.ModuleRunStatusCompleted {
			fileBytes, err := m.Config().DefaultFileStore.ReadFile(filestorage.GetPlanPrettyPath(module.TeamID, module.ID, run.ID), true)

			if err != nil {
				m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

				return
			}

			commentBody = "## Hatchet Plan\n"

			commentBody += fmt.Sprintf("```\n%s\n```", string(fileBytes))

			_, _, err = client.Checks.UpdateCheckRun(
				context.Background(),
				module.DeploymentConfig.GithubRepoOwner,
				module.DeploymentConfig.GithubRepoName,
				run.ModuleRunConfig.GithubCheckID,
				githubsdk.UpdateCheckRunOptions{
					Name:       fmt.Sprintf("Hatchet plan for %s", module.DeploymentConfig.ModulePath),
					Status:     githubsdk.String("completed"),
					Conclusion: githubsdk.String("success"),
				},
			)

			if err != nil {
				m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

				return
			}
		} else if run.Status == models.ModuleRunStatusFailed {
			// otherwise, write that the module run failed
			commentBody = "## Hatchet Plan\n"

			commentBody += fmt.Sprintf("Plan failed")

			_, _, err = client.Checks.UpdateCheckRun(
				context.Background(),
				module.DeploymentConfig.GithubRepoOwner,
				module.DeploymentConfig.GithubRepoName,
				run.ModuleRunConfig.GithubCheckID,
				githubsdk.UpdateCheckRunOptions{
					Name:       fmt.Sprintf("Hatchet plan for %s", module.DeploymentConfig.ModulePath),
					Status:     githubsdk.String("completed"),
					Conclusion: githubsdk.String("success"),
				},
			)

			if err != nil {
				m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

				return
			}
		}

		_, _, err = client.Issues.EditComment(
			context.Background(),
			module.DeploymentConfig.GithubRepoOwner,
			module.DeploymentConfig.GithubRepoName,
			run.ModuleRunConfig.GithubCommentID,
			&githubsdk.IssueComment{
				Body: &commentBody,
			},
		)

		if err != nil {
			m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

			return
		}
	}

	// add notification if necessary
	err = notifications.CreateNotificationFromModuleRun(m.Config(), module.TeamID, run)

	if err != nil {
		m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}
}
