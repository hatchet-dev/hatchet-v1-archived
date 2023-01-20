package users

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

type UserTeamListHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewUserTeamListHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &UserTeamListHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *UserTeamListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	req := &types.ListUserTeamsRequest{}

	if ok := u.DecodeAndValidate(w, r, req); !ok {
		return
	}

	teams, paginate, err := u.Repo().Team().ListTeamsByUserID(
		user.ID,
		req.OrganizationID,
		repository.WithPage(req.PaginationRequest),
	)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListUserTeamsResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.Team, 0),
	}

	for _, team := range teams {
		resp.Rows = append(resp.Rows, team.ToAPIType())
	}

	u.WriteResult(w, r, resp)
}
