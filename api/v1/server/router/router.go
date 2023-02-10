package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	moduleRegisterer := NewModuleRouteRegisterer()

	baseRoutePath := "/api/v1"

	chi.RegisterMethod("LOCK")
	chi.RegisterMethod("UNLOCK")

	r.Route(baseRoutePath, func(r chi.Router) {
		r.Use(middleware.Logger)

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

		r.Route(fmt.Sprintf("/teams/{%s}", types.URLParamTeamID), func(r chi.Router) {
			routes = append(
				routes,
				moduleRegisterer.GetRoutes(
					r,
					config,
					&endpoint.Path{
						Parent:       baseRoutePath,
						RelativePath: fmt.Sprintf("/teams/{%s}", types.URLParamTeamID),
					},
					endpointFactory,
				),
			)
		})

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
	authNBasicFactory := authn.NewAuthNBasicFactory(config)
	noAuthNFactory := authn.NewNoAuthNFactory(config)

	orgFactory := authz.NewOrgScopedFactory(config)
	orgMemberFactory := authz.NewOrgMemberScopedFactory(config)

	teamFactory := authz.NewTeamScopedFactory(config)
	teamMemberFactory := authz.NewTeamMemberScopedFactory(config)

	moduleFactory := authz.NewModuleScopedFactory(config)
	moduleRunFactory := authz.NewRunScopedFactory(config)
	moduleEnvVarsFactory := authz.NewModuleEnvVarsScopedFactory(config)
	moduleValuesFactory := authz.NewModuleValuesScopedFactory(config)

	gaiFactory := authz.NewGithubAppInstallationScopedFactory(config)

	contentTypeFactory := ContentTypeMiddleware

	for _, route := range routes {
		atomicGroup := route.Router.Group(nil)

		atomicGroup.Use(contentTypeFactory(route.Endpoint.Metadata.ContentType))

		// always register user scopes first
		for _, scope := range route.Endpoint.Metadata.Scopes {
			switch scope {
			case types.UserScope:
				if !config.AuthConfig.RequireEmailVerification || route.Endpoint.Metadata.AllowUnverifiedEmails {
					atomicGroup.Use(authNFactory.NewAuthenticatedWithoutEmailVerification)
				} else {
					atomicGroup.Use(authNFactory.NewAuthenticated)
				}
			case types.BasicAuthUserScope:
				if !config.AuthConfig.RequireEmailVerification || route.Endpoint.Metadata.AllowUnverifiedEmails {
					atomicGroup.Use(authNBasicFactory.NewAuthenticatedWithoutEmailVerification)
				} else {
					atomicGroup.Use(authNBasicFactory.NewAuthenticated)
				}
			case types.NoUserScope:
				atomicGroup.Use(noAuthNFactory.NewNotAuthenticated)
			case types.GithubAppInstallationScope:
				endpointMetaFactory := endpoint.NewEndpointMiddleware(config, route.Endpoint.Metadata)
				atomicGroup.Use(endpointMetaFactory.Middleware)

				atomicGroup.Use(gaiFactory.Middleware)
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
			case types.ModuleScope:
				atomicGroup.Use(moduleFactory.Middleware)
			case types.ModuleRunScope:
				atomicGroup.Use(moduleRunFactory.Middleware)
			case types.ModuleEnvVarScope:
				atomicGroup.Use(moduleEnvVarsFactory.Middleware)
			case types.ModuleValuesScope:
				atomicGroup.Use(moduleValuesFactory.Middleware)
			}
		}

		atomicGroup.Method(
			string(route.Endpoint.Metadata.Method),
			route.Endpoint.Metadata.Path.RelativePath,
			route.Handler,
		)
	}
}

func ContentTypeMiddleware(contentType string) func(http.Handler) http.Handler {
	if contentType == "" {
		contentType = "application/json;charset=utf8"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", contentType)
			next.ServeHTTP(w, r)
		})
	}
}
