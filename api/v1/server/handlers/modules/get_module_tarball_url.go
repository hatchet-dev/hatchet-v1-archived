package modules

import (
	"context"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/git/github"
	"github.com/hatchet-dev/hatchet/internal/models"

	githubsdk "github.com/google/go-github/v49/github"
)

type ModuleTarballURLGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleTarballURLGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleTarballURLGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleTarballURLGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)

	req := &types.GetModuleTarballURLRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	sha := req.GithubSHA

	client, err := github.GetGithubAppClientFromModule(m.Config(), module)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	if sha == "" {
		branchResp, _, err := client.Repositories.GetBranch(
			context.TODO(),
			module.DeploymentConfig.GithubRepoOwner,
			module.DeploymentConfig.GithubRepoName,
			module.DeploymentConfig.GithubRepoBranch,
			false,
		)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		sha = branchResp.GetCommit().GetSHA()
	}

	ghURL, _, err := client.Repositories.GetArchiveLink(
		context.TODO(),
		module.DeploymentConfig.GithubRepoOwner,
		module.DeploymentConfig.GithubRepoName,
		githubsdk.Zipball,
		&githubsdk.RepositoryContentGetOptions{
			Ref: sha,
		},
		false,
	)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	res := &types.GetModuleTarballURLResponse{
		URL: ghURL.String(),
	}

	m.WriteResult(w, r, res)
}
