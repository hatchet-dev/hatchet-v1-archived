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

	// ensure that org member is removed from all dependent teams
	// note that the authz methods shouldn't allow this user to authorize for a specific
	// team without an orgMember entry, but removing them manually makes this cleaner
	teams, _, err := o.Repo().Team().ListTeamsByUserID(orgMember.UserID, orgMember.OrganizationID)

	if err != nil {
		o.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	for _, team := range teams {
		teamMember, err := o.Repo().Team().ReadTeamMemberByOrgMemberID(team.ID, orgMember.ID)

		if err != nil {
			o.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		}

		teamMember, err = o.Repo().Team().DeleteTeamMember(teamMember)

		if err != nil {
			o.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		}
	}
}
