package router

import (
	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/metadata"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

func NewBaseRegisterer(children ...*router.Registerer) *router.Registerer {
	return &router.Registerer{
		GetRoutes: GetBaseRoutes,
		Children:  children,
	}
}

func GetBaseRoutes(
	r chi.Router,
	config *server.Config,
	basePath *endpoint.Path,
	factory endpoint.APIEndpointFactory,
	children ...*router.Registerer,
) []*router.Route {
	routes := make([]*router.Route, 0)

	// GET /api/v1/metadata -> metadata.NewServerMetadataGetHandler
	// swagger:operation GET /api/v1/metadata getServerMetadata
	//
	// ### Description
	//
	// Gets the metadata for the Hatchet instance.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get server metadata
	// tags:
	// - Metadata
	// responses:
	//   '200':
	//     description: Successfully got the metadata
	//     schema:
	//       $ref: '#/definitions/APIServerMetadata'
	getMetadataEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/metadata",
			},
			Scopes: []types.PermissionScope{},
		},
	)

	getMetadataHandler := metadata.NewServerMetadataGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getMetadataEndpoint,
		Handler:  getMetadataHandler,
		Router:   r,
	})

	return routes
}
