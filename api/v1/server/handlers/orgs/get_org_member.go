package orgs

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type OrgGetMemberHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewOrgGetMemberHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &OrgGetMemberHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (o *OrgGetMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orgMember, _ := r.Context().Value(types.OrgMemberScope).(*models.OrganizationMember)

	o.WriteResult(w, r, orgMember.ToAPIType(o.Config().DB.GetEncryptionKey()))
}
