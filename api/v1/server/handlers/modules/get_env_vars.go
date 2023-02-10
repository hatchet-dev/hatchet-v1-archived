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

type ModuleEnvVarsGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleEnvVarsGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleEnvVarsGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleEnvVarsGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mev, _ := r.Context().Value(types.ModuleEnvVarScope).(*models.ModuleEnvVarsVersion)

	res, err := mev.ToAPIType(m.Config().DB.GetEncryptionKey())

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	m.WriteResult(w, r, res)
}
