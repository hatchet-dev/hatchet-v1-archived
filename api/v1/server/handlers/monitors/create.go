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

	req := &types.CreateMonitorRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	targetModules, reqErr := getMonitorModulesFromRequest(m.Config(), team, req.Modules)

	if reqErr != nil {
		m.HandleAPIError(w, r, reqErr)
		return
	}

	// load the query to make sure it parses
	_, err := opa.LoadQueryFromBytes(opa.PACKAGE_HATCHET_MODULE, []byte(req.PolicyBytes))

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: fmt.Sprintf("Could not parse policy: %s", err.Error()),
		}, http.StatusBadRequest))

		return
	}

	monitor := &models.ModuleMonitor{
		TeamID:       team.ID,
		DisplayName:  req.Name,
		Description:  req.Description,
		Kind:         models.ModuleMonitorKind(req.Kind),
		CronSchedule: req.CronSchedule,
		Modules:      targetModules,
		CurrentMonitorPolicyBytesVersion: models.MonitorPolicyBytesVersion{
			Version:     1,
			PolicyBytes: []byte(req.PolicyBytes),
		},
	}

	monitor, err = m.Repo().ModuleMonitor().CreateModuleMonitor(monitor)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	if req.Kind == types.MonitorKindPlan || req.Kind == types.MonitorKindState {
		err = dispatcher.DispatchCronMonitor(m.Config().TemporalClient, team.ID, monitor.ID, req.CronSchedule)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

			return
		}
	}

	m.WriteResult(w, r, monitor.ToAPIType())
}
