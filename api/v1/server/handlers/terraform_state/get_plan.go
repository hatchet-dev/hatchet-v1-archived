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

type TerraformPlanGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTerraformPlanGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TerraformPlanGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TerraformPlanGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	path := GetPlanPrettyPath(team.ID, module.ID, run.ID)

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
