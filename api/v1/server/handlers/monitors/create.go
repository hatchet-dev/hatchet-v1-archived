package monitors

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"
)

type MonitorCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewMonitorCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &MonitorCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *MonitorCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	// TODO: cron schedule validation
	req := &types.CreateMonitorRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	monitor := &models.ModuleMonitor{
		TeamID:           team.ID,
		Kind:             models.MonitorKindState,
		PresetPolicyName: models.ModuleMonitorPresetPolicyNameDrift,
		PolicyBytes:      req.PolicyBytes,
	}

	monitor, err := m.Repo().ModuleMonitor().CreateModuleMonitor(monitor)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	err = dispatcher.DispatchCronMonitor(m.Config().TemporalClient, team.ID, monitor.ID, req.CronSchedule)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	w.WriteHeader(http.StatusOK)
}
