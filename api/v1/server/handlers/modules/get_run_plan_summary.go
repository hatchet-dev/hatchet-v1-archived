package modules

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/terraform"
)

type ModuleGetPlanSummaryHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleGetPlanSummaryHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleGetPlanSummaryHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleGetPlanSummaryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	var path string

	// if this is an apply, get the corresponding plan by the SHA, if it exists
	if run.Kind == models.ModuleRunKindApply {
		if run.ModuleRunConfig.GitCommitSHA == "" {
			m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
				Code:        types.ErrCodeBadRequest,
				Description: "cannot request plan summary for run that doesn't have a commit SHA",
			}, http.StatusBadRequest))

			return
		}

		// find the corresponding run for that SHA
		planKind := models.ModuleRunKindPlan
		planRuns, err := m.Repo().Module().ListModuleRunsByGithubSHA(module.ID, run.ModuleRunConfig.GitCommitSHA, &planKind)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		} else if planRuns == nil || len(planRuns) == 0 {
			m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
				Code:        types.ErrCodeNotFound,
				Description: "a corresponding plan was not found",
			}, http.StatusNotFound))

			return
		}

		path = filestorage.GetPlanJSONPath(module.TeamID, module.ID, planRuns[0].ID)
	} else if run.Kind == models.ModuleRunKindPlan {
		path = filestorage.GetPlanJSONPath(module.TeamID, module.ID, run.ID)
	} else {
		m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: "plan summaries not valid for this plan type",
		}, http.StatusBadRequest))

		return
	}

	jsonBytes, err := m.Config().DefaultFileStore.ReadFile(path, true)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	planSummary, err := terraform.GetPlanSummaryFromBytes(jsonBytes)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	m.WriteResult(w, r, planSummary)
}
