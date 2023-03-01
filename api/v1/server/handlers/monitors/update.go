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
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"
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

	if monitor.IsDefault {
		m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: fmt.Sprintf("Cannot update default modules"),
		}, http.StatusBadRequest))

		return
	}

	req := &types.UpdateMonitorRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	var targetModules []*models.Module

	if req.Modules != nil {
		targetModules = make([]*models.Module, 0)

		gotModules, reqErr := getMonitorModulesFromRequest(m.Config(), team, req.Modules)

		if reqErr != nil {
			m.HandleAPIError(w, r, reqErr)
			return
		}

		for _, gotModule := range gotModules {
			targetModules = append(targetModules, &gotModule)
		}
	}

	if req.Disabled != nil {
		monitor.Disabled = *req.Disabled
	}

	if req.CronSchedule != "" {
		// make sure cron schedule is only set for cron kinds
		if (req.Kind == "" && !monitor.IsCronKind()) || (req.Kind != "" && !models.IsCronKind(models.ModuleMonitorKind(req.Kind))) {
			m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
				Code:        types.ErrCodeBadRequest,
				Description: fmt.Sprintf("Cannot set cron schedule when monitor kind is not a scheduled monitor"),
			}, http.StatusBadRequest))

			return
		}

		monitor.CronSchedule = req.CronSchedule
	}

	if req.Description != "" {
		monitor.Description = req.Description
	}

	if req.Name != "" {
		monitor.DisplayName = req.Name
	}

	var shouldCreateCron bool

	if req.Kind != "" {
		shouldCreateCron = !monitor.IsCronKind() && models.IsCronKind(models.ModuleMonitorKind(req.Kind))

		monitor.Kind = models.ModuleMonitorKind(req.Kind)
	}

	if req.PolicyBytes != "" {
		// create a new version of the policy bytes
		newPolicyBytesVersion := models.MonitorPolicyBytesVersion{
			Version:     monitor.CurrentMonitorPolicyBytesVersion.Version + 1,
			PolicyBytes: []byte(req.PolicyBytes),
		}

		_, err := opa.LoadQueryFromBytes(opa.PACKAGE_HATCHET_MODULE, []byte(req.PolicyBytes))

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
				Code:        types.ErrCodeBadRequest,
				Description: fmt.Sprintf("Could not parse policy: %s", err.Error()),
			}, http.StatusBadRequest))

			return
		}

		monitor.CurrentMonitorPolicyBytesVersion = newPolicyBytesVersion
	}

	monitor, err := m.Repo().ModuleMonitor().UpdateModuleMonitor(monitor)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	if targetModules != nil {
		_monitor, err := m.Repo().ModuleMonitor().ReplaceModuleMonitorModules(monitor, targetModules)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

			return
		}

		monitor.Modules = _monitor.Modules
	}

	if shouldCreateCron {
		err = dispatcher.DispatchCronMonitor(m.Config().TemporalClient, team.ID, monitor.ID, monitor.CronSchedule)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

			return
		}
	} else if req.CronSchedule != "" {
		// update the workflow
		err = dispatcher.UpdateCronMonitor(m.Config().TemporalClient, team.ID, monitor.ID, monitor.CronSchedule)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

			return
		}
	}

	m.WriteResult(w, r, monitor.ToAPIType())
}
