package users

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/authn"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type UserDeleteCurrentHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewUserDeleteCurrentHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &UserDeleteCurrentHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *UserDeleteCurrentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	// check that the user is not an owner in any organizations. if they are, they must change
	// owners or delete those organizations
	orgs, _, err := u.Repo().Org().ListOrgsByUserID(user.ID)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	ownerOrgs := make([]string, 0)

	for _, org := range orgs {
		if org.OwnerID == user.ID {
			ownerOrgs = append(ownerOrgs, org.ID)
		}
	}

	if len(ownerOrgs) > 0 {
		u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Description: fmt.Sprintf("the current user cannot be deleted because they are an owner of the following organizations: %s", strings.Join(ownerOrgs, ", ")),
			Code:        types.ErrCodeBadRequest,
		}, http.StatusBadRequest))

		return
	}

	user, err = u.Repo().User().DeleteUser(user)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	if err := authn.SaveUserUnauthenticated(w, r, u.Config()); err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
	}

	w.WriteHeader(http.StatusAccepted)
}
