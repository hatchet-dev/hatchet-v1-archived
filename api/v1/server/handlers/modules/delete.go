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
	"github.com/hatchet-dev/hatchet/internal/terraform"
)

type ModuleDeleteHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleDeleteHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleDeleteHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)

	// if module is locally deployed, check that the state is empty before allowing destroy
	if module.DeploymentMechanism == models.DeploymentMechanismLocal {
		statePath := terraform.GetStatePath(module.TeamID, module.ID)

		fileBytes, err := m.Config().DefaultFileStore.ReadFile(statePath, true)

		if err != nil {
			state, err := terraform.GetInternalStateFromBytes(fileBytes)

			if err != nil {
				m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
				return
			}

			if len(state.Resources) > 0 {
				m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
					Code:        types.ErrCodeBadRequest,
					Description: fmt.Sprintf("Cannot delete locally deployed module as it contains non-deleted state. Please run [hatchet destroy] first."),
				}, http.StatusBadRequest))

				return
			}
		}

		module, err = m.Repo().Module().DeleteModule(module)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}
	} else {
		_, reqErr := createManualRun(m.Config(), module, models.ModuleRunKindDestroy)

		if reqErr != nil {
			m.HandleAPIError(w, r, reqErr)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
	m.WriteResult(w, r, module.ToAPIType())
}
