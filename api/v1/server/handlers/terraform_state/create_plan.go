package terraform_state

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
	"github.com/hatchet-dev/hatchet/internal/models"
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

	jsonPlanPath := filestorage.GetPlanJSONPath(team.ID, module.ID, run.ID)
	prettyPlanPath := filestorage.GetPlanPrettyPath(team.ID, module.ID, run.ID)

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
}
