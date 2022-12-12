package users

import (
	"errors"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewUserCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) *UserCreateHandler {
	return &UserCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *UserCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &types.CreateUserRequest{}

	if ok := u.DecodeAndValidate(w, r, request); !ok {
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

	// hash the password using bcrypt
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	user.Password = string(hashedPw)

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

	w.WriteHeader(http.StatusCreated)
	u.WriteResult(w, r, user.ToAPIType())
}
