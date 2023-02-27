package monitors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/opa"
	"github.com/hatchet-dev/hatchet/internal/repository"
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

	var targetModuleIDs []string

	if req.Modules == nil || len(req.Modules) == 0 {
		targetModuleIDs = make([]string, 0)
	} else {
		targetModuleIDs = req.Modules
	}

	var modErr error
	var targetModules []models.Module = make([]models.Module, 0)

	for _, modID := range targetModuleIDs {
		// ensure all modules are in the team
		mod, err := m.Repo().Module().ReadModuleByID(team.ID, modID)

		if err != nil {
			if errors.Is(err, repository.RepositoryErrorNotFound) {
				m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
					Code:        types.ErrCodeBadRequest,
					Description: fmt.Sprintf("Could not find module with id %s", modID),
				}, http.StatusBadRequest))

				return
			} else {
				modErr = multierror.Append(modErr, err)
			}
		}

		targetModules = append(targetModules, *mod)
	}

	if modErr != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(modErr))

		return
	}

	// load the query to make sure it parses
	_, err := opa.LoadQueryFromBytes(req.Name, []byte(req.PolicyBytes))

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: fmt.Sprintf("Could not parse policy: %s", err.Error()),
		}, http.StatusBadRequest))

		return
	}

	monitor := &models.ModuleMonitor{
		TeamID:           team.ID,
		DisplayName:      req.Name,
		Description:      req.Description,
		Kind:             models.ModuleMonitorKind(req.Kind),
		PresetPolicyName: models.ModuleMonitorPresetPolicyNameDrift,
		CronSchedule:     req.CronSchedule,
		PolicyBytes:      []byte(req.PolicyBytes),
		Modules:          targetModules,
	}

	monitor, err = m.Repo().ModuleMonitor().CreateModuleMonitor(monitor)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	err = dispatcher.DispatchCronMonitor(m.Config().TemporalClient, team.ID, monitor.ID, req.CronSchedule)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	m.WriteResult(w, r, monitor.ToAPIType())
}
