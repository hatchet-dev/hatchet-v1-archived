package modules

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type RunCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewRunCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &RunCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *RunCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)

	req := &types.CreateModuleRunRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	var run *models.ModuleRun
	var reqErr apierrors.RequestError

	if module.DeploymentMechanism == models.DeploymentMechanismGithub {
		run, reqErr = createManualRun(m.Config(), module, models.ModuleRunKind(req.Kind))
	} else if module.DeploymentMechanism == models.DeploymentMechanismLocal {
		run, reqErr = createLocalRun(m.Config(), module, models.ModuleRunKind(req.Kind), req.Hostname)
	}

	if reqErr != nil {
		m.HandleAPIError(w, r, reqErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
	m.WriteResult(w, r, run.ToAPIType(nil))
}
