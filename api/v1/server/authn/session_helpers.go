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
	return session.Save(r, w)
}
