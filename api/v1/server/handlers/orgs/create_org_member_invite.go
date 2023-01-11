package orgs

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/notifier"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type OrgCreateMemberInviteHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewOrgCreateMemberInviteHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &OrgCreateMemberInviteHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (o *OrgCreateMemberInviteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	org, _ := r.Context().Value(types.OrgScope).(*models.Organization)
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	req := &types.CreateOrgMemberInviteRequest{}

	if ok := o.DecodeAndValidate(w, r, req); !ok {
		return
	}

	// ensure that there are no org members with this email address
	candOrgMember, err := o.Repo().Org().ReadOrgMemberByUserOrInviteeEmail(org.ID, req.InviteeEmail)

	if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	if candOrgMember != nil {
		o.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
			types.APIError{
				Code:        types.ErrCodeBadRequest,
				Description: "There is already an organization member with this email address.",
			},
			http.StatusBadRequest,
		))

		return
	}

	policies, reqErr := getPoliciesFromRequest(o.Repo().Org(), org.ID, req.InviteePolicies)

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

	inviteLink, err := models.NewOrganizationInviteLink(o.Config().ServerRuntimeConfig.ServerURL, org.ID)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	rawToken := string(inviteLink.Token)

	inviteLink.InviteeEmail = req.InviteeEmail

	err = inviteLink.Encrypt(o.Config().DB.GetEncryptionKey())

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	orgPolicies := make([]models.OrganizationPolicy, 0)

	for _, policy := range policies {
		orgPolicies = append(orgPolicies, *policy)
	}

	// create the organization member
	orgMember := &models.OrganizationMember{
		InviteLink:  *inviteLink,
		OrgPolicies: orgPolicies,
	}

	orgMember, err = o.Repo().Org().CreateOrgMember(org, orgMember)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	tokID := orgMember.InviteLink.ID

	w.WriteHeader(http.StatusCreated)
	o.WriteResult(w, r, orgMember.ToAPIType(o.Config().DB.GetEncryptionKey()))

	// send invite link to member
	if reqErr := sendInviteEmail(o.Config(), inviteLink.InviteeEmail, rawToken, tokID, org.DisplayName, user.Email); reqErr != nil {
		// we've already written a success message, so don't rewrite it
		o.HandleAPIErrorNoWrite(w, r, reqErr)
	}
}

func sendInviteEmail(config *server.Config, targetEmail, tok, tokID, orgName, inviterAddress string) apierrors.RequestError {
	// this is the only time we'll recover the raw pw reset token, so we store it in the values
	queryVals := url.Values{
		"token":           []string{tok},
		"invite_id":       []string{tokID},
		"org_name":        []string{orgName},
		"inviter_address": []string{inviterAddress},
	}

	err := config.UserNotifier.SendInviteLinkEmail(
		&notifier.SendInviteLinkEmailOpts{
			Email:            targetEmail,
			URL:              fmt.Sprintf("%s/organization_invite/accept?%s", config.ServerRuntimeConfig.ServerURL, queryVals.Encode()),
			OrganizationName: orgName,
			InviterAddress:   inviterAddress,
		},
	)

	if err != nil {
		return apierrors.NewErrInternal(err)
	}

	return nil
}
