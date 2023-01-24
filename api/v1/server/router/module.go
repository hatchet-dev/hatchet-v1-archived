package router

import (
	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/modules"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

func NewModuleRouteRegisterer(children ...*router.Registerer) *router.Registerer {
	return &router.Registerer{
		GetRoutes: GetModuleRoutes,
		Children:  children,
	}
}

func GetModuleRoutes(
	r chi.Router,
	config *server.Config,
	basePath *endpoint.Path,
	factory endpoint.APIEndpointFactory,
	children ...*router.Registerer,
) []*router.Route {
	routes := make([]*router.Route, 0)

	// POST /api/v1/teams/{team_id}/modules -> modules.NewModuleCreateHandler
	// swagger:operation POST /api/v1/teams/{team_id}/modules createModule
	//
	// ### Description
	//
	// Creates a new module.
	//
	// ---
	// produces:
	// - application/json
	// summary: Create Module
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - in: body
	//     name: CreateModuleRequest
	//     description: The module to create
	//     schema:
	//       $ref: '#/definitions/CreateModuleRequest'
	// responses:
	//   '201':
	//     description: Successfully created the module
	//     schema:
	//       $ref: '#/definitions/CreateModuleResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	moduleCreateEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/modules",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
			},
		},
	)

	moduleCreateHandler := modules.NewModuleCreateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: moduleCreateEndpoint,
		Handler:  moduleCreateHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/modules -> modules.NewModuleListHandler
	// swagger:operation GET /api/v1/teams/{team_id}/modules listModules
	//
	// ### Description
	//
	// Lists modules for a given team.
	//
	// ---
	// produces:
	// - application/json
	// summary: List Modules
	// tags:
	// - Modules
	// responses:
	//   '200':
	//     description: Successfully listed modules
	//     schema:
	//       $ref: '#/definitions/ListModulesResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	modulesListEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/modules",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
			},
		},
	)

	modulesListHandler := modules.NewModuleListHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: modulesListEndpoint,
		Handler:  modulesListHandler,
		Router:   r,
	})

	return routes
}
