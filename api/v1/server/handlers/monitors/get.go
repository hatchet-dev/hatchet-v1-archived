package monitors

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type MonitorGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewMonitorGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &MonitorGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *MonitorGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	monitor, _ := r.Context().Value(types.MonitorScope).(*models.ModuleMonitor)

	m.WriteResult(w, r, monitor.ToAPIType())
}
