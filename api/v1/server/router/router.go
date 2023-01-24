package router

import (
	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/authn"
	"github.com/hatchet-dev/hatchet/api/v1/server/authz"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

func NewAPIRouter(config *server.Config) *chi.Mux {
	r := chi.NewRouter()

	endpointFactory := endpoint.NewAPIObjectEndpointFactory(config)

	baseRegisterer := NewBaseRegisterer()
	userRegisterer := NewUserRouteRegisterer()
	githubAppRegisterer := NewGithubAppRouteRegisterer()
	orgRegisterer := NewOrgRouteRegisterer()
	teamRegisterer := NewTeamRouteRegisterer()

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

		githubAppRoutes := githubAppRegisterer.GetRoutes(
			r,
			config,
			&endpoint.Path{
				Parent:       baseRoutePath,
				RelativePath: "",
			},
			endpointFactory,
			githubAppRegisterer.Children...,
		)

		orgRoutes := orgRegisterer.GetRoutes(
			r,
			config,
			&endpoint.Path{
				Parent:       baseRoutePath,
				RelativePath: "",
			},
			endpointFactory,
		)

		teamRoutes := teamRegisterer.GetRoutes(
			r,
			config,
			&endpoint.Path{
				Parent:       baseRoutePath,
				RelativePath: "",
			},
			endpointFactory,
		)

		routes := [][]*router.Route{
			baseRoutes,
			userRoutes,
			githubAppRoutes,
			orgRoutes,
			teamRoutes,
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
	// Create a new "user-scoped" factory which will create a new user-scoped request
	// after authentication. Each subsequent http.Handler can lookup the user in context.
	authNFactory := authn.NewAuthNFactory(config)
	noAuthNFactory := authn.NewNoAuthNFactory(config)

	orgFactory := authz.NewOrgScopedFactory(config)
	orgMemberFactory := authz.NewOrgMemberScopedFactory(config)

	teamFactory := authz.NewTeamScopedFactory(config)
	teamMemberFactory := authz.NewTeamMemberScopedFactory(config)

	for _, route := range routes {
		atomicGroup := route.Router.Group(nil)

		// always register user scopes first
		for _, scope := range route.Endpoint.Metadata.Scopes {
			switch scope {
			case types.UserScope:
				if !config.AuthConfig.RequireEmailVerification || route.Endpoint.Metadata.AllowUnverifiedEmails {
					atomicGroup.Use(authNFactory.NewAuthenticatedWithoutEmailVerification)
				} else {
					atomicGroup.Use(authNFactory.NewAuthenticated)
				}
			case types.NoUserScope:
				atomicGroup.Use(noAuthNFactory.NewNotAuthenticated)
			case types.OrgScope:
				endpointMetaFactory := endpoint.NewEndpointMiddleware(config, route.Endpoint.Metadata)
				atomicGroup.Use(endpointMetaFactory.Middleware)

				atomicGroup.Use(orgFactory.Middleware)
			case types.OrgMemberScope:
				atomicGroup.Use(orgMemberFactory.Middleware)
			case types.TeamScope:
				endpointMetaFactory := endpoint.NewEndpointMiddleware(config, route.Endpoint.Metadata)
				atomicGroup.Use(endpointMetaFactory.Middleware)

				atomicGroup.Use(teamFactory.Middleware)
			case types.TeamMemberScope:
				atomicGroup.Use(teamMemberFactory.Middleware)
			}
		}

		atomicGroup.Method(
			string(route.Endpoint.Metadata.Method),
			route.Endpoint.Metadata.Path.RelativePath,
			route.Handler,
		)
	}
}
