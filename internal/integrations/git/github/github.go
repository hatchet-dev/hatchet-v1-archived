package github

import (
	"fmt"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

// GetGithubAppClientFromGAI gets the github client from the GithubAppInstallation model
func GetGithubAppClientFromGAI(config *server.Config, gai *models.GithubAppInstallation) (*github.Client, error) {
	itr, err := ghinstallation.New(
		http.DefaultTransport,
		config.GithubApp.AppID,
		gai.InstallationID,
		config.GithubApp.Secret,
	)

	if err != nil {
		return nil, err
	}

	return github.NewClient(&http.Client{Transport: itr}), nil
}

// GetGithubAppClientFromModule gets the github client from a module and a user id
func GetGithubAppClientFromModule(config *server.Config, mod *models.Module) (*github.Client, error) {
	if mod.DeploymentConfig.GithubAppInstallationID == "" {
		return nil, fmt.Errorf("module does not have github app installation id param set")
	}

	gai, err := config.DB.Repository.GithubAppInstallation().ReadGithubAppInstallationByID(mod.DeploymentConfig.GithubAppInstallationID)

	if err != nil {
		return nil, err
	}

	return GetGithubAppClientFromGAI(config, gai)
}
