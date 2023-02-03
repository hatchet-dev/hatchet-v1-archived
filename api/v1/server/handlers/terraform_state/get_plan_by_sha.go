package terraform_state

import (
	"errors"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type TerraformPlanGetBySHAHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTerraformPlanGetBySHAHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TerraformPlanGetBySHAHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TerraformPlanGetBySHAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	if run.ModuleRunConfig.GithubCommitSHA == "" {
		t.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: "cannot request plan for run that doesn't have a commit SHA",
		}, http.StatusBadRequest))

		return
	}

	// find the corresponding run for that SHA
	planKind := models.ModuleRunKindPlan
	planRuns, err := t.Repo().Module().ListModuleRunsByGithubSHA(module.ID, run.ModuleRunConfig.GithubCommitSHA, &planKind)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	} else if planRuns == nil || len(planRuns) == 0 {
		t.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeNotFound,
			Description: "a corresponding plan was not found",
		}, http.StatusNotFound))

		return
	}

	path := getPlanZIPPath(team.ID, module.ID, planRuns[0].ID)

	fileBytes, err := t.Config().DefaultFileStore.ReadFile(path, true)

	if err != nil {
		if errors.Is(err, filestorage.FileDoesNotExist) {
			t.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
				Code:        types.ErrCodeNotFound,
				Description: "the requested terraform state was not found",
			}, http.StatusNotFound))

			return
		}

		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	if _, err = w.Write(fileBytes); err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}
}
