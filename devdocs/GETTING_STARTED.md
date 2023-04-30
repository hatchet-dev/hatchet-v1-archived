## Getting Started

This guide shows you how to get a development environment up and running with Hatchet.

Prerequisites:

- Go `1.19+` installed (can be checked with `go version`)
- Node version at least `v18.15.0` and npm version at least `9.5.0` (can be checked with `npm version`)
- `air` installed (can be installed via `go install github.com/cosmtrek/air@latest`)

First, set your environment variables. At a minimum, run the following to set your `.env` file:

```sh
echo "
DEBUG=true

DATABASE_KIND=sqlite
DATABASE_SQLITE_PATH=./hatchet.db
DATABASE_ENCRYPTION_KEY=$(cat /dev/urandom | base64 | head -c 32)

SERVER_RUNTIME_URL=http://localhost:8081
SERVER_RUNTIME_RUN_BACKGROUND_WORKER=true
SERVER_RUNTIME_RUN_RUNNER_WORKER=true
SERVER_RUNTIME_RUN_TEMPORAL_WORKER=true
SERVER_RUNTIME_RUN_STATIC_FILE_SERVER=true

SERVER_AUTH_COOKIE_SECRETS=\"$(cat /dev/urandom | base64 | head -c 16) $(cat /dev/urandom | base64 | head -c 16)\"
SERVER_AUTH_COOKIE_DOMAIN=localhost:8081
SERVER_AUTH_COOKIE_INSECURE=true
" > .env
```

Next, run the following to set `./dashboard/.env`:

```sh
echo "
DEV_SERVER_PORT=8081
DEV_SERVER_SOCKET_PORT=8081
ENABLE_PROXY=true
API_SERVER=http://localhost:8080
" > ./dashboard/.env
```

Then, run `cd ./dashboard && npm i --legacy-peer-deps && cd ..`.

Finally, run `make start-dev` - if you're building for the first time, it may take up to 2 minutes for the Go server to build. When both the webpack server and Go server have started, navigate to `http://localhost:8081` -- you should see the Hatchet login screen!

## Making API Changes

The API files can be found in `./api/v1`. To add or update an endpoint:

1. Update the corresponding type in `api/v1/types`. Unless the change is very minor, it's recommended that you add any new types to a draft PR as soon as possible (ideally before you start working on the PR) so that the changes can be discussed.
2. Update or add the handler in `./api/v1/server/handlers`. See [creating a new handler](#creating-a-new-handler).
3. Update or add the route in `./api/v1/server/router`. See [creating a route](#creating-a-route).
4. Generate new client files. See [building the api](#building-the-api).

### Creating a New Handler

To create a new handler, the file should contain a handler constructor that returns an `http.Handler`. However, you will likely want this handler to inherit `handlers.HatchetHandlerReadWriter` and initialize your handler using `handlers.NewDefaultHatchetHandler`, as the default handler contains a number of helper methods for decoding, validating, and handling API errors. As an example:

```go
type ResourceHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewResourceHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ResourceHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (r *ResourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // core handler code
}
```

Hatchet handlers should **only execute core business logic** -- they should not perform user authn/authz (this is done by the middleware packages, see `./api/v1/server/authn` or `./api/v1/server/authz`). They should not execute overly complex backend logic - you should create a package in `internal` for this. Overall, handlers should be extremely readable and easily modifiable.

You can access parent resources through the request context - all parent resources are added to the context using the `authz` middleware, as we verify permissions while populating resources. For example, if registering a route under `/modules/{module_id}/my-resource`, you can access the parent `Module` via:

```go
module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
```

### Creating a Route

Once the handler is written, you can create a new route by registering the handler in `./api/v1/server/router`. To actually register the route, append routes in the corresponding file or create a new `router.Registerer` in a new file. Additionally, make sure that you add the corresponding `go-swagger` comment above the route registration (see existing files as a reference). For example, a new `GET` endpoint for the `ResourceHandler` we wrote above would look like:

```go
// GET /api/v1/users/current -> users.UserGetCurrentHandler
// swagger:operation GET /api/v1/users/current getCurrentUser
//
// ### Description
//
// Retrieves the API resource.
//
// ---
// produces:
// - application/json
// summary: Retrieve the API Resource
// tags:
// - APIResource
// responses:
//   '200':
//     description: Successfully got the API resource
//     schema:
//       $ref: '#/definitions/GetUserResponse'
//   '403':
//     description: Forbidden
//     schema:
//       $ref: '#/definitions/APIErrorForbiddenExample'
getAPIResource := factory.NewAPIEndpoint(
    &endpoint.EndpointMetadata{
        Verb:   types.APIVerbGet,
        Method: types.HTTPVerbGet,
        Path: &endpoint.Path{
            Parent:       basePath,
            RelativePath: "/api-resource",
		},
        // Scopes sets the set of authn/authz permissions required for this endpoint. Most routes should
        // have `UserScope`.
		Scopes: []types.PermissionScope{
			types.UserScope,
		},
	},
)

getAPIResourceHandler := NewResourceHandler(
    config,
    factory.GetDecoderValidator(),
    factory.GetResultWriter(),
)

routes = append(routes, &router.Route{
    Endpoint: getAPIResource,
    Handler:  getAPIResourceHandler,
    Router:   r,
})
```

### Building the API

If you make changes to the API, you will need to rebuild the API clients. This can be done via:

```sh
make build-api-client # builds typescript client
make build-api-client-golang # builds golang client
```
