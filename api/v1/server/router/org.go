package router

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/orgs"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// swagger:parameters getOrganization
type orgPathParams struct {
	// The org id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Org string `json:"org_id"`
}

func NewOrgRouteRegisterer(children ...*router.Registerer) *router.Registerer {
	return &router.Registerer{
		GetRoutes: GetOrganizationRoutes,
		Children:  children,
	}
}

func GetOrganizationRoutes(
	r chi.Router,
	config *server.Config,
	basePath *endpoint.Path,
	factory endpoint.APIEndpointFactory,
	children ...*router.Registerer,
) []*router.Route {
	routes := make([]*router.Route, 0)

	// POST /api/v1/organizations -> orgs.NewOrgCreateHandler
	// swagger:operation POST /api/v1/organizations createOrganization
	//
	// Creates a new organization, with the authenticated user set as the organization owner.
	//
	// ---
	// produces:
	// - application/json
	// summary: Create a new organization
	// tags:
	// - Organizations
	// parameters:
	//   - in: body
	//     name: CreateOrganizationRequest
	//     description: The organization to create
	//     schema:
	//       $ref: '#/definitions/CreateOrganizationRequest'
	// responses:
	//   '201':
	//     description: Successfully created the organization
	//     schema:
	//       $ref: '#/definitions/CreateOrganizationResponse'
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
	createOrgEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/organizations",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	createOrgHandler := orgs.NewOrgCreateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: createOrgEndpoint,
		Handler:  createOrgHandler,
		Router:   r,
	})

	// GET /api/v1/organizations/{org_id} -> orgs.NewOrgGetHandler
	// swagger:operation GET /api/v1/organizations/{org_id} getOrganization
	//
	// Retrieves an organization by the `org_id`.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get an organization
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	// responses:
	//   '200':
	//     description: Successfully got the organization
	//     schema:
	//       $ref: '#/definitions/GetOrganizationResponse'
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
	getOrgEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}", string(types.URLParamOrgID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
			},
		},
	)

	getOrgHandler := orgs.NewOrgGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getOrgEndpoint,
		Handler:  getOrgHandler,
		Router:   r,
	})

	return routes
}
