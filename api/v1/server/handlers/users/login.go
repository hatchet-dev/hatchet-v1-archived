package users

import (
	"errors"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/authn"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type UserLoginHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewUserLoginHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &UserLoginHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *UserLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &types.LoginUserRequest{}

	if ok := u.DecodeAndValidate(w, r, request); !ok {
		return
	}

	if !u.Config().AuthConfig.IsEmailAllowed(request.Email) {
		u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.InvalidEmailOrPassword, http.StatusUnauthorized, "email is not in restricted domain list"))
		return
	}

	// determine if the user exists
	existingUser, err := u.Repo().User().ReadUserByEmail(request.Email)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.InvalidEmailOrPassword, http.StatusUnauthorized, "user does not exists"))
			return
		}

		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	if verified, err := existingUser.VerifyPassword(request.Password); !verified || err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.InvalidEmailOrPassword, http.StatusUnauthorized, "bad password"))
		return
	}

	// save the user as authenticated in the session
	_, err = authn.SaveUserAuthenticated(w, r, u.Config(), existingUser)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	u.WriteResult(w, r, existingUser.ToAPIType())
}
