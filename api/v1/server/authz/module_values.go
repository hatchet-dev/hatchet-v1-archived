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

type ModuleValuesScopedFactory struct {
	config *server.Config
}

func NewModuleValuesScopedFactory(
	config *server.Config,
) *ModuleValuesScopedFactory {
	return &ModuleValuesScopedFactory{config}
}

func (p *ModuleValuesScopedFactory) Middleware(next http.Handler) http.Handler {
	return &ModuleValuesScopedMiddleware{next, p.config}
}

type ModuleValuesScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *ModuleValuesScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)

	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	valuesID := reqScopes[types.ModuleValuesScope].ResourceID

	mvv, err := p.config.DB.Repository.ModuleValues().ReadModuleValuesVersionByID(module.ID, valuesID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("values with id %s not found ", valuesID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	ctx := NewModuleValuesContext(r.Context(), mvv)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewModuleValuesContext(ctx context.Context, mvv *models.ModuleValuesVersion) context.Context {
	ctx = context.WithValue(ctx, types.ModuleValuesScope, mvv)

	return ctx
}
