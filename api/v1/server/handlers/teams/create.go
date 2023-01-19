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

type TeamCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTeamCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TeamCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TeamCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	org, _ := r.Context().Value(types.OrgScope).(*models.Organization)
	orgMember, _ := r.Context().Value(types.OrgMemberLookupKey).(*models.OrganizationMember)

	request := &types.CreateTeamRequest{}

	if ok := t.DecodeAndValidate(w, r, request); !ok {
		return
	}

	team := &models.Team{
		OrganizationID: org.ID,
		DisplayName:    request.DisplayName,
		TeamMembers: []models.TeamMember{
			{
				OrgMemberID: orgMember.ID,
			},
		},
	}

	team, err := t.Repo().Team().CreateTeam(team)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	teamPolicy, err := t.Repo().Team().ReadPresetTeamPolicyByName(team.ID, models.PresetTeamPolicyNameAdmin)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	teamMember, err := t.Repo().Team().ReadTeamMemberByOrgMemberID(team.ID, orgMember.ID)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	teamMember, err = t.Repo().Team().AppendTeamPolicyToTeamMember(teamMember, teamPolicy)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	t.WriteResult(w, r, team.ToAPIType())
}
