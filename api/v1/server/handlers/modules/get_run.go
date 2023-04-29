package modules

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

type ModuleRunGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleRunGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleRunGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleRunGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	var res *types.ModuleRun

	if run.Kind == models.ModuleRunKindPlan && run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindVCS {
		// vcsRepo, err := vcs.GetVCSRepositoryFromModule(m.Config().VCSProviders, module)

		// if err != nil {
		// 	m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

		// 	return
		// }

		// vcsRepo.GetPR(module, run)

		pr, err := m.Repo().GithubPullRequest().ReadGithubPullRequestByGithubID(run.ModuleID, run.ModuleRunConfig.GithubPullRequestID)

		if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		} else if errors.Is(err, repository.RepositoryErrorNotFound) {
			res = run.ToAPIType(nil)
		} else {
			res = run.ToAPIType(pr)
		}

		// m.Repo().GithubPullRequest().ReadGithubPullRequestCommentByGithubID(run.ModuleID, run.ModuleRunConfig.GithubPullRequestID)

		// fmt.Println("ERR", err, errors.Is(err, repository.RepositoryErrorNotFound))

		// if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
		// 	m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		// 	return
		// } else if errors.Is(err, repository.RepositoryErrorNotFound) {
		// 	res = run.ToAPIType(nil)
		// } else {
		// 	pr, err := m.Repo().GithubPullRequest().ReadGithubPullRequestByID(module.TeamID, prComment.GithubPullRequestID)

		// 	if err != nil {
		// 		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		// 		return
		// 	}

		// 	res = run.ToAPIType(pr)
		// }
	} else {
		res = run.ToAPIType(nil)
	}

	m.WriteResult(w, r, res)
}
