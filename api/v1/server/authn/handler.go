package authn

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// AuthNFactory generates a middleware handler `AuthN`
type AuthNFactory struct {
	config *server.Config
}

// NewAuthNFactory returns an `AuthNFactory` that uses the passed-in server
// config
func NewAuthNFactory(
	config *server.Config,
) *AuthNFactory {
	return &AuthNFactory{config}
}

// NewAuthenticated creates a new instance of `AuthN` that implements the http.Handler
// interface.
func (f *AuthNFactory) NewAuthenticated(next http.Handler) http.Handler {
	return &AuthN{next, f.config, false}
}

// NewAuthenticatedWithRedirect creates a new instance of `AuthN` that implements the http.Handler
// interface. This handler redirects the user to login if the user is not attached, and stores a
// redirect URI in the session, if the session exists.
func (f *AuthNFactory) NewAuthenticatedWithRedirect(next http.Handler) http.Handler {
	return &AuthN{next, f.config, true}
}

// AuthN implements the authentication middleware
type AuthN struct {
	next     http.Handler
	config   *server.Config
	redirect bool
}

// ServeHTTP attaches an authenticated subject to the request context,
// or serves a forbidden error. If authenticated, it calls the next handler.
func (authn *AuthN) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store := authn.config.UserSessionStore

	session, err := store.Get(r, store.GetName())

	if err != nil {
		session.Values["authenticated"] = false

		// we attempt to save the session, but do not catch the error since we send the
		// forbidden error regardless
		session.Save(r, w)

		authn.sendForbiddenError(err, w, r)
		return
	}

	if auth, ok := session.Values["authenticated"].(bool); !auth || !ok {
		authn.handleForbiddenForSession(w, r, fmt.Errorf("stored cookie was not authenticated"), session)
		return
	}

	// read the user id in the token
	userID, ok := session.Values["user_id"].(string)

	if !ok {
		authn.handleForbiddenForSession(w, r, fmt.Errorf("could not cast user_id to string"), session)
		return
	}

	authn.nextWithUserID(w, r, userID)
}

func (authn *AuthN) handleForbiddenForSession(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	session *sessions.Session,
) {
	if authn.redirect {
		// need state parameter to validate when redirected
		if r.URL.RawQuery == "" {
			session.Values["redirect_uri"] = r.URL.Path
		} else {
			session.Values["redirect_uri"] = r.URL.Path + "?" + r.URL.RawQuery
		}

		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", 302)
	} else {
		authn.sendForbiddenError(err, w, r)
	}

	return
}

// nextWithUserID calls the next handler with the user set in the context with key
// `types.UserScope`.
func (authn *AuthN) nextWithUserID(w http.ResponseWriter, r *http.Request, userID string) {
	// search for the user
	user, err := authn.config.DB.Repository.User().ReadUserByID(userID)

	if err != nil {
		authn.sendForbiddenError(fmt.Errorf("user with id %s not found in database", userID), w, r)
		return
	}

	// add the user to the context
	ctx := r.Context()
	ctx = context.WithValue(ctx, types.UserScope, user)

	r = r.Clone(ctx)
	authn.next.ServeHTTP(w, r)
}

// sendForbiddenError sends a 403 Forbidden error to the end user while logging a
// specific error
func (authn *AuthN) sendForbiddenError(err error, w http.ResponseWriter, r *http.Request) {
	reqErr := apierrors.NewErrForbidden(err)

	apierrors.HandleAPIError(authn.config.Logger, authn.config.ErrorAlerter, w, r, reqErr, true)
}
