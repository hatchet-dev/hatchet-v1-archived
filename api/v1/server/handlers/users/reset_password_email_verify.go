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

type ResetPasswordEmailVerifyHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewResetPasswordEmailVerifyHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ResetPasswordEmailVerifyHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *ResetPasswordEmailVerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &types.ResetPasswordEmailVerifyTokenRequest{}

	if ok := u.DecodeAndValidate(w, r, request); !ok {
		return
	}

	// find the token in the database by token id
	token, err := u.Repo().PasswordResetToken().ReadPasswordResetTokenByEmailAndTokenID(request.Email, request.TokenID)

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

	// if we've passed all checks, return 200-level response code
	w.WriteHeader(http.StatusOK)
}
