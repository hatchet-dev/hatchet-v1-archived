package github_app

import (
	"errors"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type ListGithubAppInstallationsHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewListGithubAppInstallationsHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ListGithubAppInstallationsHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (g *ListGithubAppInstallationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	req := &types.ListGithubAppInstallationsRequest{}

	if ok := g.DecodeAndValidate(w, r, req); !ok {
		return
	}

	gais, paginate, err := g.Repo().GithubAppInstallation().ListGithubAppInstallationsByUserID(user.ID, repository.WithPage(req.PaginationRequest))

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			resp := &types.ListGithubAppInstallationsResponse{
				Pagination: &types.PaginationResponse{
					NumPages:    1,
					CurrentPage: 1,
					NextPage:    1,
				},
				Rows: make([]*types.GithubAppInstallation, 0),
			}

			g.WriteResult(w, r, resp)
			return
		}

		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	resp := &types.ListGithubAppInstallationsResponse{
		Pagination: paginate.ToAPIType(),
		Rows:       make([]*types.GithubAppInstallation, 0),
	}

	for _, gai := range gais {
		resp.Rows = append(resp.Rows, gai.ToAPIType())
	}

	g.WriteResult(w, r, resp)
}
