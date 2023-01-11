package metadata

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

type ServerMetadataGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewServerMetadataGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ServerMetadataGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (o *ServerMetadataGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o.WriteResult(w, r, o.Config().ToAPIServerMetadataType())
}
