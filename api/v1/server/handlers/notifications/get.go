package notifications

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type NotificationGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewNotificationGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &NotificationGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *NotificationGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	notif, _ := r.Context().Value(types.NotificationScope).(*models.Notification)

	t.WriteResult(w, r, notif.ToAPIType())
}
