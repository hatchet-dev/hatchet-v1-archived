package pats

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type PATDeleteHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewPATDeleteHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &PATDeleteHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *PATDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	tokenID, reqErr := handlerutils.GetURLParamString(r, types.PersonalAccessTokenURLParam)

	if reqErr != nil {
		u.HandleAPIError(w, r, reqErr)
		return
	}

	pat, err := u.Repo().PersonalAccessToken().ReadPersonalAccessToken(user.ID, tokenID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
				types.APIError{
					Code:        types.ErrCodeNotFound,
					Description: types.GenericResourceNotFound,
				},
				http.StatusNotFound,
			))

			return
		}

		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	pat, err = u.Repo().PersonalAccessToken().DeletePersonalAccessToken(pat)

	fmt.Println("RESULT OF DELETE CALL IS", pat, err)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			u.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
				types.APIError{
					Code:        types.ErrCodeNotFound,
					Description: types.GenericResourceNotFound,
				},
				http.StatusNotFound,
			))

			return
		}

		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	u.WriteResult(w, r, pat.ToAPIType())
}
