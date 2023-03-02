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

type NotificationScopedFactory struct {
	config *server.Config
}

func NewNotificationScopedFactory(
	config *server.Config,
) *NotificationScopedFactory {
	return &NotificationScopedFactory{config}
}

func (p *NotificationScopedFactory) Middleware(next http.Handler) http.Handler {
	return &NotificationScopedMiddleware{next, p.config}
}

type NotificationScopedMiddleware struct {
	next   http.Handler
	config *server.Config
}

func (p *NotificationScopedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)
	notifID := reqScopes[types.NotificationScope].ResourceID

	notif, err := p.config.DB.Repository.Notification().ReadNotificationByID(team.ID, notifID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrForbidden(
				fmt.Errorf("notification with id %s not found ", notifID),
			), true)
		} else {
			apierrors.HandleAPIError(p.config.Logger, p.config.ErrorAlerter, w, r, apierrors.NewErrInternal(err), true)
		}

		return
	}

	ctx := NewNotificationContext(r.Context(), notif)
	r = r.Clone(ctx)
	p.next.ServeHTTP(w, r)
}

func NewNotificationContext(ctx context.Context, notif *models.Notification) context.Context {
	ctx = context.WithValue(ctx, types.NotificationScope, notif)

	return ctx
}
