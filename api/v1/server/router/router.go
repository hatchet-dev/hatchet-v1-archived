package router

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

func NewAPIRouter(config *server.Config) *chi.Mux {
	r := chi.NewRouter()

	endpointFactory := endpoint.NewAPIObjectEndpointFactory(config)

	baseRegisterer := NewBaseRegisterer()
	userRegisterer := NewUserScopedRegisterer()

	baseRoutePath := "/api/v1"

	r.Route(baseRoutePath, func(r chi.Router) {
		baseRoutePath := &endpoint.Path{
			RelativePath: baseRoutePath,
		}

		baseRoutes := baseRegisterer.GetRoutes(
			r,
			config,
			baseRoutePath,
			endpointFactory,
		)

		userRoutes := userRegisterer.GetRoutes(
			r,
			config,
			&endpoint.Path{
				Parent:       baseRoutePath,
				RelativePath: "",
			},
			endpointFactory,
			userRegisterer.Children...,
		)

		routes := [][]*router.Route{
			baseRoutes,
			userRoutes,
		}

		var allRoutes []*router.Route

		for _, r := range routes {
			allRoutes = append(allRoutes, r...)
		}

		registerRoutes(config, allRoutes)
	})

	return r
}

func registerRoutes(config *server.Config, routes []*router.Route) {
	for _, route := range routes {

		fmt.Println("parsing route:", route, route.Endpoint.Metadata.Path.GetPathString())

		atomicGroup := route.Router.Group(nil)

		atomicGroup.Method(
			string(route.Endpoint.Metadata.Method),
			route.Endpoint.Metadata.Path.RelativePath,
			route.Handler,
		)
	}
}
