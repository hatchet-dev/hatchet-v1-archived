package notifications

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type NotificationListHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewNotificationListHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &NotificationListHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *NotificationListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)
	org, _ := r.Context().Value(types.OrgScope).(*models.Organization)

	req := &types.ListNotificationsRequest{}

	if ok := t.DecodeAndValidate(w, r, req); !ok {
		return
	}

	var teamIDs []string

	if req.TeamID != "" {
		teamIDs = []string{req.TeamID}

	} else {
		// get all teams for the organization that the user is a part of
		teams, _, err := t.Repo().Team().ListTeamsByUserID(
			user.ID,
			org.ID,
		)

		if err != nil {
			t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		teamIDs = make([]string, 0)

		for _, team := range teams {
			teamIDs = append(teamIDs, team.ID)
		}
	}

	notifs, paginate, err := t.Repo().Notification().ListNotificationsByTeamIDs(
		teamIDs,
		&repository.ListNotificationOpts{},
		repository.WithPage(req.PaginationRequest),
	)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListNotificationsResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.NotificationMeta, 0),
	}

	for _, notif := range notifs {
		resp.Rows = append(resp.Rows, notif.ToAPITypeMeta())
	}

	t.WriteResult(w, r, resp)
}
