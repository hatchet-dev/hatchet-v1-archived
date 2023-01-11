package users

import (
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
)

type ResendVerifyEmailHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewResendVerifyEmailHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ResendVerifyEmailHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *ResendVerifyEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	if reqErr := sendVerifyEmail(u.Config(), user); reqErr != nil {
		u.HandleAPIError(w, r, reqErr)
	}
}

func sendVerifyEmail(config *server.Config, user *models.User) apierrors.RequestError {
	emailVerifyTok, err := models.NewVerifyEmailTokenFromEmail(user.Email)

	if err != nil {
		return apierrors.NewErrInternal(err)
	}

	// this is the only time we'll recover the raw pw reset token, so we store it in the values
	queryVals := url.Values{
		"token": []string{string(emailVerifyTok.Token)},
	}

	emailVerifyTok, err = config.DB.Repository.VerifyEmailToken().CreateVerifyEmailToken(emailVerifyTok)

	if err != nil {
		return apierrors.NewErrInternal(err)
	}

	queryVals.Add("token_id", emailVerifyTok.ID)

	err = config.UserNotifier.SendVerificationEmail(
		&notifier.SendVerificationEmailOpts{
			Email: user.Email,
			URL:   fmt.Sprintf("%s/verify_email/finalize?%s", config.ServerRuntimeConfig.ServerURL, queryVals.Encode()),
		},
	)

	if err != nil {
		return apierrors.NewErrInternal(err)
	}

	return nil
}
