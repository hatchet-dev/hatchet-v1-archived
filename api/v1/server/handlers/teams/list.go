package teams

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

type TeamListHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTeamListHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TeamListHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TeamListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	org, _ := r.Context().Value(types.OrgScope).(*models.Organization)

	req := &types.ListTeamsRequest{}

	if ok := t.DecodeAndValidate(w, r, req); !ok {
		return
	}

	teams, paginate, err := t.Repo().Team().ListTeamsByOrgID(
		org.ID,
		repository.WithPage(req.PaginationRequest),
	)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListTeamsResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.Team, 0),
	}

	for _, team := range teams {
		resp.Rows = append(resp.Rows, team.ToAPIType())
	}

	t.WriteResult(w, r, resp)
}
