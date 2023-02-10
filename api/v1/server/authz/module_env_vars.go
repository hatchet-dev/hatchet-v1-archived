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

type ModuleEnvVarsScopedFactory struct {
	config *server.Config
}

func NewModuleEnvVarsScopedFactory(
	config *server.Config,
) *ModuleEnvVarsScopedFactory {
	return &ModuleEnvVarsScopedFactory{config}
}

func (p *ModuleEnvVarsScopedFactory) Middleware(next http.Handler) http.Handler {
	return &ModuleEnvVarsScopedMiddleware{next, p.config}
}

type ModuleEnvVarsScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *ModuleEnvVarsScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)

	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	envVarsID := reqScopes[types.ModuleEnvVarScope].ResourceID

	mev, err := p.config.DB.Repository.ModuleEnvVars().ReadModuleEnvVarsVersionByID(module.ID, envVarsID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("env vars with id %s not found ", envVarsID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	ctx := NewModuleEnvVarsContext(r.Context(), mev)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewModuleEnvVarsContext(ctx context.Context, mev *models.ModuleEnvVarsVersion) context.Context {
	ctx = context.WithValue(ctx, types.ModuleEnvVarScope, mev)

	return ctx
}
