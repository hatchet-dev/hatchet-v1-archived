package github_app

import (
	"errors"
	"net/http"

	"github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type GithubAppWebhookHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewGithubAppWebhookHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &GithubAppWebhookHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (g *GithubAppWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ghApp, reqErr := GetGithubAppConfig(g.Config())

	if reqErr != nil {
		g.HandleAPIError(w, r, reqErr)
		return
	}

	payload, err := github.ValidatePayload(r, []byte(ghApp.GetWebhookSecret()))

	if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)

	if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	switch e := event.(type) {
	case *github.InstallationEvent:
		if *e.Action == "created" || *e.Action == "added" {
			// make sure the sender exists in the database
			gao, err := g.Repo().GithubAppOAuth().ReadGithubAppOAuthByGithubUserID(*e.Sender.ID)

			if err != nil {
				g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
				return
			}

			_, err = g.Repo().GithubAppInstallation().ReadGithubAppInstallationByInstallationAndAccountID(*e.Installation.ID, *e.Installation.Account.ID)

			if err != nil && errors.Is(err, repository.RepositoryErrorNotFound) {
				// insert account/installation pair into database
				_, err := g.Repo().GithubAppInstallation().CreateGithubAppInstallation(&models.GithubAppInstallation{
					GithubAppOAuthID:        gao.ID,
					AccountName:             *e.Installation.Account.Login,
					AccountAvatarURL:        *e.Installation.Account.AvatarURL,
					AccountID:               *e.Installation.Account.ID,
					InstallationID:          *e.Installation.ID,
					InstallationSettingsURL: *e.Installation.HTMLURL,
				})

				if err != nil {
					g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
					return
				}

				return
			} else if err != nil {
				g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
				return
			}
		}

		if *e.Action == "deleted" {
			gai, err := g.Repo().GithubAppInstallation().ReadGithubAppInstallationByInstallationAndAccountID(*e.Installation.ID, *e.Installation.Account.ID)

			if err != nil {
				g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
				return
			}

			gai, err = g.Repo().GithubAppInstallation().DeleteGithubAppInstallation(gai)

			if err != nil {
				g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
				return
			}
		}
	}
}
