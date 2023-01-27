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

type OrgScopedFactory struct {
	config *server.Config
}

func NewOrgScopedFactory(
	config *server.Config,
) *OrgScopedFactory {
	return &OrgScopedFactory{config}
}

func (p *OrgScopedFactory) Middleware(next http.Handler) http.Handler {
	return &OrgScopedMiddleware{next, p.config}
}

type OrgScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *OrgScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	orgID := reqScopes[types.OrgScope].ResourceID

	org, err := p.config.DB.Repository.Org().ReadOrgByID(orgID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("org with id %s not found ", orgID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	// read the org members to verify that the user has access to this resource
	orgMember, err := p.config.DB.Repository.Org().ReadOrgMemberByUserID(orgID, user.ID, user.UserAccountKind == models.UserAccountService)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("org member with org id %s and user id %s not found ", orgID, user.ID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	// validate against policies
	isValid := false
	policyInput := policies.GetInputFromRequest(r)

	for _, policy := range orgMember.OrgPolicies {
		if !policy.IsCustom {
			switch policy.PolicyName {
			case string(models.PresetPolicyNameOwner):
				allow, err := opa.RunAllowQuery(policies.PresetOrgPolicies.OrgOwnerPolicy.Query, policyInput)

				if err == nil && allow {
					isValid = true
					break
				}
			case string(models.PresetPolicyNameAdmin):
				allow, err := opa.RunAllowQuery(policies.PresetOrgPolicies.OrgAdminPolicy.Query, policyInput)

				if err == nil && allow {
					isValid = true
					break
				}
			case string(models.PresetPolicyNameMember):
				allow, err := opa.RunAllowQuery(policies.PresetOrgPolicies.OrgMemberPolicy.Query, policyInput)

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

	ctx := NewOrganizationContext(r.Context(), org, orgMember)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewOrganizationContext(ctx context.Context, org *models.Organization, orgMember *models.OrganizationMember) context.Context {
	ctx = context.WithValue(ctx, types.OrgScope, org)
	ctx = context.WithValue(ctx, types.OrgMemberLookupKey, orgMember)

	return ctx
}
