package authn

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/encryption"
	"github.com/hatchet-dev/hatchet/internal/models"
)

func SaveUserAuthenticated(
	w http.ResponseWriter,
	r *http.Request,
	config *server.Config,
	user *models.User,
) (string, error) {
	session, err := config.UserSessionStore.Get(r, config.UserSessionStore.GetName())

	if err != nil {
		return "", err
	}

	var redirect string

	if valR := session.Values["redirect_uri"]; valR != nil {
		redirect = session.Values["redirect_uri"].(string)
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID
	session.Values["email"] = user.Email

	// we unset the redirect uri after login
	session.Values["redirect_uri"] = ""

	return redirect, session.Save(r, w)
}

func SaveOAuthState(
	w http.ResponseWriter,
	r *http.Request,
	config *server.Config,
) (string, error) {
	stateBytes, err := encryption.GenerateRandomBytes(16)

	if err != nil {
		return "", err
	}

	session, err := config.UserSessionStore.Get(r, config.ServerRuntimeConfig.CookieName)

	if err != nil {
		return "", err
	}

	// need state parameter to validate when redirected
	session.Values["state"] = string(stateBytes)

	// need a parameter to indicate that this was triggered through the oauth flow
	session.Values["oauth_triggered"] = true

	if err := session.Save(r, w); err != nil {
		return "", err
	}

	return string(stateBytes), nil
}

func ValidateOAuthState(
	w http.ResponseWriter,
	r *http.Request,
	config *server.Config,
) (isValidated bool, isOAuthTriggered bool, err error) {
	session, err := config.UserSessionStore.Get(r, config.ServerRuntimeConfig.CookieName)

	if err != nil {
		return false, false, err
	}

	if _, ok := session.Values["state"]; !ok {
		return false, false, fmt.Errorf("state parameter not found in session")
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		return false, false, fmt.Errorf("state parameters do not match")
	}

	if isOAuthTriggeredVal, exists := session.Values["oauth_triggered"]; exists {
		isOAuthTriggered, ok := isOAuthTriggeredVal.(bool)

		isOAuthTriggered = ok && isOAuthTriggered
	}

	// need state parameter to validate when redirected
	session.Values["state"] = ""
	session.Values["oauth_triggered"] = false

	if err := session.Save(r, w); err != nil {
		return false, false, fmt.Errorf("could not clear session")
	}

	return true, isOAuthTriggered, nil
}

func SaveUserUnauthenticated(
	w http.ResponseWriter,
	r *http.Request,
	config *server.Config,
) error {
	session, err := config.UserSessionStore.Get(r, config.UserSessionStore.GetName())

	if err != nil {
		return err
	}

	session.Values["authenticated"] = false
	session.Values["user_id"] = nil
	session.Values["email"] = nil

	// we set the maxage of the session so that the session gets deleted. This avoids cases
	// where the same cookie can get re-authed to a different user, which would be problematic
	// if the session values weren't properly cleared on logout.
	session.Options.MaxAge = -1

	return session.Save(r, w)
}
