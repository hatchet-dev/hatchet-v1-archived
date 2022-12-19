package pats

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

type PATListHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewPATListHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &PATListHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *PATListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	req := &types.ListPATsRequest{}

	if ok := u.DecodeAndValidate(w, r, req); !ok {
		return
	}

	pats, paginate, err := u.Repo().PersonalAccessToken().ListPersonalAccessTokensByUserID(
		user.ID,
		repository.WithPage(req.PaginationRequest),
	)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListPATsResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.PersonalAccessToken, 0),
	}

	for _, pat := range pats {
		resp.Rows = append(resp.Rows, pat.ToAPIType())
	}

	u.WriteResult(w, r, resp)
}
