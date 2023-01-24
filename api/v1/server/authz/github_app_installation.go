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

type GithubAppInstallationScopedFactory struct {
	config *server.Config
}

func NewGithubAppInstallationScopedFactory(
	config *server.Config,
) *GithubAppInstallationScopedFactory {
	return &GithubAppInstallationScopedFactory{config}
}

func (p *GithubAppInstallationScopedFactory) Middleware(next http.Handler) http.Handler {
	return &GithubAppInstallationScopedMiddleware{next, p.config}
}

type GithubAppInstallationScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *GithubAppInstallationScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	gaiID := reqScopes[types.GithubAppInstallationScope].ResourceID

	gai, err := p.config.DB.Repository.GithubAppInstallation().ReadGithubAppInstallationByID(gaiID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("github app installation with id %s not found ", gaiID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	if gai.GithubAppOAuth.UserID != user.ID {
		apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
			fmt.Errorf("user with id %s does not have access to github app installation %s ", user.ID, gai.ID),
		), true)

		return
	}

	ctx := NewGithubAppInstallationContext(r.Context(), gai)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewGithubAppInstallationContext(ctx context.Context, gai *models.GithubAppInstallation) context.Context {
	ctx = context.WithValue(ctx, types.GithubAppInstallationScope, gai)

	return ctx
}
