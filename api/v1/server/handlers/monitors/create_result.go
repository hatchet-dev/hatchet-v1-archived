package monitors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/runutils"
)

type MonitorResultCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewMonitorResultCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &MonitorResultCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *MonitorResultCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	req := &types.CreateMonitorResultRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	// read the monitor ID and check that it belongs to this team
	monitor, err := m.Repo().ModuleMonitor().ReadModuleMonitorByID(team.ID, req.MonitorID)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrForbidden(fmt.Errorf("monitor %s does not belong to team %s", req.MonitorID, team.ID)))

		return
	}

	// create the result in the database, and associate it with this module run
	result := &models.ModuleMonitorResult{
		TeamID:          team.ID,
		ModuleID:        module.ID,
		ModuleRunID:     run.ID,
		ModuleMonitorID: monitor.ID,
		Status:          models.MonitorResultStatus(req.Status),
		Title:           req.Title,
		Severity:        models.MonitorResultSeverity(req.Severity),
	}

	if req.SuccessMessage != "" {
		result.Message = req.SuccessMessage
	} else if len(req.FailureMessages) > 0 {
		result.Message = strings.Join(req.FailureMessages, ",")
	}

	result, err = m.Repo().ModuleMonitor().CreateModuleMonitorResult(monitor, result)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	// if the module run has a "monitor" kind, update the module run as well
	if run.Kind == models.ModuleRunKindMonitor {
		run.Status = models.ModuleRunStatusCompleted

		desc, err := runutils.GenerateRunDescription(m.Config(), module, run, run.Status)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		run.StatusDescription = desc

		run, err = m.Repo().Module().UpdateModuleRun(run)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
