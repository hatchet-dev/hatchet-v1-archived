package authn

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// AuthNBasicFactory generates a middleware handler `AuthNBasic`
type AuthNBasicFactory struct {
	config *server.Config
}

// NewAuthNBasicFactory returns an `AuthNBasicFactory` that uses the passed-in server
// config
func NewAuthNBasicFactory(
	config *server.Config,
) *AuthNBasicFactory {
	return &AuthNBasicFactory{config}
}

// NewAuthenticated creates a new instance of `AuthNBasic` that implements the http.Handler
// interface.
func (f *AuthNBasicFactory) NewAuthenticated(next http.Handler) http.Handler {
	return &AuthNBasic{next, f.config, true, false}
}

func (f *AuthNBasicFactory) NewAuthenticatedWithoutEmailVerification(next http.Handler) http.Handler {
	return &AuthNBasic{next, f.config, false, false}
}

// NewAuthenticatedWithRedirect creates a new instance of `AuthNBasic` that implements the http.Handler
// interface. This handler redirects the user to login if the user is not attached, and stores a
// redirect URI in the session, if the session exists.
func (f *AuthNBasicFactory) NewAuthenticatedWithRedirect(requireEmailVerification bool, next http.Handler) http.Handler {
	return &AuthNBasic{next, f.config, requireEmailVerification, true}
}

// AuthNBasic implements the authentication middleware
type AuthNBasic struct {
	next                     http.Handler
	config                   *server.Config
	requireEmailVerification bool
	redirect                 bool
}

// ServeHTTP attaches an authenticated subject to the request context,
// or serves a forbidden error. If authenticated, it calls the next handler.
func (authn *AuthNBasic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tokenKind, hatchetToken, ok := r.BasicAuth()

	if ok {
		if tokenKind != "pat" && tokenKind != "mrt" {
			authn.sendForbiddenError(fmt.Errorf("basic auth username must be either pat or mrt"), w, r)
			return
		}

		if hatchetToken == "" {
			authn.sendForbiddenError(fmt.Errorf("hatchet token does not exist"), w, r)
			return
		}

		if tokenKind == "pat" {
			pat, err := token.GetPATFromEncoded(hatchetToken, authn.config.DB.Repository.PersonalAccessToken(), authn.config.TokenOpts)

			if err != nil {
				authn.sendForbiddenError(err, w, r)
				return
			}

			if pat.Revoked || pat.IsExpired() {
				authn.sendForbiddenError(fmt.Errorf("token with id %s not valid", pat.ID), w, r)
				return
			}

			authn.nextWithUserID(w, r, pat.UserID)
			return
		} else if tokenKind == "mrt" {
			mrt, err := token.GetMRTFromEncoded(hatchetToken, authn.config.DB.Repository.Module(), authn.config.TokenOpts)

			if err != nil {
				authn.sendForbiddenError(err, w, r)
				return
			}

			if mrt.Revoked || mrt.IsExpired() {
				authn.sendForbiddenError(fmt.Errorf("token with id %s not valid", mrt.ID), w, r)
				return
			}

			authn.nextWithUserID(w, r, mrt.UserID)
			return
		}

		return
	}

	authn.sendForbiddenError(fmt.Errorf("no basic auth credentials"), w, r)
}

// sendEmailNotVerifiedError sends a 400 Bad Request error to the end user indicating that the email
// has not been verified
func (authn *AuthNBasic) sendEmailNotVerifiedError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	reqErr := apierrors.NewErrPassThroughToClient(types.APIError{
		Description: "Email is not verified. Please verify your email and try again.",
		Code:        types.ErrCodeEmailNotVerified,
	}, http.StatusUnprocessableEntity)

	apierrors.HandleAPIError(authn.config.Logger, authn.config.ErrorAlerter, w, r, reqErr, true)
}

// nextWithUserID calls the next handler with the user set in the context with key
// `types.UserScope`.
func (authn *AuthNBasic) nextWithUserID(w http.ResponseWriter, r *http.Request, userID string) {
	// search for the user
	user, err := authn.config.DB.Repository.User().ReadUserByID(userID)

	if err != nil {
		authn.sendForbiddenError(fmt.Errorf("user with id %s not found in database", userID), w, r)
		return
	}

	if !authn.config.AuthConfig.IsEmailAllowed(user.Email) {
		authn.sendForbiddenError(fmt.Errorf("email is not in restricted domain list"), w, r)
		return
	}

	if authn.requireEmailVerification && !user.EmailVerified {
		authn.sendEmailNotVerifiedError(w, r)
		return
	}

	// add the user to the context
	ctx := r.Context()
	ctx = context.WithValue(ctx, types.UserScope, user)

	r = r.Clone(ctx)
	authn.next.ServeHTTP(w, r)
}

func (authn *AuthNBasic) sendForbiddenError(err error, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	reqErr := apierrors.NewErrForbidden(err)
	apierrors.HandleAPIError(authn.config.Logger, authn.config.ErrorAlerter, w, r, reqErr, true)
	return
}
