package terraform_state

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/git/github"
	"github.com/hatchet-dev/hatchet/internal/models"

	githubsdk "github.com/google/go-github/v49/github"
)

type TerraformPlanCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTerraformPlanCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TerraformPlanCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TerraformPlanCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	req := &types.CreateTerraformPlanRequest{}

	if ok := t.DecodeAndValidate(w, r, req); !ok {
		return
	}

	jsonPlanPath := getPlanJSONPath(team.ID, module.ID, run.ID)
	prettyPlanPath := getPlanPrettyPath(team.ID, module.ID, run.ID)

	err := t.Config().DefaultFileStore.WriteFile(jsonPlanPath, []byte(req.PlanJSON), true)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	err = t.Config().DefaultFileStore.WriteFile(prettyPlanPath, []byte(req.PlanPretty), true)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	t.WriteResult(w, r, nil)

	// TODO: update module run status
	if run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindGithub {
		client, err := github.GetGithubAppClientFromModule(t.Config(), module)

		if err != nil {
			t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

			return
		}

		commentBody := "## Hatchet Plan\n"

		commentBody += fmt.Sprintf("```\n%s\n```", req.PlanPretty)

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
			t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

			return
		}

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
			t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

			return
		}
	}
}

func getPlanJSONPath(teamID, moduleID, runID string) string {
	return fmt.Sprintf("%s/%s/%s/plan.json", teamID, moduleID, runID)
}

func getPlanPrettyPath(teamID, moduleID, runID string) string {
	return fmt.Sprintf("%s/%s/%s/plan.txt", teamID, moduleID, runID)
}
