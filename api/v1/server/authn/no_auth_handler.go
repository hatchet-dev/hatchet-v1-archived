package authn

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// NoAuthNFactory generates a middleware handler `NoAuthN` which verifies that
// there are NO valid auth credentials attached to this request. This is useful for things
// like creating a new user or logging in.
type NoAuthNFactory struct {
	config *server.Config
}

// NewNoAuthNFactory returns an `AuthNFactory` that uses the passed-in server
// config
func NewNoAuthNFactory(
	config *server.Config,
) *NoAuthNFactory {
	return &NoAuthNFactory{config}
}

// NewNotAuthenticated creates a new instance of `NoAuthN` that implements the http.Handler
// interface.
func (f *NoAuthNFactory) NewNotAuthenticated(next http.Handler) http.Handler {
	return &NoAuthN{next, f.config, false}
}

// NoAuthN implements the authentication middleware
type NoAuthN struct {
	next     http.Handler
	config   *server.Config
	redirect bool
}

// ServeHTTP checks for credentials. If they exist, it throws bad request to the user.
func (noauthn *NoAuthN) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// first check for a bearer token
	tok, err := getPATFromRequest(r, noauthn.config)

	if err == nil && tok != nil {
		noauthn.sendBadRequest(w, r)
		return
	}

	// check for a cookie-based user session
	store := noauthn.config.UserSessionStore

	session, err := store.Get(r, store.GetName())

	if err != nil {
		noauthn.next.ServeHTTP(w, r)
		return
	}

	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		noauthn.sendBadRequest(w, r)
		return
	}

	noauthn.next.ServeHTTP(w, r)
}

// sendBadRequest sends a 400 bad request error to the end user if they have auth credentials
func (noauthn *NoAuthN) sendBadRequest(w http.ResponseWriter, r *http.Request) {
	reqErr := apierrors.NewErrPassThroughToClient(types.APIError{
		Description: "Valid credentials are already attached to this request. Clear your credentials and try again.",
		Code:        types.ErrCodeBadRequest,
	}, http.StatusBadRequest)

	apierrors.HandleAPIError(noauthn.config.Logger, noauthn.config.ErrorAlerter, w, r, reqErr, true)
}
