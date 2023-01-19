package teams

import (
	"errors"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type TeamAddMemberHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTeamAddMemberHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TeamAddMemberHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TeamAddMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	req := &types.TeamAddMemberRequest{}

	if ok := t.DecodeAndValidate(w, r, req); !ok {
		return
	}

	// look up the org member to make sure they exist
	orgMember, err := t.Repo().Org().ReadOrgMemberByID(team.OrganizationID, req.OrgMemberID)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// make sure the org member is not already a part of the team
	candTeamMember, err := t.Repo().Team().ReadTeamMemberByOrgMemberID(team.ID, orgMember.ID)

	if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	if candTeamMember != nil {
		t.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
			types.APIError{
				Code:        types.ErrCodeBadRequest,
				Description: "There is already a team member that corresponds to this organization member.",
			},
			http.StatusBadRequest,
		))

		return
	}

	policies, reqErr := getPoliciesFromRequest(t.Repo().Team(), team.ID, req.Policies)

	if reqErr != nil {
		t.HandleAPIError(w, r, reqErr)

		return
	} else if len(policies) == 0 {
		t.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
			types.APIError{
				Description: "At least one policy must be requested",
				Code:        types.ErrCodeBadRequest,
			},
			http.StatusBadRequest,
		))

		return
	}

	teamPolicies := make([]models.TeamPolicy, 0)

	for _, policy := range policies {
		teamPolicies = append(teamPolicies, *policy)
	}

	// create the team member
	teamMember := &models.TeamMember{
		OrgMemberID:  orgMember.ID,
		TeamPolicies: teamPolicies,
	}

	teamMember, err = t.Repo().Team().CreateTeamMember(team, teamMember)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	t.WriteResult(w, r, teamMember.ToAPIType())
}
