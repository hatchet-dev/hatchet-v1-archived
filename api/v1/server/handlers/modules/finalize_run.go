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

	coreRunFailed := run.Status == models.ModuleRunStatusFailed || (req.ReportKind == types.ModuleRunReportKindCore && req.Status == types.ModuleRunStatusFailed)
	coreRunSucceeded := run.Status == models.ModuleRunStatusCompleted || (req.ReportKind == types.ModuleRunReportKindCore && req.Status == types.ModuleRunStatusCompleted)
	monitorsInProgress := !coreRunFailed && len(run.ModuleMonitorResults) != len(run.Monitors)

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
	if run.Kind == models.ModuleRunKindPlan && run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindVCS {
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

			monitorResultAPIType := failedMonitorResult.ToAPIType()

			commentBody = getFailedMonitorGithubComment(m.Config(), module, monitor, run, monitorResultAPIType.Message)

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

			commentBody = getSuccessfulPlanGithubComment(m.Config(), module, run, fileBytes)

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
			commentBody = getFailedPlanGithubComment(m.Config(), module, run, run.StatusDescription)

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

func getSuccessfulPlanGithubComment(config *server.Config, module *models.Module, run *models.ModuleRun, planBytes []byte) string {
	shortSHA := getShortSHA(run)
	shaLink := getCommitSHALink(module, run)
	moduleLink := getModuleLink(config, module.TeamID, module.ID)

	return fmt.Sprintf(
		"## Hatchet Plan\n"+
			"Successfully created a plan from commit [`%s`](%s). Details:\n"+
			"||Run information|\n"+
			"|-|-|\n"+
			"| Module | [`%s`](%s) |\n"+
			"| Commit SHA | [`%s`](%s) |\n"+
			"```\n%s\n```",
		shortSHA, shaLink, module.Name, moduleLink, shortSHA, shaLink, string(planBytes),
	)
}

func getFailedMonitorGithubComment(config *server.Config, module *models.Module, monitor *models.ModuleMonitor, run *models.ModuleRun, msg string) string {
	shortSHA := getShortSHA(run)
	shaLink := getCommitSHALink(module, run)
	moduleLink := getModuleLink(config, module.TeamID, module.ID)
	monitorLink := getMonitorLink(config, module.TeamID, monitor.ID)

	return fmt.Sprintf(
		"## Hatchet Plan\n"+
			"Could not create a plan from commit [`%s`](%s). Details:\n"+
			"||Run information|\n"+
			"|-|-|\n"+
			"| Module | [`%s`](%s) |\n"+
			"| Monitor | [`%s`](%s) |\n"+
			"| Commit SHA | [`%s`](%s) |\n"+
			"| Error message | %s |\n",
		shortSHA, shaLink, module.Name, moduleLink, monitor.DisplayName, monitorLink, shortSHA, shaLink, msg,
	)
}

func getFailedPlanGithubComment(config *server.Config, module *models.Module, run *models.ModuleRun, msg string) string {
	shortSHA := getShortSHA(run)
	shaLink := getCommitSHALink(module, run)
	moduleLink := getModuleLink(config, module.TeamID, module.ID)

	return fmt.Sprintf(
		"## Hatchet Plan\n"+
			"Could not create a plan from commit [`%s`](%s). Details:\n"+
			"||Run information|\n"+
			"|-|-|\n"+
			"| Module | [`%s`](%s) |\n"+
			"| Commit SHA | [`%s`](%s) |\n"+
			"| Error message | %s |\n",
		shortSHA, shaLink, module.Name, moduleLink, shortSHA, shaLink, msg,
	)
}

func getModuleLink(config *server.Config, teamID, moduleID string) string {
	return fmt.Sprintf("%s/teams/%s/modules/%s", config.ServerRuntimeConfig.ServerURL, teamID, moduleID)
}

func getMonitorLink(config *server.Config, teamID, monitorID string) string {
	return fmt.Sprintf("%s/teams/%s/monitors/%s", config.ServerRuntimeConfig.ServerURL, teamID, monitorID)
}

func getShortSHA(run *models.ModuleRun) string {
	return run.ModuleRunConfig.GitCommitSHA[0:7]
}

func getCommitSHALink(module *models.Module, run *models.ModuleRun) string {
	return fmt.Sprintf(
		"https://github.com/%s/%s/commit/%s",
		module.DeploymentConfig.GithubRepoOwner,
		module.DeploymentConfig.GithubRepoName,
		run.ModuleRunConfig.GitCommitSHA,
	)
}
