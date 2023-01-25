package router

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/modules"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/terraform_state"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// swagger:parameters createModuleRun
type modulePathParams struct {
	// The team id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Team string `json:"team_id"`

	// The module id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Module string `json:"module_id"`
}

// swagger:parameters createTerraformState getTerraformState
type moduleRunPathParams struct {
	// The team id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Team string `json:"team_id"`

	// The module id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Module string `json:"module_id"`

	// The module run id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Run string `json:"module_run_id"`
}

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

	// POST /api/v1/teams/{team_id}/modules/{module_id}/runs -> modules.NewRunCreateHandler
	// swagger:operation POST /api/v1/teams/{team_id}/modules/{module_id}/runs createModuleRun
	//
	// ### Description
	//
	// Creates a new module run.
	//
	// ---
	// produces:
	// - application/json
	// summary: Create Module Run
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	// responses:
	//   '201':
	//     description: Successfully created the module
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	runCreateEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs", types.URLParamModuleID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
			},
		},
	)

	runCreateHandler := modules.NewRunCreateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: runCreateEndpoint,
		Handler:  runCreateHandler,
		Router:   r,
	})

	// POST /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate -> terraform_state
	// swagger:operation POST /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate createTerraformState
	//
	// ### Description
	//
	// Creates a `POST` request for Terraform state. **Should only be called by Terraform in automation.**
	//
	// ---
	// produces:
	// - application/json
	// summary: Create or Update Terraform State
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	//   - name: module_run_id
	// responses:
	//   '200':
	//     description: Successfully got the TF state
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	//   '423':
	//     description: Locked
	tfStateCreateEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/tfstate", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.BasicAuthUserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
			},
		},
	)

	tfStateCreateHandler := terraform_state.NewTerraformStateCreateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: tfStateCreateEndpoint,
		Handler:  tfStateCreateHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate -> terraform_state
	// swagger:operation POST /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate getTerraformState
	//
	// ### Description
	//
	// Creates a `GET` request for Terraform state. **Should only be called by Terraform in automation.**
	//
	// ---
	// produces:
	// - application/json
	// summary: Create or Update Terraform State
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	//   - name: module_run_id
	// responses:
	//   '200':
	//     description: Successfully got the TF state
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	tfStateGetEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/tfstate", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.BasicAuthUserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
			},
		},
	)

	tfStateGetHandler := terraform_state.NewTerraformStateGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: tfStateGetEndpoint,
		Handler:  tfStateGetHandler,
		Router:   r,
	})

	// LOCK /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate -> terraform_state
	// swagger:operation LOCK /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate lockTerraformState
	//
	// ### Description
	//
	// Creates a `LOCK` request for Terraform state. **Should only be called by Terraform in automation.**
	//
	// ---
	// produces:
	// - application/json
	// summary: Lock Terraform State
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	//   - name: module_run_id
	// responses:
	//   '200':
	//     description: Successfully locked
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	//   '423':
	//     description: Locked
	tfStateLockEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbLock,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/tfstate", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.BasicAuthUserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
			},
		},
	)

	tfStateLockHandler := terraform_state.NewTerraformStateLockHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: tfStateLockEndpoint,
		Handler:  tfStateLockHandler,
		Router:   r,
	})

	// UNLOCK /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate -> terraform_state
	// swagger:operation UNLOCK /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate unlockTerraformState
	//
	// ### Description
	//
	// Creates an `UNLOCK` request for Terraform state. **Should only be called by Terraform in automation.**
	//
	// ---
	// produces:
	// - application/json
	// summary: Unlock Terraform State
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	//   - name: module_run_id
	// responses:
	//   '200':
	//     description: Successfully unlocked
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	tfStateUnlockEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbUnlock,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/tfstate", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.BasicAuthUserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
			},
		},
	)

	tfStateUnlockHandler := terraform_state.NewTerraformStateUnlockHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: tfStateUnlockEndpoint,
		Handler:  tfStateUnlockHandler,
		Router:   r,
	})

	return routes
}