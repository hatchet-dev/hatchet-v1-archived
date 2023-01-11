package users

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type VerifyEmailHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewVerifyEmailHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &VerifyEmailHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *VerifyEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	request := &types.VerifyEmailRequest{}

	if ok := u.DecodeAndValidate(w, r, request); !ok {
		return
	}

	// find the token in the database by token id
	token, err := u.Repo().VerifyEmailToken().ReadVerifyEmailTokenByEmailAndTokenID(user.Email, request.TokenID)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrForbidden(
			fmt.Errorf("token was not found"),
		))

		return
	}

	// verify that the tokens match the hashed
	if verified, err := token.VerifyToken(request.Token); !verified || err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrForbidden(
			fmt.Errorf("token did not match hashed value"),
		))

		return
	}

	// verify that the token is valid
	if token.IsExpired() {
		u.HandleAPIError(w, r, apierrors.NewErrForbidden(
			fmt.Errorf("token was expired"),
		))

		return
	}

	// verify that the token has not been used
	if token.Revoked {
		u.HandleAPIError(w, r, apierrors.NewErrForbidden(
			fmt.Errorf("token was revoked (already used)"),
		))

		return
	}

	user.EmailVerified = true

	user, err = u.Repo().User().UpdateUser(user)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	token.Revoked = true

	token, err = u.Repo().VerifyEmailToken().UpdateVerifyEmailToken(token)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
