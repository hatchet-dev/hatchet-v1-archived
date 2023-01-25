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

type ModuleScopedFactory struct {
	config *server.Config
}

func NewModuleScopedFactory(
	config *server.Config,
) *ModuleScopedFactory {
	return &ModuleScopedFactory{config}
}

func (p *ModuleScopedFactory) Middleware(next http.Handler) http.Handler {
	return &ModuleScopedMiddleware{next, p.config}
}

type ModuleScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *ModuleScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	moduleID := reqScopes[types.ModuleScope].ResourceID

	mod, err := p.config.DB.Repository.Module().ReadModuleByID(moduleID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("module with id %s not found ", moduleID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	ctx := NewModuleContext(r.Context(), mod)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewModuleContext(ctx context.Context, mod *models.Module) context.Context {
	ctx = context.WithValue(ctx, types.ModuleScope, mod)

	return ctx
}
