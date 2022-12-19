package authz

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type OrgMemberScopedFactory struct {
	config *server.Config
}

func NewOrgMemberScopedFactory(
	config *server.Config,
) *OrgMemberScopedFactory {
	return &OrgMemberScopedFactory{config}
}

func (p *OrgMemberScopedFactory) Middleware(next http.Handler) http.Handler {
	return &OrgMemberScopedMiddleware{next, p.config}
}

type OrgMemberScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *OrgMemberScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	org, _ := r.Context().Value(types.OrgScope).(*models.Organization)

	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	orgMemberID := reqScopes[types.OrgMemberScope].ResourceID

	orgMember, err := p.config.DB.Repository.Org().ReadOrgMemberByID(org.ID, orgMemberID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("org member with id %s not found ", orgMemberID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	ctx := NewOrganizationMemberContext(r.Context(), orgMember)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewOrganizationMemberContext(ctx context.Context, orgMember *models.OrganizationMember) context.Context {
	ctx = context.WithValue(ctx, types.OrgMemberScope, orgMember)

	return ctx
}
