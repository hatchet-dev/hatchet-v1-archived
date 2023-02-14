package modules

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type ModuleCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)

	request := &types.CreateModuleRequest{}

	if ok := m.DecodeAndValidate(w, r, request); !ok {
		return
	}

	mod := &models.Module{
		TeamID:              team.ID,
		Name:                request.Name,
		DeploymentMechanism: models.DeploymentMechanismGithub,
	}

	if request.DeploymentGithub == nil {
		m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: fmt.Sprintf("at least one deployment mechanism must be specified"),
		}, http.StatusBadRequest))

		return
	}

	if github := request.DeploymentGithub; github != nil {
		deplConfig, reqErr := setupGithubDeploymentConfig(m.Config(), request.DeploymentGithub, team, user)

		if reqErr != nil {
			m.HandleAPIError(w, r, reqErr)
			return
		}

		mod.DeploymentConfig = *deplConfig
	}

	mod, err := m.Repo().Module().CreateModule(mod)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	var mvv *models.ModuleValuesVersion

	if valuesGithub := request.ValuesGithub; valuesGithub != nil {
		// ensure that the app installation id exists and the user has access to it
		_, reqErr := canAccessGithubAppInstallation(m.Config(), valuesGithub.GithubAppInstallationID, user)

		if reqErr != nil {
			m.HandleAPIError(w, r, reqErr)
			return
		}

		mvv, err = createModuleValuesGithub(m.Config(), mod, valuesGithub, 0)
	} else {
		mvv, err = createModuleValuesRaw(m.Config(), mod, request.ValuesRaw, 0)
	}

	// set values version, this is updated later to reduce DB queries
	mod.CurrentModuleValuesVersionID = mvv.ID

	// create env vars
	mev, err := models.NewModuleEnvVarsVersion(mod.ID, 0, request.EnvVars)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	mev, err = m.Repo().ModuleEnvVars().CreateModuleEnvVarsVersion(mev)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	mod.CurrentModuleEnvVarsVersionID = mev.ID

	mod, err = m.Repo().Module().UpdateModule(mod)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	m.WriteResult(w, r, mod.ToAPIType())
}
