package orgs

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

type OrgListMembersHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewOrgListMembersHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &OrgListMembersHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (o *OrgListMembersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	org, _ := r.Context().Value(types.OrgScope).(*models.Organization)

	req := &types.ListOrgMembersRequest{}

	if ok := o.DecodeAndValidate(w, r, req); !ok {
		return
	}

	members, paginate, err := o.Repo().Org().ListOrgMembersByOrgID(
		org.ID,
		false,
		repository.WithPage(req.PaginationRequest),
	)

	if err != nil {
		o.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListOrgMembersResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.OrganizationMemberSanitized, 0),
	}

	for _, member := range members {
		resp.Rows = append(resp.Rows, member.ToAPITypeSanitized())
	}

	o.WriteResult(w, r, resp)
}
