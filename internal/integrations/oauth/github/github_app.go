package github

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/hatchet-dev/hatchet/internal/integrations/oauth"
	"golang.org/x/oauth2"
)

const (
	GithubAuthURL  string = "https://github.com/login/oauth/authorize"
	GithubTokenURL string = "https://github.com/login/oauth/access_token"
)

type GithubAppConf struct {
	oauth2.Config

	AppName       string
	WebhookSecret string
	Secret        []byte
	AppID         int64
}

func NewGithubAppConf(cfg *oauth.Config, appName, appSecretPath, appWebhookSecret, appID string) (*GithubAppConf, error) {
	intAppID, err := strconv.ParseInt(appID, 10, 64)

	if err != nil {
		return nil, err
	}

	appSecret, err := ioutil.ReadFile(appSecretPath)

	if err != nil {
		return nil, fmt.Errorf("could not read github app secret: %s", err)
	}

	return &GithubAppConf{
		AppName:       appName,
		WebhookSecret: appWebhookSecret,
		Secret:        appSecret,
		AppID:         intAppID,
		Config: oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  GithubAuthURL,
				TokenURL: GithubTokenURL,
			},
			RedirectURL: cfg.BaseURL + "/api/v1/oauth/github_app/callback",
			Scopes:      cfg.Scopes,
		},
	}, nil
}
