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

type OrgCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewOrgCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &OrgCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (u *OrgCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	request := &types.CreateOrganizationRequest{}

	if ok := u.DecodeAndValidate(w, r, request); !ok {
		return
	}

	org := &models.Organization{
		DisplayName: request.DisplayName,
		OwnerID:     user.ID,
		OrgMembers: []models.OrganizationMember{
			{
				InviteAccepted: true,
				UserID:         user.ID,
			},
		},
	}

	org, err := u.Repo().Org().CreateOrg(org)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// TODO(abelanger5): create a transaction for the org creation, in case subsequent policy attachment
	// fails and org is left without a member with an owner. This can likely be done with a gorm hook of some
	// kind.
	orgPolicy, err := u.Repo().Org().ReadPresetPolicyByName(org.ID, models.PresetPolicyNameOwner)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	orgMember, err := u.Repo().Org().ReadOrgMemberByUserID(org.ID, user.ID, false)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	orgMember, err = u.Repo().Org().AppendOrgPolicyToOrgMember(orgMember, orgPolicy)

	if err != nil {
		u.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	u.WriteResult(w, r, org.ToAPIType())
}
