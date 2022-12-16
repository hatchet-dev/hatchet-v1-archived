package orgs

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type OrgGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewOrgGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &OrgGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (o *OrgGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	org, _ := r.Context().Value(types.OrgScope).(*models.Organization)

	o.WriteResult(w, r, org.ToAPIType())
}
