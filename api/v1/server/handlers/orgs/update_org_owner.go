package orgs

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

type OrgUpdateOwnerHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewOrgUpdateOwnerHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &OrgUpdateOwnerHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (o *OrgUpdateOwnerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	org, _ := r.Context().Value(types.OrgScope).(*models.Organization)
	orgMember, _ := r.Context().Value(types.OrgMemberLookupKey).(*models.OrganizationMember)

	req := &types.UpdateOrgOwnerRequest{}

	if ok := o.DecodeAndValidate(w, r, req); !ok {
		return
	}

	ownerPolicy, err := o.Repo().Org().ReadPresetPolicyByName(org.ID, models.PresetPolicyNameOwner)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	adminPolicy, err := o.Repo().Org().ReadPresetPolicyByName(org.ID, models.PresetPolicyNameAdmin)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// read the new member
	newOrgMember, err := o.Repo().Org().ReadOrgMemberByUserID(org.ID, req.NewOwnerMemberID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			o.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
				types.APIError{
					Code:        types.ErrCodeNotFound,
					Description: fmt.Sprintf("could not find organization member with id %s in this organization", req.NewOwnerMemberID),
				},
				http.StatusNotFound,
			))

			return
		}

		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	org.OwnerID = orgMember.UserID

	// attach the owner policy to the new org member
	newOrgMember, err = o.Repo().Org().AppendOrgPolicyToOrgMember(newOrgMember, ownerPolicy)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// attach the admin policy to the old org member
	orgMember, err = o.Repo().Org().RemoveOrgPolicyFromOrgMember(orgMember, adminPolicy)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// remove the owner policy from the old org member
	orgMember, err = o.Repo().Org().RemoveOrgPolicyFromOrgMember(orgMember, ownerPolicy)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}
}
