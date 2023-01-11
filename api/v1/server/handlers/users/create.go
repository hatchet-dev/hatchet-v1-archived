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
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type UserCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewUserCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &UserCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *UserCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &types.CreateUserRequest{}

	if ok := u.DecodeAndValidate(w, r, request); !ok {
		return
	}

	if !u.Config().AuthConfig.IsEmailAllowed(request.Email) {
		u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.InvalidEmailOrPassword, http.StatusUnauthorized, "email is not in restricted domain list"))
		return
	}

	// determine if the user exists before attempting to write the user
	existingUser, err := u.Repo().User().ReadUserByEmail(request.Email)

	if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
		// if this is not a "Not-Found" error, this is an internal error as we cannot read from the database
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	if existingUser != nil {
		u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.InvalidEmailOrPassword, http.StatusUnauthorized, "user already exists"))
		return
	}

	user := &models.User{
		DisplayName: request.DisplayName,
		Email:       request.Email,
		Password:    request.Password,
	}

	// write the user to the db
	user, err = u.Repo().User().CreateUser(user)

	if err != nil {
		if errors.Is(err, repository.RepositoryUniqueConstraintFailed) {
			u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.InvalidEmailOrPassword, http.StatusUnauthorized, "unique constraint failed"))
			return
		}

		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// save the user as authenticated in the session
	_, err = authn.SaveUserAuthenticated(w, r, u.Config(), user)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	u.WriteResult(w, r, user.ToAPIType())

	// send a verification email after creation
	if reqErr := sendVerifyEmail(u.Config(), user); reqErr != nil {
		// we've already written a success message, so don't rewrite it
		u.HandleAPIErrorNoWrite(w, r, reqErr)
	}
}
