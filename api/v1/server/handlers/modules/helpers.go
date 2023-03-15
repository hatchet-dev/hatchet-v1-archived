package modules

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	githubsdk "github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/git/github"
	"github.com/hatchet-dev/hatchet/internal/integrations/valuesstorage/db"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/monitors"
	"github.com/hatchet-dev/hatchet/internal/queuemanager"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/runutils"
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
)

func createManualRun(config *server.Config, module *models.Module, kind models.ModuleRunKind) (*models.ModuleRun, apierrors.RequestError) {
	repo := config.DB.Repository

	run := &models.ModuleRun{
		ModuleID:    module.ID,
		Status:      models.ModuleRunStatusQueued,
		Kind:        kind,
		LogLocation: config.DefaultLogStore.GetID(),
		ModuleRunConfig: models.ModuleRunConfig{
			TriggerKind:            models.ModuleRunTriggerKindManual,
			ModuleValuesVersionID:  module.CurrentModuleValuesVersionID,
			ModuleEnvVarsVersionID: module.CurrentModuleEnvVarsVersionID,
		},
	}

	desc, err := runutils.GenerateRunDescription(config, module, run, models.ModuleRunStatusInProgress, nil)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	run.StatusDescription = desc

	run, err = repo.Module().CreateModuleRun(run)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	// get all monitors for this run
	runMonitors, err := monitors.GetAllMonitorsForModuleRun(repo, module.TeamID, run)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	run, err = repo.Module().AppendModuleRunMonitors(run, runMonitors)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	err = config.ModuleRunQueueManager.Enqueue(module, run, &queuemanager.LockOpts{})

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	err = dispatcher.DispatchModuleRunQueueChecker(config.TemporalClient, &modulequeuechecker.CheckQueueInput{
		TeamID:   module.TeamID,
		ModuleID: module.ID,
	})

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	return run, nil
}

func createLocalRun(config *server.Config, module *models.Module, kind models.ModuleRunKind, hostname string) (*models.ModuleRun, apierrors.RequestError) {
	repo := config.DB.Repository

	run := &models.ModuleRun{
		ModuleID:    module.ID,
		Status:      models.ModuleRunStatusQueued,
		Kind:        kind,
		LogLocation: config.DefaultLogStore.GetID(),
		ModuleRunConfig: models.ModuleRunConfig{
			TriggerKind:            models.ModuleRunTriggerKindManual,
			ModuleValuesVersionID:  module.CurrentModuleValuesVersionID,
			ModuleEnvVarsVersionID: module.CurrentModuleEnvVarsVersionID,
			LocalHostname:          hostname,
		},
	}

	desc, err := runutils.GenerateRunDescription(config, module, run, models.ModuleRunStatusInProgress, nil)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	run.StatusDescription = desc

	run, err = repo.Module().CreateModuleRun(run)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	// get all monitors for this run
	runMonitors, err := monitors.GetAllMonitorsForModuleRun(repo, module.TeamID, run)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	run, err = repo.Module().AppendModuleRunMonitors(run, runMonitors)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	return run, nil
}

func setupGithubDeploymentConfig(config *server.Config, req *types.CreateModuleRequestGithub, team *models.Team, user *models.User) (*models.ModuleDeploymentConfig, apierrors.RequestError) {
	res := &models.ModuleDeploymentConfig{
		ModulePath:              req.Path,
		GithubRepoName:          req.GithubRepositoryName,
		GithubRepoOwner:         req.GithubRepositoryOwner,
		GithubRepoBranch:        req.GithubRepositoryBranch,
		GithubAppInstallationID: req.GithubAppInstallationID,
	}

	gai, reqErr := canAccessGithubAppInstallation(config, req.GithubAppInstallationID, user)

	if reqErr != nil {
		return nil, reqErr
	}

	_, err := createGithubWebhookIfNotExists(config, gai, team.ID, req.GithubRepositoryOwner, req.GithubRepositoryName)

	// TODO(abelanger5): clean up github webhook on subsequent errors

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	return res, nil
}

func getLocalDeploymentConfig(config *server.Config, req *types.CreateModuleRequestLocal, team *models.Team, user *models.User) (*models.ModuleDeploymentConfig, apierrors.RequestError) {
	res := &models.ModuleDeploymentConfig{
		ModulePath: req.LocalPath,
		UserID:     user.ID,
	}

	return res, nil
}

func canAccessGithubAppInstallation(config *server.Config, reqID string, user *models.User) (*models.GithubAppInstallation, apierrors.RequestError) {
	repo := config.DB.Repository

	// ensure that the app installation id exists and the user has access to it
	gai, err := repo.GithubAppInstallation().ReadGithubAppInstallationByID(reqID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			return nil, apierrors.NewErrPassThroughToClient(
				types.APIError{
					Description: "github installation id not found",
					Code:        types.ErrCodeNotFound,
				},
				http.StatusNotFound,
			)
		}

		return nil, apierrors.NewErrInternal(err)
	}

	if gai.GithubAppOAuth.UserID != user.ID {
		return nil, apierrors.NewErrForbidden(
			fmt.Errorf("user %s does not have access to github app installation %s", user.ID, gai.ID),
		)
	}

	return gai, nil
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

func createModuleValuesRaw(config *server.Config, module *models.Module, vals map[string]interface{}, prevVersion uint) (*models.ModuleValuesVersion, error) {
	valuesManager := db.NewDatabaseValuesStore(config.DB.Repository)

	mvv := &models.ModuleValuesVersion{
		ModuleID: module.ID,
		Version:  prevVersion + 1,
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

func createModuleValuesGithub(config *server.Config, module *models.Module, req *types.CreateModuleValuesRequestGithub, prevVersion uint) (*models.ModuleValuesVersion, error) {
	mvv := &models.ModuleValuesVersion{
		ModuleID:                module.ID,
		Version:                 prevVersion + 1,
		Kind:                    models.ModuleValuesVersionKindGithub,
		GithubValuesPath:        req.Path,
		GithubRepoOwner:         req.GithubRepositoryOwner,
		GithubRepoName:          req.GithubRepositoryName,
		GithubRepoBranch:        req.GithubRepositoryBranch,
		GithubAppInstallationID: req.GithubAppInstallationID,
	}

	return config.DB.Repository.ModuleValues().CreateModuleValuesVersion(mvv)
}

func isAllowedDeploymentMechanism(config *server.Config, mechanism string) bool {
	for _, permittedMechanism := range config.ServerRuntimeConfig.PermittedModuleDeploymentMechanisms {
		if permittedMechanism == mechanism {
			return true
		}
	}

	return false
}
