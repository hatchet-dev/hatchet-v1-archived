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

type TeamMemberScopedFactory struct {
	config *server.Config
}

func NewTeamMemberScopedFactory(
	config *server.Config,
) *TeamMemberScopedFactory {
	return &TeamMemberScopedFactory{config}
}

func (p *TeamMemberScopedFactory) Middleware(next http.Handler) http.Handler {
	return &TeamMemberScopedMiddleware{next, p.config}
}

type TeamMemberScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *TeamMemberScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	teamMemberID := reqScopes[types.TeamMemberScope].ResourceID

	teamMember, err := p.config.DB.Repository.Team().ReadTeamMemberByID(team.ID, teamMemberID, false)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("team member with id %s not found ", teamMemberID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	ctx := NewTeamMemberContext(r.Context(), teamMember)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewTeamMemberContext(ctx context.Context, teamMember *models.TeamMember) context.Context {
	ctx = context.WithValue(ctx, types.TeamMemberScope, teamMember)

	return ctx
}
