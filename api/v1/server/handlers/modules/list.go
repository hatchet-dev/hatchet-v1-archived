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

type ModuleListHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleListHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleListHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	req := &types.ListModulesRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	mods, paginate, err := m.Repo().Module().ListModulesByTeamID(
		team.ID,
		repository.WithPage(req.PaginationRequest),
	)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListModulesResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.Module, 0),
	}

	for _, mod := range mods {
		resp.Rows = append(resp.Rows, mod.ToAPIType())
	}

	m.WriteResult(w, r, resp)
}
