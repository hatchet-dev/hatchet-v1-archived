package monitors

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"
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
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	monitor, _ := r.Context().Value(types.MonitorScope).(*models.ModuleMonitor)

	if monitor.IsDefault {
		m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: fmt.Sprintf("Cannot delete default modules"),
		}, http.StatusBadRequest))

		return
	}

	monitor, err := m.Repo().ModuleMonitor().DeleteModuleMonitor(monitor)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	if monitor.IsCronKind() {
		// terminate the workflow
		err = dispatcher.DeleteCronMonitor(m.Config().TemporalClient, team.ID, monitor.ID)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

			return
		}
	}

	m.WriteResult(w, r, monitor.ToAPIType())
}
