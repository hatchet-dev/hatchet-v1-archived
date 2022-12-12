package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

type Route struct {
	Endpoint *endpoint.APIEndpoint
	Handler  http.Handler
	Router   chi.Router
}

type Registerer struct {
	GetRoutes func(
		r chi.Router,
		config *server.Config,
		basePath *endpoint.Path,
		factory endpoint.APIEndpointFactory,
		children ...*Registerer,
	) []*Route

	Children []*Registerer
}
