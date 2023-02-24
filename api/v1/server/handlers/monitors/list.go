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

type MonitorListHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewMonitorListHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &MonitorListHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *MonitorListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	req := &types.ListMonitorsRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	monitors, paginate, err := m.Repo().ModuleMonitor().ListModuleMonitorsByTeamID(
		team.ID,
		repository.WithPage(req.PaginationRequest),
	)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListMonitorsResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.ModuleMonitorMeta, 0),
	}

	for _, monitor := range monitors {
		resp.Rows = append(resp.Rows, monitor.ToAPITypeMeta())
	}

	m.WriteResult(w, r, resp)
}
