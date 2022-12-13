package users

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type UserGetCurrentHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewUserGetCurrentHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &UserGetCurrentHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *UserGetCurrentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	u.WriteResult(w, r, user.ToAPIType())
}
