package modules

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/valuesstorage/db"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs/github"
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

	desc, err := runutils.GenerateRunDescription(config, module, run, models.ModuleRunStatusInProgress, nil, nil)

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

	desc, err := runutils.GenerateRunDescription(config, module, run, models.ModuleRunStatusInProgress, nil, nil)

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
		GitRepoName:             req.GithubRepositoryName,
		GitRepoOwner:            req.GithubRepositoryOwner,
		GitRepoBranch:           req.GithubRepositoryBranch,
		GithubAppInstallationID: req.GithubAppInstallationID,
	}

	_, reqErr := canAccessGithubAppInstallation(config, req.GithubAppInstallationID, user)

	if reqErr != nil {
		return nil, reqErr
	}

	fact := config.VCSProviders[vcs.VCSRepositoryKindGithub]
	provider, err := github.ToGithubVCSProvider(fact)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	githubVCS, err := provider.GetVCSRepositoryFromModule(res)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	err = githubVCS.SetupRepository(team.ID)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

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
		Kind:                    models.ModuleValuesVersionKindVCS,
		GitValuesPath:           req.Path,
		GitRepoOwner:            req.GithubRepositoryOwner,
		GitRepoName:             req.GithubRepositoryName,
		GitRepoBranch:           req.GithubRepositoryBranch,
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
