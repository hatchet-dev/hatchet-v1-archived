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

type RunScopedFactory struct {
	config *server.Config
}

func NewRunScopedFactory(
	config *server.Config,
) *RunScopedFactory {
	return &RunScopedFactory{config}
}

func (p *RunScopedFactory) Middleware(next http.Handler) http.Handler {
	return &RunScopedMiddleware{next, p.config}
}

type RunScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *RunScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	runID := reqScopes[types.ModuleRunScope].ResourceID

	mod, err := p.config.DB.Repository.Module().ReadModuleRunByID(runID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("module run with id %s not found ", runID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	ctx := NewRunContext(r.Context(), mod)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewRunContext(ctx context.Context, mod *models.ModuleRun) context.Context {
	ctx = context.WithValue(ctx, types.ModuleRunScope, mod)

	return ctx
}
