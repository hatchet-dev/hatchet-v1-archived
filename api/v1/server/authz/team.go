package authz

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/v1/server/authz/policies"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/opa"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type TeamScopedFactory struct {
	config *server.Config
}

func NewTeamScopedFactory(
	config *server.Config,
) *TeamScopedFactory {
	return &TeamScopedFactory{config}
}

func (p *TeamScopedFactory) Middleware(next http.Handler) http.Handler {
	return &TeamScopedMiddleware{next, p.config}
}

type TeamScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *TeamScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	teamID := reqScopes[types.TeamScope].ResourceID

	team, err := p.config.DB.Repository.Team().ReadTeamByID(teamID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("team with id %s not found ", teamID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	orgMember, err := p.config.DB.Repository.Org().ReadOrgMemberByUserID(team.OrganizationID, user.ID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("org member with org id %s and user id %s not found ", team.OrganizationID, user.ID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	// read the team members to verify that the user has access to this resource
	teamMember, err := p.config.DB.Repository.Team().ReadTeamMemberByOrgMemberID(teamID, orgMember.ID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("team member with team id %s and user id %s not found ", teamID, user.ID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	// validate against policies
	isValid := false
	policyInput := policies.GetInputFromRequest(r)

	for _, policy := range teamMember.TeamPolicies {
		if !policy.IsCustom {
			switch policy.PolicyName {
			case string(models.PresetTeamPolicyNameAdmin):
				allow, err := opa.RunAllowQuery(policies.PresetTeamPolicies.TeamAdminPolicy.Query, policyInput)

				if err == nil && allow {
					isValid = true
					break
				}
			case string(models.PresetTeamPolicyNameMember):
				allow, err := opa.RunAllowQuery(policies.PresetTeamPolicies.TeamMemberPolicy.Query, policyInput)

				if err == nil && allow {
					isValid = true
					break
				}
			default:

			}
		}
	}

	if !isValid {
		apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
			fmt.Errorf("no policies permit this action"),
		), true)

		return
	}

	ctx := NewTeamContext(r.Context(), team, teamMember)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewTeamContext(ctx context.Context, team *models.Team, teamMember *models.TeamMember) context.Context {
	ctx = context.WithValue(ctx, types.TeamScope, team)
	ctx = context.WithValue(ctx, types.TeamMemberLookupKey, teamMember)

	return ctx
}
