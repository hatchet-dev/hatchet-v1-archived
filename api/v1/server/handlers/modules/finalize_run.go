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

	monitorsInProgress := len(run.ModuleMonitorResults) != len(run.Monitors)
	coreRunSucceeded := run.Status == models.ModuleRunStatusCompleted || (req.ReportKind == types.ModuleRunReportKindCore && req.Status == types.ModuleRunStatusCompleted)
	var failedMonitorResult *models.ModuleMonitorResult

	for _, monitorResult := range run.ModuleMonitorResults {
		if monitorResult.Status == models.MonitorResultStatusFailed {
			failedMonitorResult = &monitorResult
			break
		}
	}

	// if there are still monitors in progress, we simply do nothing and wait for additional reports from
	// the monitors
	if monitorsInProgress {
		m.WriteResult(w, r, run.ToAPITypeOverview())
		return
	}

	if coreRunSucceeded {
		run.Status = models.ModuleRunStatusCompleted
	} else {
		run.Status = models.ModuleRunStatusFailed
	}

	if req.Description == "" {
		desc, err := runutils.GenerateRunDescription(m.Config(), module, run, run.Status, failedMonitorResult)

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

	// write the github comment if all monitors have succeeded for the plan, or one monitor has failed
	if run.Kind == models.ModuleRunKindPlan && run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindGithub {
		client, err := github.GetGithubAppClientFromModule(m.Config(), module)

		if err != nil {
			m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

			return
		}

		var commentBody string

		// if there's a failed monitor, we write that to github
		if failedMonitorResult != nil {
			monitor, err := m.Repo().ModuleMonitor().ReadModuleMonitorByID(module.TeamID, failedMonitorResult.ModuleMonitorID)

			if err != nil {
				m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

				return
			}

			monitorAPIType := monitor.ToAPIType()
			monitorResultAPIType := failedMonitorResult.ToAPIType()

			commentBody = "## Hatchet Plan\n"

			commentBody += fmt.Sprintf("Monitor %s failed for this plan with message: %s", monitorAPIType.Name, monitorResultAPIType.Message)

			_, _, err = client.Checks.UpdateCheckRun(
				context.Background(),
				module.DeploymentConfig.GithubRepoOwner,
				module.DeploymentConfig.GithubRepoName,
				run.ModuleRunConfig.GithubCheckID,
				githubsdk.UpdateCheckRunOptions{
					Name:       fmt.Sprintf("Hatchet plan for %s", module.DeploymentConfig.ModulePath),
					Status:     githubsdk.String("completed"),
					Conclusion: githubsdk.String("failure"),
				},
			)

			if err != nil {
				m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

				return
			}
		} else if run.Status == models.ModuleRunStatusCompleted {
			// otherwise, if the run was successful, write the prettified plan to github
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
					Conclusion: githubsdk.String("failure"),
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
