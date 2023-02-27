package monitors

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type MonitorDeleteHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewMonitorDeleteHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &MonitorDeleteHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *MonitorDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	monitor, _ := r.Context().Value(types.MonitorScope).(*models.ModuleMonitor)

	monitor, err := m.Repo().ModuleMonitor().DeleteModuleMonitor(monitor)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	m.WriteResult(w, r, monitor.ToAPIType())
}
