package router

import (
	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/pats"
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
		//     schema:
		//       $ref: '#/definitions/CreateUserResponse'
		//   '400':
		//     description: A malformed or bad request
		//     schema:
		//       $ref: '#/definitions/APIErrors'
		//   '403':
		//     description: Forbidden
		//     schema:
		//       $ref: '#/definitions/APIErrors'
		//   '405':
		//     description: This endpoint is not supported on this Hatchet instance.
		//     schema:
		//       $ref: '#/definitions/APIErrors'
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

	// GET /api/v1/users/current -> user.UserGetCurrentHandler
	// swagger:operation GET /api/v1/users/current getCurrentUser
	//
	// Retrieves the current user object based on the data passed in auth.
	//
	// ---
	// produces:
	// - application/json
	// summary: Retrieve the current user.
	// tags:
	// - Users
	// responses:
	//   '200':
	//     description: Successfully got the user
	//     schema:
	//       $ref: '#/definitions/GetUserResponse'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrors'
	getUserCurrentEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	getUserCurrentHandler := users.NewUserGetCurrentHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getUserCurrentEndpoint,
		Handler:  getUserCurrentHandler,
		Router:   r,
	})

	// DELETE /api/v1/users/current -> user.UserDeleteCurrentHandler
	// swagger:operation DELETE /api/v1/users/current deleteCurrentUser
	//
	// Deletes the current user.
	//
	// ---
	// produces:
	// - application/json
	// summary: Delete the current user.
	// tags:
	// - Users
	// responses:
	//   '202':
	//     description: Successfully triggered user deletion
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrors'
	deleteUserCurrentEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbDelete,
			Method: types.HTTPVerbDelete,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	deleteUserCurrentHandler := users.NewUserDeleteCurrentHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: deleteUserCurrentEndpoint,
		Handler:  deleteUserCurrentHandler,
		Router:   r,
	})

	// POST /api/v1/users/current/settings/pat -> user.NewPATCreateHandler
	// swagger:operation POST /api/v1/users/current/settings/pat createPersonalAccessToken
	//
	// Creates a new personal access token for a user.
	//
	// ---
	// produces:
	// - application/json
	// summary: Create a new user
	// tags:
	// - Users
	// parameters:
	//   - in: body
	//     name: CreatePATRequest
	//     description: The user to create
	//     schema:
	//       $ref: '#/definitions/CreatePATRequest'
	// responses:
	//   '201':
	//     description: Successfully created the user
	//     schema:
	//       $ref: '#/definitions/CreatePATResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrors'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrors'
	//   '405':
	//     description: This endpoint is not supported on this Hatchet instance.
	//     schema:
	//       $ref: '#/definitions/APIErrors'
	createPATEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current/settings/pat",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	createPATHandler := pats.NewPATCreateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: createPATEndpoint,
		Handler:  createPATHandler,
		Router:   r,
	})

	return routes
}
