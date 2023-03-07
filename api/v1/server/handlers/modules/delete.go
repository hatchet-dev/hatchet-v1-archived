package modules

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
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

	_, reqErr := createManualRun(m.Config(), module, models.ModuleRunKindDestroy)

	if reqErr != nil {
		m.HandleAPIError(w, r, reqErr)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	m.WriteResult(w, r, module.ToAPIType())
}
