package teams

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/monitors"
	"github.com/hatchet-dev/hatchet/internal/temporal"
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"
	"go.temporal.io/api/workflowservice/v1"
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

	teamMember, err := t.Repo().Team().ReadTeamMemberByOrgMemberID(team.ID, orgMember.ID, false)

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

	// create a temporal namespace for the team
	tc, err := t.Config().TemporalClient.GetClient(temporal.DefaultQueueName)

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	retention := 24 * 30 * time.Hour

	_, err = tc.WorkflowService().RegisterNamespace(context.Background(), &workflowservice.RegisterNamespaceRequest{
		Namespace:                        team.ID,
		WorkflowExecutionRetentionPeriod: &retention,
	})

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	// create service account team member by creating a new user and team member and appending them to
	// the team
	saUser, err := t.Repo().User().CreateUser(&models.User{
		DisplayName:     fmt.Sprintf("%s-runner-service-account", team.ID),
		UserAccountKind: models.UserAccountService,
		Email:           fmt.Sprintf("%s-runner-service-account@hatchet.run", team.ID),
		EmailVerified:   true,
	})

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	// create org member without an organization policy. this still gives the service account access
	// to team-scoped endpoints, but not organization endpoints
	saOrgMember, err := t.Repo().Org().CreateOrgMember(org, &models.OrganizationMember{
		IsServiceAccountRunner: true,
		UserID:                 saUser.ID,
	})

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	saTeamMember, err := t.Repo().Team().CreateTeamMember(team, &models.TeamMember{
		IsServiceAccountRunner: true,
		OrgMemberID:            saOrgMember.ID,
	})

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	// read member policy
	saMemberPolicy, err := t.Repo().Team().ReadPresetTeamPolicyByName(team.ID, models.PresetTeamPolicyNameMember)

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	saTeamMember, err = t.Repo().Team().AppendTeamPolicyToTeamMember(saTeamMember, saMemberPolicy)

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	team.ServiceAccountRunnerID = saUser.ID

	team, err = t.Repo().Team().UpdateTeam(team)

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	// create notification inbox
	inbox := &models.NotificationInbox{
		TeamID: team.ID,
	}

	inbox, err = t.Repo().Notification().CreateNotificationInbox(inbox)

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	// create a default drift detection monitor for the team
	monitor := &models.ModuleMonitor{
		TeamID:           team.ID,
		Kind:             models.MonitorKindPlan,
		PresetPolicyName: models.ModuleMonitorPresetPolicyNameDrift,
		CurrentMonitorPolicyBytesVersion: models.MonitorPolicyBytesVersion{
			Version:     1,
			PolicyBytes: monitors.PresetDriftDetectionPolicy,
		},
	}

	monitor, err = t.Repo().ModuleMonitor().CreateModuleMonitor(monitor)

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}

	err = dispatcher.DispatchCronMonitor(t.Config().TemporalClient, team.ID, monitor.ID, "0 */6 * * *")

	if err != nil {
		t.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))
		return
	}
}
