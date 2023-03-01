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

type MonitorScopedFactory struct {
	config *server.Config
}

func NewMonitorScopedFactory(
	config *server.Config,
) *MonitorScopedFactory {
	return &MonitorScopedFactory{config}
}

func (p *MonitorScopedFactory) Middleware(next http.Handler) http.Handler {
	return &MonitorScopedMiddleware{next, p.config}
}

type MonitorScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *MonitorScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	monitorID := reqScopes[types.MonitorScope].ResourceID

	monitor, err := p.config.DB.Repository.ModuleMonitor().ReadModuleMonitorByID(team.ID, monitorID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("monitor with id %s not found ", monitorID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	ctx := NewMonitorContext(r.Context(), monitor)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewMonitorContext(ctx context.Context, monitor *models.ModuleMonitor) context.Context {
	ctx = context.WithValue(ctx, types.MonitorScope, monitor)

	return ctx
}
