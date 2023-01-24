package router

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/github_app"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// swagger:parameters listGithubRepos
type githubAppInstallationPathParams struct {
	// The github app installation id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	GithubAppInstallation string `json:"github_app_installation_id"`
}

func NewGithubAppRouteRegisterer(children ...*router.Registerer) *router.Registerer {
	return &router.Registerer{
		GetRoutes: GetGithubAppRoutes,
		Children:  children,
	}
}

func GetGithubAppRoutes(
	r chi.Router,
	config *server.Config,
	basePath *endpoint.Path,
	factory endpoint.APIEndpointFactory,
	children ...*router.Registerer,
) []*router.Route {
	routes := make([]*router.Route, 0)

	// GET /api/v1/oauth/github_app -> github_app.NewGithubAppOAuthStartHandler
	// swagger:operation GET /api/v1/oauth/github_app startGithubAppOAuth
	//
	// ### Description
	//
	// Starts the OAuth flow to authenticate with a Github App.
	//
	// ---
	// produces:
	// - application/json
	// summary: Start Github App OAuth
	// tags:
	// - Github Apps
	// responses:
	//   '302':
	//     description: Successfully triggered Github App oauth
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	//   '405':
	//     description: This endpoint is not supported on this Hatchet instance.
	//     schema:
	//       $ref: '#/definitions/APIErrorNotSupportedExample'
	startOAuthEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/oauth/github_app",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	startOAuthHandler := github_app.NewGithubAppOAuthStartHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: startOAuthEndpoint,
		Handler:  startOAuthHandler,
		Router:   r,
	})

	// GET /api/v1/oauth/github_app/callback -> github_app.NewGithubAppOAuthCallbackHandler
	// swagger:operation GET /api/v1/oauth/github_app/callback finishGithubAppOAuth
	//
	// ### Description
	//
	// Finishes the OAuth flow to authenticate with a Github App.
	//
	// ---
	// produces:
	// - application/json
	// summary: Start Github App OAuth
	// tags:
	// - Github Apps
	// responses:
	//   '302':
	//     description: Successfully authenticated OR error state
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	//   '405':
	//     description: This endpoint is not supported on this Hatchet instance.
	//     schema:
	//       $ref: '#/definitions/APIErrorNotSupportedExample'
	callbackOAuthEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/oauth/github_app/callback",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	callbackOAuthHandler := github_app.NewGithubAppOAuthCallbackHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: callbackOAuthEndpoint,
		Handler:  callbackOAuthHandler,
		Router:   r,
	})

	// GET /api/v1/github_app/install -> github_app.NewGithubAppOAuthInstallHandler
	// swagger:operation GET /api/v1/github_app/install installGithubApp
	//
	// ### Description
	//
	// Redirects the user to Github to install the Github App.
	//
	// ---
	// produces:
	// - application/json
	// summary: Install Github App
	// tags:
	// - Github Apps
	// responses:
	//   '302':
	//     description: Successfully redirected to Github app installation
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	//   '405':
	//     description: This endpoint is not supported on this Hatchet instance.
	//     schema:
	//       $ref: '#/definitions/APIErrorNotSupportedExample'
	installEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/github_app/install",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	installHandler := github_app.NewGithubAppOAuthInstallHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: installEndpoint,
		Handler:  installHandler,
		Router:   r,
	})

	// POST /api/v1/webhooks/github_app -> github_app.NewGithubAppWebhookHandler
	// swagger:operation POST /api/v1/webhooks/github_app githubAppWebhook
	//
	// ### Description
	//
	// Implements a Github App webhook.
	//
	// ---
	// produces:
	// - application/json
	// summary: Github App Webhook
	// tags:
	// - Github Apps
	// responses:
	//   '200':
	//     description: Successfully processed app webhook
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	//   '405':
	//     description: This endpoint is not supported on this Hatchet instance.
	//     schema:
	//       $ref: '#/definitions/APIErrorNotSupportedExample'
	webhookEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/webhooks/github_app",
			},
			Scopes: []types.PermissionScope{},
		},
	)

	webhookHandler := github_app.NewGithubAppWebhookHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: webhookEndpoint,
		Handler:  webhookHandler,
		Router:   r,
	})

	// GET /api/v1/github_app/{github_app_installation_id}/repos -> github_app.NewListGithubReposHandler
	// swagger:operation GET /api/v1/github_app/{github_app_installation_id}/repos listGithubRepos
	//
	// ### Description
	//
	// Lists the Github repos that the github app installation has access to.
	//
	// ---
	// produces:
	// - application/json
	// summary: List Github Repos
	// tags:
	// - Github Apps
	// responses:
	//   '200':
	//     description: Successfully listed Github repositories
	//     schema:
	//       $ref: '#/definitions/ListGithubReposResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	listReposEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/github_app/{%s}/repos", types.URLParamGithubAppInstallationID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.GithubAppInstallationScope,
			},
		},
	)

	listReposHandler := github_app.NewListGithubReposHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: listReposEndpoint,
		Handler:  listReposHandler,
		Router:   r,
	})

	// GET /api/v1/github_app/{github_app_installation_id}/repos/{github_repo_owner}/{github_repo_name}/branches -> github_app.NewListGithubRepoBranchesHandler
	// swagger:operation GET /api/v1/github_app/{github_app_installation_id}/repos/{github_repo_owner}/{github_repo_name}/branches listGithubRepoBranches
	//
	// ### Description
	//
	// Lists the Github repo branches.
	//
	// ---
	// produces:
	// - application/json
	// summary: List Github Branches
	// tags:
	// - Github Apps
	// responses:
	//   '200':
	//     description: Successfully listed Github repo branches
	//     schema:
	//       $ref: '#/definitions/ListGithubRepoBranchesResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	listBranchesEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/github_app/{%s}/repos/{%s}/{%s}/branches", types.URLParamGithubAppInstallationID, types.URLParamGithubRepoOwner, types.URLParamGithubRepoName),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.GithubAppInstallationScope,
			},
		},
	)

	listBranchesHandler := github_app.NewListGithubRepoBranchesHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: listBranchesEndpoint,
		Handler:  listBranchesHandler,
		Router:   r,
	})

	return routes
}
