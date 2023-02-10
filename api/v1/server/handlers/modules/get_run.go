package modules

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
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
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	var res *types.ModuleRun

	if run.Kind == models.ModuleRunKindPlan && run.ModuleRunConfig.TriggerKind == models.ModuleRunTriggerKindGithub {
		prComment, err := m.Repo().GithubPullRequest().ReadGithubPullRequestCommentByGithubID(run.ModuleID, run.ModuleRunConfig.GithubCommentID)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		pr, err := m.Repo().GithubPullRequest().ReadGithubPullRequestByID(module.TeamID, prComment.GithubPullRequestID)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		res = run.ToAPIType(pr)
	}

	m.WriteResult(w, r, res)
}
