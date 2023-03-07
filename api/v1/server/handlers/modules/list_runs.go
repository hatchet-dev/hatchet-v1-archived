package modules

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

type ModuleRunsListHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleRunsListHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleRunsListHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleRunsListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)

	req := &types.ListModuleRunsRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	var statuses []string

	if req.Status != nil {
		for _, reqStatus := range req.Status {
			statuses = append(statuses, string(reqStatus))
		}
	}

	var kinds []string

	if req.Kind != nil {
		for _, reqKind := range req.Kind {
			kinds = append(kinds, string(reqKind))
		}
	}

	modRuns, paginate, err := m.Repo().Module().ListRunsByModuleID(
		module.ID,
		&repository.ModuleRunsFilter{
			Status: statuses,
			Kind:   kinds,
		},
		repository.WithPage(req.PaginationRequest),
		repository.WithSortBy("created_at"),
		repository.WithOrder(repository.OrderDesc),
	)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListModuleRunsResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.ModuleRunOverview, 0),
	}

	for _, modRun := range modRuns {
		resp.Rows = append(resp.Rows, modRun.ToAPITypeOverview())
	}

	m.WriteResult(w, r, resp)
}
