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

type OrgAcceptMemberInviteHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewOrgAcceptMemberInviteHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &OrgAcceptMemberInviteHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (o *OrgAcceptMemberInviteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)
	urlParamInviteID, reqErr := handlerutils.GetURLParamString(r, types.URLParamOrgMemberInviteID)

	if reqErr != nil {
		o.HandleAPIError(w, r, reqErr)
		return
	}

	urlParamInviteTok, reqErr := handlerutils.GetURLParamString(r, types.URLParamOrgMemberInviteTok)

	if reqErr != nil {
		o.HandleAPIError(w, r, reqErr)
		return
	}

	invite, err := o.Repo().Org().ReadOrgInviteByID(urlParamInviteID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			o.HandleAPIError(w, r, apierrors.NewErrForbidden(
				fmt.Errorf("could not find invite with id %s", urlParamInviteID),
			))

			return
		}

		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	invite.Decrypt(o.Config().DB.GetEncryptionKey())

	// validate the hashed token
	valid, err := invite.VerifyToken([]byte(urlParamInviteTok))

	if !valid || err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrForbidden(
			fmt.Errorf("invalid token for invite with id %s", invite.ID),
		))

		return
	}

	if invite.Used {
		o.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
			types.APIError{
				Code:        types.ErrCodeBadRequest,
				Description: "Invite has already been used. Please contact your organization administrator to generate a new invite.",
			},
			http.StatusBadRequest,
		))

		return
	}

	if invite.InviteeEmail != user.Email {
		o.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
			types.APIError{
				Code:        types.ErrCodeBadRequest,
				Description: "Wrong email for invite. Please contact your organization administrator to check the invite email.",
			},
			http.StatusBadRequest,
			fmt.Sprintf("got email %s, expected email %s", user.Email, invite.InviteeEmail),
		))

		return
	}

	if invite.IsExpired() {
		o.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
			types.APIError{
				Code:        types.ErrCodeBadRequest,
				Description: "Invite has expired. Please contact your organization administrator to generate a new invite.",
			},
			http.StatusBadRequest,
		))

		return
	}

	invite.Used = true
	invite, err = o.Repo().Org().UpdateOrgInvite(invite)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// add the user to the organization member list
	orgMember, err := o.Repo().Org().ReadOrgMemberByID(invite.OrganizationID, invite.OrganizationMemberID, false)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	orgMember.UserID = user.ID
	orgMember.InviteAccepted = true

	orgMember, err = o.Repo().Org().UpdateOrgMember(orgMember)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}
}
