package handlers

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type HatchetHandler interface {
	Config() *server.Config
	Repo() repository.Repository
	HandleAPIError(w http.ResponseWriter, r *http.Request, err apierrors.RequestError)
	HandleAPIErrorNoWrite(w http.ResponseWriter, r *http.Request, err apierrors.RequestError)
	// PopulateOAuthSession(
	// 	w http.ResponseWriter,
	// 	r *http.Request,
	// 	state string,
	// 	isProject, isUser bool,
	// 	integrationClient types.OAuthIntegrationClient,
	// 	integrationID uint,
	// ) error
}

type HatchetHandlerWriter interface {
	HatchetHandler
	handlerutils.ResultWriter
}

type HatchetHandlerReader interface {
	HatchetHandler
	handlerutils.RequestDecoderValidator
}

type HatchetHandlerReadWriter interface {
	HatchetHandlerWriter
	HatchetHandlerReader
}

type DefaultHatchetHandler struct {
	config           *server.Config
	decoderValidator handlerutils.RequestDecoderValidator
	writer           handlerutils.ResultWriter
}

func NewDefaultHatchetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) HatchetHandlerReadWriter {
	return &DefaultHatchetHandler{config, decoderValidator, writer}
}

func (d *DefaultHatchetHandler) Config() *server.Config {
	return d.config
}

func (d *DefaultHatchetHandler) Repo() repository.Repository {
	return d.config.DB.Repository
}

func (d *DefaultHatchetHandler) HandleAPIError(w http.ResponseWriter, r *http.Request, err apierrors.RequestError) {
	apierrors.HandleAPIError(d.Config().Logger, d.Config().ErrorAlerter, w, r, err, true)
}

func (d *DefaultHatchetHandler) HandleAPIErrorNoWrite(w http.ResponseWriter, r *http.Request, err apierrors.RequestError) {
	apierrors.HandleAPIError(d.Config().Logger, d.Config().ErrorAlerter, w, r, err, false)
}

func (d *DefaultHatchetHandler) WriteResult(w http.ResponseWriter, r *http.Request, v interface{}) {
	d.writer.WriteResult(w, r, v)
}

func (d *DefaultHatchetHandler) DecodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	return d.decoderValidator.DecodeAndValidate(w, r, v)
}

func (d *DefaultHatchetHandler) DecodeAndValidateNoWrite(r *http.Request, v interface{}) error {
	return d.decoderValidator.DecodeAndValidateNoWrite(r, v)
}

func IgnoreAPIError(w http.ResponseWriter, r *http.Request, err apierrors.RequestError) {
	return
}

// func (d *DefaultHatchetHandler) PopulateOAuthSession(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	state string,
// 	isProject, isUser bool,
// 	integrationClient types.OAuthIntegrationClient,
// 	integrationID uint,
// ) error {
// 	session, err := d.Config().Store.Get(r, d.Config().ServerConf.CookieName)

// 	if err != nil {
// 		return err
// 	}

// 	// need state parameter to validate when redirected
// 	session.Values["state"] = state

// 	// check if redirect uri is populated, then overwrite
// 	if redirect := r.URL.Query().Get("redirect_uri"); redirect != "" {
// 		session.Values["redirect_uri"] = redirect
// 	}

// 	if isProject {
// 		project, _ := r.Context().Value(types.ProjectScope).(*models.Project)

// 		if project == nil {
// 			return fmt.Errorf("could not read project")
// 		}

// 		session.Values["project_id"] = project.ID
// 	}

// 	if isUser {
// 		user, _ := r.Context().Value(types.UserScope).(*models.User)

// 		if user == nil {
// 			return fmt.Errorf("could not read user")
// 		}

// 		session.Values["user_id"] = user.ID
// 	}

// 	if integrationID != 0 && len(integrationClient) > 0 {
// 		session.Values["integration_id"] = integrationID
// 		session.Values["integration_client"] = string(integrationClient)
// 	}

// 	if err := session.Save(r, w); err != nil {
// 		return err
// 	}

// 	return nil
// }

type Unavailable struct {
	config    *server.Config
	handlerID string
}

func NewUnavailable(config *server.Config, handlerID string) *Unavailable {
	return &Unavailable{config, handlerID}
}

func (u *Unavailable) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	apierrors.HandleAPIError(u.config.Logger, u.config.ErrorAlerter, w, r, apierrors.NewErrPassThroughToClient(
		types.APIError{
			Description: fmt.Sprintf("%s not available in community edition", u.handlerID),
			Code:        types.ErrCodeUnavailable,
		},
		http.StatusMethodNotAllowed,
	), true)
}
