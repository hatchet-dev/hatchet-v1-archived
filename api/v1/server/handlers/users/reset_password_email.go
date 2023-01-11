package users

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/notifier"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type ResetPasswordEmailHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewResetPasswordEmailHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ResetPasswordEmailHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *ResetPasswordEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &types.ResetPasswordEmailRequest{}

	if ok := u.DecodeAndValidate(w, r, request); !ok {
		return
	}

	// find the user in the database so we don't generate reset tokens for users that don't exist
	user, err := u.Repo().User().ReadUserByEmail(request.Email)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			w.WriteHeader(http.StatusOK)
			return
		}

		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	pwResetToken, err := models.NewPasswordResetTokenFromEmail(request.Email)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// this is the only time we'll recover the raw pw reset token, so we store it in the values
	queryVals := url.Values{
		"token": []string{string(pwResetToken.Token)},
		"email": []string{request.Email},
	}

	pwResetToken, err = u.Repo().PasswordResetToken().CreatePasswordResetToken(pwResetToken)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	queryVals.Add("token_id", pwResetToken.ID)

	err = u.Config().UserNotifier.SendPasswordResetEmail(
		&notifier.SendPasswordResetEmailOpts{
			Email: user.Email,
			URL:   fmt.Sprintf("%s/reset_password/finalize?%s", u.Config().ServerRuntimeConfig.ServerURL, queryVals.Encode()),
		},
	)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}
}
