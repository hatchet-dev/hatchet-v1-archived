package github_app

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/authn"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"golang.org/x/oauth2"
)

type GithubAppOAuthStartHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewGithubAppOAuthStartHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &GithubAppOAuthStartHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (g *GithubAppOAuthStartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state, err := authn.SaveOAuthState(w, r, g.Config())

	if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	http.Redirect(w, r, g.Config().GithubApp.AuthCodeURL(state, oauth2.AccessTypeOffline), 302)
}
