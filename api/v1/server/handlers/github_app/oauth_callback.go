package github_app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/authn"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"golang.org/x/oauth2"
)

type GithubAppOAuthCallbackHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewGithubAppOAuthCallbackHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &GithubAppOAuthCallbackHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (g *GithubAppOAuthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)

	validate, isOAuthTriggered, err := authn.ValidateOAuthState(w, r, g.Config())

	if !validate || err != nil {
		if err != nil {
			g.HandleAPIError(w, r, apierrors.NewErrForbidden(err))
		} else {
			g.HandleAPIError(w, r, apierrors.NewErrForbidden(fmt.Errorf("could not validate state parameter")))
		}

		return
	}

	token, err := g.Config().GithubApp.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))

	if err != nil || !token.Valid() {
		http.Redirect(w, r, "/", 302)

		return
	}

	sharedFields := &models.SharedOAuthFields{
		ClientID:     []byte(g.Config().GithubApp.ClientID),
		AccessToken:  []byte(token.AccessToken),
		RefreshToken: []byte(token.RefreshToken),
		Expiry:       token.Expiry,
		UserID:       user.ID,
	}

	ghClient := github.NewClient(g.Config().GithubApp.Client(oauth2.NoContext, token))

	githubUser, _, err := ghClient.Users.Get(context.Background(), "")

	if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	githubAppOAuth := &models.GithubAppOAuth{
		SharedOAuthFields: sharedFields,
		GithubUserID:      *githubUser.ID,
	}

	// if the user already has an oauth integration, we update it when they re-auth
	foundGHAppOAuth, err := g.Repo().GithubAppOAuth().ReadGithubAppOAuthByUserID(user.ID)

	if err != nil && errors.Is(err, repository.RepositoryErrorNotFound) {
		githubAppOAuth, err = g.Repo().GithubAppOAuth().CreateGithubAppOAuth(githubAppOAuth)

		if err != nil {
			g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}
	} else if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	} else if foundGHAppOAuth != nil {
		foundGHAppOAuth.SharedOAuthFields = sharedFields
		foundGHAppOAuth.GithubUserID = *githubUser.ID

		foundGHAppOAuth, err = g.Repo().GithubAppOAuth().UpdateGithubAppOAuth(foundGHAppOAuth)

		if err != nil {
			g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}
	}

	if isOAuthTriggered {
		http.Redirect(w, r, fmt.Sprintf("https://github.com/apps/%s/installations/new", g.Config().GithubApp.AppName), 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}
