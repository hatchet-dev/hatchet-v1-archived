package users

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type ResetPasswordManualHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewResetPasswordManualHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ResetPasswordManualHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *ResetPasswordManualHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	request := &types.ResetPasswordManualRequest{}

	if ok := u.DecodeAndValidate(w, r, request); !ok {
		return
	}

	if request.OldPassword == request.NewPassword {
		u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Description: "Old and new passwords cannot match. Please try again.",
			Code:        types.ErrCodeBadRequest,
		}, http.StatusUnauthorized, "old and new passwords matched"))

		return
	}

	if verified, err := user.VerifyPassword(request.OldPassword); !verified || err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Description: "Your old password is incorrect. Please try again.",
			Code:        types.ErrCodeBadRequest,
		}, http.StatusUnauthorized, "invalid old password"))

		return
	}

	user.Password = request.NewPassword

	if err := user.HashPassword(); err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	user, err := u.Repo().User().UpdateUser(user)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}
}
