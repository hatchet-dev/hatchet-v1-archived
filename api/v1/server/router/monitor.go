package router

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/monitors"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// swagger:parameters getMonitor updateMonitor deleteMonitor
type monitorPathParams struct {
	// The team id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Team string `json:"team_id"`

	// The monitor id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Monitor string `json:"monitor_id"`
}

func NewMonitorRouteRegisterer(children ...*router.Registerer) *router.Registerer {
	return &router.Registerer{
		GetRoutes: GetMonitorRoutes,
		Children:  children,
	}
}

func GetMonitorRoutes(
	r chi.Router,
	config *server.Config,
	basePath *endpoint.Path,
	factory endpoint.APIEndpointFactory,
	children ...*router.Registerer,
) []*router.Route {
	routes := make([]*router.Route, 0)

	// POST /api/v1/teams/{team_id}/monitors -> monitors.NewMonitorCreateHandler
	// swagger:operation POST /api/v1/teams/{team_id}/monitors createMonitor
	//
	// ### Description
	//
	// Creates a new monitor.
	//
	// ---
	// produces:
	// - application/json
	// summary: Create Monitor
	// tags:
	// - Monitors
	// parameters:
	//   - name: team_id
	//   - in: body
	//     name: CreateMonitorRequest
	//     description: The monitor to create
	//     schema:
	//       $ref: '#/definitions/CreateMonitorRequest'
	// responses:
	//   '200':
	//     description: Successfully created the monitor
	//     schema:
	//       $ref: '#/definitions/CreateMonitorResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	monitorCreateEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/monitors",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
			},
		},
	)

	monitorCreateHandler := monitors.NewMonitorCreateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: monitorCreateEndpoint,
		Handler:  monitorCreateHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/monitors/{monitor_id} -> monitors.NewMonitorGetHandler
	// swagger:operation GET /api/v1/teams/{team_id}/monitors/{monitor_id} getMonitor
	//
	// ### Description
	//
	// Gets a monitor by id.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get Monitor
	// tags:
	// - Monitors
	// parameters:
	//   - name: team_id
	//   - name: monitor_id
	// responses:
	//   '200':
	//     description: Successfully got the monitor
	//     schema:
	//       $ref: '#/definitions/GetMonitorResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	getMonitorEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/monitors/{%s}", types.URLParamMonitorID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.MonitorScope,
			},
		},
	)

	getMonitorHandler := monitors.NewMonitorGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getMonitorEndpoint,
		Handler:  getMonitorHandler,
		Router:   r,
	})

	// POST /api/v1/teams/{team_id}/monitors/{monitor_id} -> monitors.NewMonitorUpdateHandler
	// swagger:operation POST /api/v1/teams/{team_id}/monitors/{monitor_id} updateMonitor
	//
	// ### Description
	//
	// Updates a monitor.
	//
	// ---
	// produces:
	// - application/json
	// summary: Update Monitor
	// tags:
	// - Monitors
	// parameters:
	//   - name: team_id
	//   - name: monitor_id
	//   - in: body
	//     name: UpdateMonitorRequest
	//     description: The monitor to update
	//     schema:
	//       $ref: '#/definitions/UpdateMonitorRequest'
	// responses:
	//   '200':
	//     description: Successfully updated the monitor
	//     schema:
	//       $ref: '#/definitions/UpdateMonitorResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	updateMonitorEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/monitors/{%s}", types.URLParamMonitorID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.MonitorScope,
			},
		},
	)

	updateMonitorHandler := monitors.NewMonitorUpdateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: updateMonitorEndpoint,
		Handler:  updateMonitorHandler,
		Router:   r,
	})

	// DELETE /api/v1/teams/{team_id}/monitors/{monitor_id} -> monitors.NewMonitorDeleteHandler
	// swagger:operation DELETE /api/v1/teams/{team_id}/monitors/{monitor_id} deleteMonitor
	//
	// ### Description
	//
	// Deletes a monitor.
	//
	// ---
	// produces:
	// - application/json
	// summary: Delete Monitor
	// tags:
	// - Monitors
	// parameters:
	//   - name: team_id
	//   - name: monitor_id
	// responses:
	//   '200':
	//     description: Successfully delete the monitor
	//     schema:
	//       $ref: '#/definitions/DeleteMonitorResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	deleteMonitorEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbDelete,
			Method: types.HTTPVerbDelete,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/monitors/{%s}", types.URLParamMonitorID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.MonitorScope,
			},
		},
	)

	deleteMonitorHandler := monitors.NewMonitorDeleteHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: deleteMonitorEndpoint,
		Handler:  deleteMonitorHandler,
		Router:   r,
	})

	policyGetEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			ContentType: "application/octet-stream",
			Verb:        types.APIVerbGet,
			Method:      types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/monitors/{%s}/policy", types.URLParamMonitorID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.MonitorScope,
				types.ModuleServiceAccountScope,
			},
		},
	)

	policyGetHandler := monitors.NewMonitorPolicyGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: policyGetEndpoint,
		Handler:  policyGetHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/monitors -> modules.NewMonitorListHandler
	// swagger:operation GET /api/v1/teams/{team_id}/monitors listMonitors
	//
	// ### Description
	//
	// Lists monitors for a given team.
	//
	// ---
	// produces:
	// - application/json
	// summary: List Monitors
	// tags:
	// - Monitors
	// responses:
	//   '200':
	//     description: Successfully listed monitors
	//     schema:
	//       $ref: '#/definitions/ListMonitorsResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	monitorsListEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/monitors",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
			},
		},
	)

	monitorsListHandler := monitors.NewMonitorListHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: monitorsListEndpoint,
		Handler:  monitorsListHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/monitor_results -> modules.NewMonitorListResultsHandler
	// swagger:operation GET /api/v1/teams/{team_id}/monitor_results listMonitorResults
	//
	// ### Description
	//
	// Lists monitor results for a given team, optionally filtered by module or monitor id.
	//
	// ---
	// produces:
	// - application/json
	// summary: List Monitor Results
	// tags:
	// - Monitors
	// responses:
	//   '200':
	//     description: Successfully listed monitor results
	//     schema:
	//       $ref: '#/definitions/ListMonitorResultsResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	resultsListEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/monitor_results",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
			},
		},
	)

	resultsListHandler := monitors.NewMonitorResultListHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: resultsListEndpoint,
		Handler:  resultsListHandler,
		Router:   r,
	})

	return routes
}
