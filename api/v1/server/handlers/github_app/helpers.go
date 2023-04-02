package github_app

import (
	"fmt"
	"net/http"

	githubsdk "github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs/github"
	"github.com/hatchet-dev/hatchet/internal/models"
)

func GetGithubProvider(config *server.Config) (res github.GithubVCSProvider, reqErr apierrors.RequestError) {
	vcsFact, exists := config.VCSProviders[vcs.VCSRepositoryKindGithub]

	if !exists {
		return res, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        1406,
			Description: fmt.Sprintf("No Github app set up on this Hatchet instance."),
		}, http.StatusNotAcceptable)
	}

	res, err := github.ToGithubVCSProviderFactory(vcsFact)

	if err != nil {
		return res, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        1406,
			Description: fmt.Sprintf("Github app is improperly set up on this Hatchet instance."),
		}, http.StatusNotAcceptable)
	}

	return res, nil
}

func GetGithubAppConfig(config *server.Config) (*github.GithubAppConf, apierrors.RequestError) {
	githubFact, reqErr := GetGithubProvider(config)

	if reqErr != nil {
		return nil, reqErr
	}

	return githubFact.GetGithubAppConfig(), nil
}

// GetGithubAppClientFromRequest gets the github app installation id from the request and authenticates
// using it and the private key
func GetGithubAppClientFromRequest(config *server.Config, r *http.Request) (*githubsdk.Client, apierrors.RequestError) {
	// get installation id from context
	gai, _ := r.Context().Value(types.GithubAppInstallationScope).(*models.GithubAppInstallation)

	githubFact, reqErr := GetGithubProvider(config)

	if reqErr != nil {
		return nil, reqErr
	}

	res, err := githubFact.GetGithubAppConfig().GetGithubClient(gai.InstallationID)

	if err != nil {
		return nil, apierrors.NewErrInternal(err)
	}

	return res, nil
}
