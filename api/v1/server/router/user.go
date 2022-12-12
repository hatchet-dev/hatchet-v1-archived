package router

import (
	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/users"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

func NewUserScopedRegisterer(children ...*router.Registerer) *router.Registerer {
	return &router.Registerer{
		GetRoutes: GetUserScopedRoutes,
		Children:  children,
	}
}

func GetUserScopedRoutes(
	r chi.Router,
	config *server.Config,
	basePath *endpoint.Path,
	factory endpoint.APIEndpointFactory,
	children ...*router.Registerer,
) []*router.Route {
	routes := make([]*router.Route, 0)

	if config.AuthConfig.BasicAuthEnabled {
		// POST /api/v1/users -> user.NewCreateUserHandler
		// swagger:operation POST /api/v1/users createUser
		//
		// Creates a new user via email and password-based authentication. This endpoint is only registered if the
		// environment variable `BASIC_AUTH_ENABLED` is set.
		//
		// ---
		// produces:
		// - application/json
		// summary: Create a new user
		// tags:
		// - Users
		// parameters:
		//   - in: body
		//     name: CreateUserRequest
		//     description: The user to create
		//     schema:
		//       $ref: '#/definitions/CreateUserRequest'
		// responses:
		//   '201':
		//     description: Successfully created the user
		//   '400':
		//     description: A malformed or bad request
		//   '403':
		//     description: Forbidden
		//   '405':
		//     description: This endpoint is not supported on this Hatchet instance.
		createUserEndpoint := factory.NewAPIEndpoint(
			&endpoint.EndpointMetadata{
				Verb:   types.APIVerbCreate,
				Method: types.HTTPVerbPost,
				Path: &endpoint.Path{
					Parent:       basePath,
					RelativePath: "/users",
				},
				Scopes: []types.PermissionScope{},
			},
		)

		createUserHandler := users.NewUserCreateHandler(
			config,
			factory.GetDecoderValidator(),
			factory.GetResultWriter(),
		)

		routes = append(routes, &router.Route{
			Endpoint: createUserEndpoint,
			Handler:  createUserHandler,
			Router:   r,
		})
	}

	return routes
}
