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

type MonitorPolicyGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewMonitorPolicyGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &MonitorPolicyGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *MonitorPolicyGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	monitor, _ := r.Context().Value(types.MonitorScope).(*models.ModuleMonitor)

	if _, err := w.Write(monitor.CurrentMonitorPolicyBytesVersion.PolicyBytes); err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}
}
