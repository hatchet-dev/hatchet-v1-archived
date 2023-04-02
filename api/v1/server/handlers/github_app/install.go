package github_app

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

type GithubAppOAuthInstallHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewGithubAppOAuthInstallHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &GithubAppOAuthInstallHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (g *GithubAppOAuthInstallHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	githubProvider, reqErr := GetGithubProvider(g.Config())

	if reqErr != nil {
		g.HandleAPIError(w, r, reqErr)
	}

	http.Redirect(w, r, fmt.Sprintf("https://github.com/apps/%s/installations/new", githubProvider.GetGithubAppConfig().GetAppName()), 302)
}
