package modules

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/git/github"
	"github.com/hatchet-dev/hatchet/internal/integrations/valuesstorage/db"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"

	githubsdk "github.com/google/go-github/v49/github"
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

		if err != nil {
			if errors.Is(err, repository.RepositoryErrorNotFound) {
				m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
					types.APIError{
						Description: "github installation id not found",
						Code:        types.ErrCodeNotFound,
					},
					http.StatusNotFound,
				))

				return
			}

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

		_, err = createGithubWebhookIfNotExists(m.Config(), gai, team.ID, github.GithubRepositoryOwner, github.GithubRepositoryName)

		// TODO(abelanger5): clean up github webhook on subsequent errors

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}
	}

	mod, err := m.Repo().Module().CreateModule(mod)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	var mvv *models.ModuleValuesVersion

	if valuesGithub := request.ValuesGithub; valuesGithub != nil {
		// ensure that the app installation id exists and the user has access to it
		gai, err := m.Repo().GithubAppInstallation().ReadGithubAppInstallationByID(valuesGithub.GithubAppInstallationID)

		if err != nil {
			if errors.Is(err, repository.RepositoryErrorNotFound) {
				m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
					types.APIError{
						Description: "github installation id not found",
						Code:        types.ErrCodeNotFound,
					},
					http.StatusNotFound,
				))

				return
			}

			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		if gai.GithubAppOAuth.UserID != user.ID {
			m.HandleAPIError(w, r, apierrors.NewErrForbidden(
				fmt.Errorf("user %s does not have access to github app installation %s", user.ID, gai.ID),
			))

			return
		}

		mvv, err = createModuleValuesGithub(m.Config(), mod, valuesGithub)
	} else {
		mvv, err = createModuleValuesRaw(m.Config(), mod, request.ValuesRaw)
	}

	// set values version, this is updated later to reduce DB queries
	mod.CurrentModuleValuesVersionID = mvv.ID

	// create env vars
	mev, err := models.NewModuleEnvVarsVersion(mod.ID, 0, request.EnvVars)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	mev, err = m.Repo().ModuleEnvVars().CreateModuleEnvVarsVersion(mev)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	mod.CurrentModuleEnvVarsVersionID = mev.ID

	mod, err = m.Repo().Module().UpdateModule(mod)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	m.WriteResult(w, r, mod.ToAPIType())
}

func createGithubWebhookIfNotExists(config *server.Config, gai *models.GithubAppInstallation, teamID, repoOwner, repoName string) (*models.GithubWebhook, error) {
	gw, err := config.DB.Repository.GithubWebhook().ReadGithubWebhookByTeamID(teamID, repoOwner, repoName)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			return createGithubWebhook(config, gai, teamID, repoOwner, repoName)
		}

		return nil, err
	}

	return gw, err

}

func createGithubWebhook(config *server.Config, gai *models.GithubAppInstallation, teamID, repoOwner, repoName string) (*models.GithubWebhook, error) {
	gw, err := models.NewGithubWebhook(teamID, repoOwner, repoName)

	if err != nil {
		return nil, err
	}

	gw, err = config.DB.Repository.GithubWebhook().CreateGithubWebhook(gw)

	if err != nil {
		return nil, err
	}

	webhookURL := fmt.Sprintf("%s/api/v1/teams/%s/github_incoming/%s", config.ServerRuntimeConfig.ServerURL, teamID, gw.ID)

	// config.DB.Repository.GithubWebhook().Create(teamID, repoName)
	client, err := github.GetGithubAppClientFromGAI(config, gai)

	if err != nil {
		return nil, err
	}

	_, _, err = client.Repositories.CreateHook(
		context.Background(), repoOwner, repoName, &githubsdk.Hook{
			Config: map[string]interface{}{
				"url":          webhookURL,
				"content_type": "json",
				"secret":       string(gw.SigningSecret),
			},
			Events: []string{"pull_request", "push"},
			Active: githubsdk.Bool(true),
		},
	)

	if err != nil {
		return nil, err
	}

	return gw, nil
}

func createModuleValuesRaw(config *server.Config, module *models.Module, vals map[string]interface{}) (*models.ModuleValuesVersion, error) {
	valuesManager := db.NewDatabaseValuesStore(config.DB.Repository)

	mvv := &models.ModuleValuesVersion{
		ModuleID: module.ID,
		Version:  1,
		Kind:     models.ModuleValuesVersionKindDatabase,
	}

	mvv, err := config.DB.Repository.ModuleValues().CreateModuleValuesVersion(mvv)

	if err != nil {
		return nil, err
	}

	err = valuesManager.WriteValues(mvv, vals)

	if err != nil {
		return nil, err
	}

	return mvv, nil
}

func createModuleValuesGithub(config *server.Config, module *models.Module, req *types.CreateModuleValuesRequestGithub) (*models.ModuleValuesVersion, error) {
	mvv := &models.ModuleValuesVersion{
		ModuleID:                module.ID,
		Version:                 1,
		Kind:                    models.ModuleValuesVersionKindGithub,
		GithubValuesPath:        req.Path,
		GithubRepoOwner:         req.GithubRepositoryOwner,
		GithubRepoName:          req.GithubRepositoryName,
		GithubRepoBranch:        req.GithubRepositoryBranch,
		GithubAppInstallationID: req.GithubAppInstallationID,
	}

	return config.DB.Repository.ModuleValues().CreateModuleValuesVersion(mvv)
}
