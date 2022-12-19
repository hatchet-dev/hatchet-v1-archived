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

type UserOrgListHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewUserOrgListHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &UserOrgListHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *UserOrgListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	req := &types.ListUserOrgsRequest{}

	if ok := u.DecodeAndValidate(w, r, req); !ok {
		return
	}

	orgs, paginate, err := u.Repo().Org().ListOrgsByUserID(
		user.ID,
		repository.WithPage(req.PaginationRequest),
	)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListUserOrgsResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.Organization, 0),
	}

	for _, org := range orgs {
		resp.Rows = append(resp.Rows, org.ToAPIType())
	}

	u.WriteResult(w, r, resp)
}
