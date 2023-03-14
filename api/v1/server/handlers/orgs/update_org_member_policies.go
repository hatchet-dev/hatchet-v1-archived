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

type OrgUpdateMemberPoliciesHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewOrgUpdateMemberPoliciesHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &OrgUpdateMemberPoliciesHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (o *OrgUpdateMemberPoliciesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	org, _ := r.Context().Value(types.OrgScope).(*models.Organization)
	orgMember, _ := r.Context().Value(types.OrgMemberScope).(*models.OrganizationMember)

	if reqErr := verifyNotOwner(orgMember); reqErr != nil {
		o.HandleAPIError(w, r, reqErr)
		return
	}

	req := &types.UpdateOrgMemberPoliciesRequest{}

	if ok := o.DecodeAndValidate(w, r, req); !ok {
		return
	}

	policies, reqErr := getPoliciesFromRequest(o.Repo().Org(), org.ID, req.Policies)

	if reqErr != nil {
		o.HandleAPIError(w, r, reqErr)

		return
	} else if len(policies) == 0 {
		o.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
			types.APIError{
				Description: "At least one policy must be requested",
				Code:        types.ErrCodeBadRequest,
			},
			http.StatusBadRequest,
		))

		return
	}

	orgMember, err := o.Repo().Org().ReplaceOrgPoliciesForOrgMember(orgMember, policies)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	o.WriteResult(w, r, orgMember.ToAPIType(o.Config().DB.GetEncryptionKey(), o.Config().ServerRuntimeConfig.ServerURL, org.DisplayName))
}
