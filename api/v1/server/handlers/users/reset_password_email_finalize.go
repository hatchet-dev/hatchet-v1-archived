package users

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

type ResetPasswordEmailFinalizeHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewResetPasswordEmailFinalizeHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ResetPasswordEmailFinalizeHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *ResetPasswordEmailFinalizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &types.ResetPasswordEmailFinalizeRequest{}

	if ok := u.DecodeAndValidate(w, r, request); !ok {
		return
	}

	// find the token in the database by token id
	token, err := u.Repo().PasswordResetToken().ReadPasswordResetTokenByEmailAndTokenID(request.Email, request.TokenID)

	if err != nil {
		// TODO(abelanger5): error handling on token not found -- throw generic error
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

	user, err := u.Repo().User().ReadUserByEmail(request.Email)

	if err != nil {
		// TODO(abelanger5): error handling user not read
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// TODO(abelanger5): make sure password doesn't match user's old password

	// update the user's password
	user.Password = request.NewPassword

	if err := user.HashPassword(); err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	user, err = u.Repo().User().UpdateUser(user)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	token.Revoked = true

	token, err = u.Repo().PasswordResetToken().UpdatePasswordResetToken(token)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
