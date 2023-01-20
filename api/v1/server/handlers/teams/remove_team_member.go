package teams

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type TeamRemoveMemberHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTeamRemoveMemberHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TeamRemoveMemberHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TeamRemoveMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	teamMember, _ := r.Context().Value(types.TeamMemberScope).(*models.TeamMember)

	teamMember, err := t.Repo().Team().DeleteTeamMember(teamMember)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
