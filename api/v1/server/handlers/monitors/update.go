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
	"github.com/hatchet-dev/hatchet/internal/opa"
)

type MonitorUpdateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewMonitorUpdateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &MonitorUpdateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *MonitorUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	monitor, _ := r.Context().Value(types.MonitorScope).(*models.ModuleMonitor)

	req := &types.UpdateMonitorRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	if req.Modules != nil {
		targetModules, reqErr := getMonitorModulesFromRequest(m.Config(), team, req.Modules)

		if reqErr != nil {
			m.HandleAPIError(w, r, reqErr)
			return
		}

		monitor.Modules = targetModules
	}

	if req.CronSchedule != "" {
		monitor.CronSchedule = req.CronSchedule
	}

	if req.Description != "" {
		monitor.Description = req.Description
	}

	if req.Name != "" {
		monitor.DisplayName = req.Name
	}

	if req.Kind != "" {
		monitor.Kind = models.ModuleMonitorKind(req.Kind)
	}

	if req.PolicyBytes != "" {
		_, err := opa.LoadQueryFromBytes(req.Name, []byte(req.PolicyBytes))

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
				Code:        types.ErrCodeBadRequest,
				Description: fmt.Sprintf("Could not parse policy: %s", err.Error()),
			}, http.StatusBadRequest))

			return
		}

		monitor.PolicyBytes = []byte(req.PolicyBytes)
	}

	monitor, err := m.Repo().ModuleMonitor().UpdateModuleMonitor(monitor)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	m.WriteResult(w, r, monitor.ToAPIType())
}
