package pats

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type PATCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewPATCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &PATCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *PATCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	request := &types.CreatePATRequest{}

	if ok := u.DecodeAndValidate(w, r, request); !ok {
		return
	}

	// ensure that there are no PATs with the same display name
	existingPAT, _ := u.Repo().PersonalAccessToken().ReadPersonalAccessTokenByDisplayName(user.ID, request.DisplayName)

	if existingPAT != nil {
		u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: fmt.Sprintf("Personal access token already exists with display_name %s for this user", request.DisplayName),
		}, http.StatusBadRequest))

		return
	}

	pat, err := models.NewPATFromUserID(request.DisplayName, user.ID)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	rawToken, err := token.GenerateTokenFromPAT(pat, u.Config().TokenOpts)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// write the user to the db
	pat, err = u.Repo().PersonalAccessToken().CreatePersonalAccessToken(pat)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	u.WriteResult(w, r, &types.CreatePATResponse{
		PersonalAccessToken: *pat.ToAPIType(),
		Token:               rawToken,
	})
}
