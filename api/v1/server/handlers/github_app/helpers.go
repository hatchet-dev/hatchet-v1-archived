package github_app

import (
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

// GetGithubAppClientFromRequest gets the github app installation id from the request and authenticates
// using it and the private key
func GetGithubAppClientFromRequest(config *server.Config, r *http.Request) (*github.Client, error) {
	// get installation id from context
	gai, _ := r.Context().Value(types.GithubAppInstallationScope).(*models.GithubAppInstallation)

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
