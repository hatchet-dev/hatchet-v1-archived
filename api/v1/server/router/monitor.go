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

	return routes
}
