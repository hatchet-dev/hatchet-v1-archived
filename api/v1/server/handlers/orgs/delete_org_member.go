package orgs

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type OrgDeleteMemberHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewOrgDeleteMemberHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &OrgDeleteMemberHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (o *OrgDeleteMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orgMember, _ := r.Context().Value(types.OrgMemberScope).(*models.OrganizationMember)

	if reqErr := verifyNotOwner(orgMember); reqErr != nil {
		o.HandleAPIError(w, r, reqErr)
		return
	}

	orgMember, err := o.Repo().Org().DeleteOrgMember(orgMember)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
