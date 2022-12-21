package authn

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/internal/config/server"
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
