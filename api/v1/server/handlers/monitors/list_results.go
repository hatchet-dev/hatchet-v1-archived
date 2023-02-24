package monitors

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type MonitorResultListHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewMonitorResultListHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &MonitorResultListHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *MonitorResultListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	req := &types.ListMonitorResultsRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	filter := &repository.ListModuleMonitorResultsOpts{
		ModuleID:        req.ModuleID,
		ModuleMonitorID: req.ModuleMonitorID,
	}

	results, paginate, err := m.Repo().ModuleMonitor().ListModuleMonitorResults(
		team.ID,
		filter,
		repository.WithPage(req.PaginationRequest),
	)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListMonitorResultsResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.ModuleMonitorResult, 0),
	}

	for _, result := range results {
		resp.Rows = append(resp.Rows, result.ToAPIType())
	}

	m.WriteResult(w, r, resp)
}
