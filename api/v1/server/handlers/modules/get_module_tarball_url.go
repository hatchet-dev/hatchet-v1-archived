package modules

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
	"github.com/hatchet-dev/hatchet/internal/models"
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

	vcsRepo, err := vcs.GetVCSRepositoryFromModule(m.Config().VCSProviders, module)

	if err != nil {
		m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

		return
	}

	if sha == "" {
		branchResp, err := vcsRepo.GetBranch(module.DeploymentConfig.GitRepoBranch)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		sha = branchResp.GetLatestRef()
	}

	archiveLink, err := vcsRepo.GetArchiveLink(sha)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	res := &types.GetModuleTarballURLResponse{
		URL: archiveLink.String(),
	}

	m.WriteResult(w, r, res)
}
