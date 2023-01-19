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

type TeamListMemberHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTeamListMemberHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TeamListMemberHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TeamListMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	req := &types.ListTeamMembersRequest{}

	if ok := t.DecodeAndValidate(w, r, req); !ok {
		return
	}

	teamMembers, paginate, err := t.Repo().Team().ListTeamMembersByTeamID(
		team.ID,
		repository.WithPage(req.PaginationRequest),
	)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListTeamMembersResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.TeamMember, 0),
	}

	for _, teamMember := range teamMembers {
		resp.Rows = append(resp.Rows, teamMember.ToAPIType())
	}

	t.WriteResult(w, r, resp)
}
