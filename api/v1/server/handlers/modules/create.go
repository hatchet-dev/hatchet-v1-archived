package modules

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type ModuleCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	request := &types.CreateModuleRequest{}

	if ok := m.DecodeAndValidate(w, r, request); !ok {
		return
	}

	mod := &models.Module{
		TeamID:              team.ID,
		Name:                request.Name,
		DeploymentMechanism: models.DeploymentMechanismGithub,
	}

	if request.DeploymentGithub == nil {
		m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: fmt.Sprintf("at least one deployment mechanism must be specified"),
		}, http.StatusBadRequest))

		return
	}

	if github := request.DeploymentGithub; github != nil {
		// ensure that the app installation id exists and the user has access to it
		gai, err := m.Repo().GithubAppInstallation().ReadGithubAppInstallationByID(github.GithubAppInstallationID)

		// TODO(abelanger5): handle this gracefully
		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		if gai.GithubAppOAuth.UserID != user.ID {
			m.HandleAPIError(w, r, apierrors.NewErrForbidden(
				fmt.Errorf("user %s does not have access to github app installation %s", user.ID, gai.ID),
			))

			return
		}

		mod.DeploymentConfig = models.ModuleDeploymentConfig{
			ModulePath:              github.Path,
			GithubRepoName:          github.GithubRepositoryName,
			GithubRepoOwner:         github.GithubRepositoryOwner,
			GithubRepoBranch:        github.GithubRepositoryBranch,
			GithubAppInstallationID: github.GithubAppInstallationID,
		}
	}

	mod, err := m.Repo().Module().CreateModule(mod)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	m.WriteResult(w, r, mod.ToAPIType())
}
