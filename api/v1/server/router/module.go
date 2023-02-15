package router

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/logs"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/modules"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/terraform_state"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// swagger:parameters createModuleRun getModuleTarballURL listModuleRuns getCurrentModuleValues getModule updateModule deleteModule
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

// swagger:parameters createTerraformState getTerraformState lockTerraformState unlockTerraformState createTerraformPlan uploadTerraformPlan getTerraformPlan getTerraformPlanBySHA finalizeModuleRun getModuleRun getModuleRunLogs
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

// swagger:parameters getModuleEnvVars
type moduleEnvVarsPathParams struct {
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

	// The module env vars version id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	ModuleEnvVar string `json:"module_env_vars_id"`
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

	// POST /api/v1/teams/{team_id}/modules/{module_id} -> modules.NewModuleUpdateHandler
	// swagger:operation POST /api/v1/teams/{team_id}/modules/{module_id} updateModule
	//
	// ### Description
	//
	// Updates a module.
	//
	// ---
	// produces:
	// - application/json
	// summary: Update Module Run
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	//   - in: body
	//     name: UpdateModuleRequest
	//     description: The module params to update
	//     schema:
	//       $ref: '#/definitions/UpdateModuleRequest'
	// responses:
	//   '200':
	//     description: Successfully updated the module
	//     schema:
	//       $ref: '#/definitions/UpdateModuleResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	updateModuleEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}", types.URLParamModuleID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
			},
		},
	)

	updateModuleHandler := modules.NewModuleUpdateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: updateModuleEndpoint,
		Handler:  updateModuleHandler,
		Router:   r,
	})

	// POST /api/v1/teams/{team_id}/modules/{module_id}/force_unlock -> modules.NewModuleForceUnlockHandler
	// swagger:operation POST /api/v1/teams/{team_id}/modules/{module_id}/force_unlock forceUnlockModule
	//
	// ### Description
	//
	// Force unlocks a module.
	//
	// ---
	// produces:
	// - application/json
	// summary: Unlock Module
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	// responses:
	//   '200':
	//     description: Successfully unlocked the module
	//     schema:
	//       $ref: '#/definitions/ForceUnlockModuleResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	forceUnlockModuleEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/force_unlock", types.URLParamModuleID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
			},
		},
	)

	forceUnlockModuleHandler := modules.NewModuleForceUnlockHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: forceUnlockModuleEndpoint,
		Handler:  forceUnlockModuleHandler,
		Router:   r,
	})

	// DELETE /api/v1/teams/{team_id}/modules/{module_id} -> modules.NewModuleDeleteHandler
	// swagger:operation DELETE /api/v1/teams/{team_id}/modules/{module_id} deleteModule
	//
	// ### Description
	//
	// Deletes a module.
	//
	// ---
	// produces:
	// - application/json
	// summary: Delete Module
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	// responses:
	//   '200':
	//     description: Successfully deleted the module
	//     schema:
	//       $ref: '#/definitions/DeleteModuleResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	deleteModuleEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbDelete,
			Method: types.HTTPVerbDelete,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}", types.URLParamModuleID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
			},
		},
	)

	deleteModuleHandler := modules.NewModuleDeleteHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: deleteModuleEndpoint,
		Handler:  deleteModuleHandler,
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

	// GET /api/v1/teams/{team_id}/modules/{module_id} -> modules.NewModuleGetHandler
	// swagger:operation GET /api/v1/teams/{team_id}/modules/{module_id} getModule
	//
	// ### Description
	//
	// Gets a module by module id.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get module
	// tags:
	// - Modules
	// responses:
	//   '200':
	//     description: Successfully got the module
	//     schema:
	//       $ref: '#/definitions/GetModuleResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	getModuleEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}", types.URLParamModuleID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
			},
		},
	)

	getModuleHandler := modules.NewModuleGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getModuleEndpoint,
		Handler:  getModuleHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/modules/{module_id}/runs -> modules.NewModuleRunsListHandler
	// swagger:operation GET /api/v1/teams/{team_id}/modules/{module_id}/runs listModuleRuns
	//
	// ### Description
	//
	// Lists module runs for a given module id.
	//
	// ---
	// produces:
	// - application/json
	// summary: List Module Runs
	// tags:
	// - Modules
	// responses:
	//   '200':
	//     description: Successfully listed the module runs
	//     schema:
	//       $ref: '#/definitions/ListModuleRunsResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	listRunsEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbList,
			Method: types.HTTPVerbGet,
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

	listRunsHandler := modules.NewModuleRunsListHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: listRunsEndpoint,
		Handler:  listRunsHandler,
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

	// POST /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/plan -> terraform_state.NewTerraformPlanCreateHandler
	// swagger:operation POST /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/plan createTerraformPlan
	//
	// ### Description
	//
	// Creates a `POST` request for a Terraform plan. **Should only be called by Terraform in automation.**
	//
	// ---
	// produces:
	// - application/json
	// summary: Create Terraform plan
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	//   - name: module_run_id
	//   - in: body
	//     required: true
	//     name: CreateTerraformPlanRequest
	//     description: The terraform plan contents
	//     schema:
	//       $ref: '#/definitions/CreateTerraformPlanRequest'
	// responses:
	//   '200':
	//     description: Successfully created the TF plan
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
	tfPlanCreateEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/plan", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
				types.ModuleServiceAccountScope,
			},
		},
	)

	tfPlanCreateHandler := terraform_state.NewTerraformPlanCreateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: tfPlanCreateEndpoint,
		Handler:  tfPlanCreateHandler,
		Router:   r,
	})

	tfPlanUploadEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/plan/zip", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
				types.ModuleServiceAccountScope,
			},
		},
	)

	tfPlanUploadHandler := terraform_state.NewTerraformPlanUploadZIPHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: tfPlanUploadEndpoint,
		Handler:  tfPlanUploadHandler,
		Router:   r,
	})

	tfPlanGetEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			ContentType: "application/octet-stream",
			Verb:        types.APIVerbGet,
			Method:      types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/plan", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
				types.ModuleServiceAccountScope,
			},
		},
	)

	tfPlanGetHandler := terraform_state.NewTerraformPlanGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: tfPlanGetEndpoint,
		Handler:  tfPlanGetHandler,
		Router:   r,
	})

	tfPlanGetBySHAEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			ContentType: "application/octet-stream",
			Verb:        types.APIVerbGet,
			Method:      types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/plan/sha", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
				types.ModuleServiceAccountScope,
			},
		},
	)

	tfPlanGetBySHAHandler := terraform_state.NewTerraformPlanGetBySHAHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: tfPlanGetBySHAEndpoint,
		Handler:  tfPlanGetBySHAHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/modules/{module_id}/tarball_url -> modules.NewModuleTarballURLGetHandler
	// swagger:operation GET /api/v1/teams/{team_id}/modules/{module_id}/tarball_url getModuleTarballURL
	//
	// ### Description
	//
	// Gets the Github tarball URL for the module.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get Module Tarball URL
	// tags:
	// - Modules
	// responses:
	//   '200':
	//     description: Successfully got tarball url
	//     schema:
	//       $ref: '#/definitions/GetModuleTarballURLResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	getTarballURLEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/tarball_url", types.URLParamModuleID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleServiceAccountScope,
			},
		},
	)

	getTarballURLHandler := modules.NewModuleTarballURLGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getTarballURLEndpoint,
		Handler:  getTarballURLHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/modules/{module_id}/values/current -> modules.NewModuleValuesGetCurrentHandler
	// swagger:operation GET /api/v1/teams/{team_id}/modules/{module_id}/values/current getCurrentModuleValues
	//
	// ### Description
	//
	// Gets the current module values for the given module, by github reference or SHA.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get Current Module Values
	// tags:
	// - Modules
	// responses:
	//   '200':
	//     description: Successfully got module values
	//     schema:
	//       $ref: '#/definitions/GetModuleValuesCurrentResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	getCurrentModuleValuesEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/values/current", types.URLParamModuleID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleServiceAccountScope,
			},
		},
	)

	getCurrentModuleValuesHandler := modules.NewModuleValuesCurrentGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getCurrentModuleValuesEndpoint,
		Handler:  getCurrentModuleValuesHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/modules/{module_id}/values/{module_values_id} -> modules.NewModuleValuesGetHandler
	// swagger:operation GET /api/v1/teams/{team_id}/modules/{module_id}/values/{module_values_id} getModuleValues
	//
	// ### Description
	//
	// Gets the current module values by ID.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get Module Values
	// tags:
	// - Modules
	// responses:
	//   '200':
	//     description: Successfully got module values
	//     schema:
	//       $ref: '#/definitions/GetModuleValuesResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	getModuleValuesEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/values/{%s}", types.URLParamModuleID, types.URLParamModuleValuesID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleValuesScope,
			},
		},
	)

	getModuleValuesHandler := modules.NewModuleValuesGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getModuleValuesEndpoint,
		Handler:  getModuleValuesHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/modules/{module_id}/env_vars/{module_env_vars_id} -> modules.NewModuleEnvVarsGetHandler
	// swagger:operation GET /api/v1/teams/{team_id}/modules/{module_id}/env_vars/{module_env_vars_id} getModuleEnvVars
	//
	// ### Description
	//
	// Gets the module env vars version by id.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get Module Env Vars
	// tags:
	// - Modules
	// responses:
	//   '200':
	//     description: Successfully got module env vars
	//     schema:
	//       $ref: '#/definitions/GetModuleEnvVarsVersionResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	getModuleEnvVarsEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/env_vars/{%s}", types.URLParamModuleID, types.URLParamModuleEnvVarsID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleEnvVarScope,
			},
		},
	)

	getModuleEnvVarsHandler := modules.NewModuleEnvVarsGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getModuleEnvVarsEndpoint,
		Handler:  getModuleEnvVarsHandler,
		Router:   r,
	})

	// POST /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/finalize -> modules.NewModuleRunFinalizeHandler
	// swagger:operation POST /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/finalize finalizeModuleRun
	//
	// ### Description
	//
	// Updates a module run with a finalized status.
	//
	// ---
	// produces:
	// - application/json
	// summary: Finalize module run
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	//   - name: module_run_id
	//   - in: body
	//     required: true
	//     name: FinalizeModuleRunRequest
	//     description: The module run status to update
	//     schema:
	//       $ref: '#/definitions/FinalizeModuleRunRequest'
	// responses:
	//   '200':
	//     description: Successfully updated the module run
	//     schema:
	//       $ref: '#/definitions/FinalizeModuleRunResponse'
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
	finalizeRunEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/finalize", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
				types.ModuleServiceAccountScope,
			},
		},
	)

	finalizeRunHandler := modules.NewModuleRunFinalizeHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: finalizeRunEndpoint,
		Handler:  finalizeRunHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id} -> modules.NewModuleRunGetHandler
	// swagger:operation GET /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id} getModuleRun
	//
	// ### Description
	//
	// Gets a module run.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get module run
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	//   - name: module_run_id
	// responses:
	//   '200':
	//     description: Successfully got the module run
	//     schema:
	//       $ref: '#/definitions/GetModuleRunResponse'
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
	getRunEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
			},
		},
	)

	getRunHandler := modules.NewModuleRunGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getRunEndpoint,
		Handler:  getRunHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/plan_summary -> modules.NewModuleGetPlanSummaryHandler
	// swagger:operation GET /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/plan_summary getModuleRunPlanSummary
	//
	// ### Description
	//
	// Gets the plan summary for a module run.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get plan summary
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	//   - name: module_run_id
	// responses:
	//   '200':
	//     description: Successfully got the plan summary
	//     schema:
	//       $ref: '#/definitions/GetModulePlanSummaryResponse'
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
	getRunPlanSummaryEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/plan_summary", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
			},
		},
	)

	getRunPlanSummaryHandler := modules.NewModuleGetPlanSummaryHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getRunPlanSummaryEndpoint,
		Handler:  getRunPlanSummaryHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/logs -> logs.NewLogGetHandler
	// swagger:operation GET /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/logs getModuleRunLogs
	//
	// ### Description
	//
	// Gets the logs for a module run.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get logs
	// tags:
	// - Modules
	// parameters:
	//   - name: team_id
	//   - name: module_id
	//   - name: module_run_id
	// responses:
	//   '200':
	//     description: Successfully got the logs
	//     schema:
	//       $ref: '#/definitions/GetLogsResponse'
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
	getLogsEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/modules/{%s}/runs/{%s}/logs", types.URLParamModuleID, types.URLParamModuleRunID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.ModuleScope,
				types.ModuleRunScope,
			},
		},
	)

	getLogsHandler := logs.NewLogGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getLogsEndpoint,
		Handler:  getLogsHandler,
		Router:   r,
	})

	return routes
}
