package modules

import (
	"errors"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type ModuleUpdateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleUpdateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleUpdateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)

	request := &types.UpdateModuleRequest{}

	if ok := m.DecodeAndValidate(w, r, request); !ok {
		return
	}

	if request.Name != "" && request.Name != module.Name {
		module.Name = request.Name
	}

	if github := request.DeploymentGithub; github != nil {
		deplConfig, reqErr := setupGithubDeploymentConfig(m.Config(), request.DeploymentGithub, team, user)

		if reqErr != nil {
			m.HandleAPIError(w, r, reqErr)
			return
		}

		module.DeploymentConfig = *deplConfig
	}

	valuesGithub := request.ValuesGithub
	valuesRaw := request.ValuesRaw
	var prevMVV *models.ModuleValuesVersion
	var err error

	if valuesGithub != nil && valuesRaw != nil {
		m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
			types.APIError{
				Description: "Cannot update a module with both raw and github-based values",
				Code:        types.ErrCodeBadRequest,
			}, http.StatusBadRequest,
		))

		return
	}

	if valuesGithub != nil || valuesRaw != nil {
		prevMVV, err = m.Repo().ModuleValues().ReadModuleValuesVersionByID(module.ID, module.CurrentModuleValuesVersionID)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}
	}

	if valuesGithub != nil {
		mvv, err := createModuleValuesGithub(m.Config(), module, valuesGithub, prevMVV.Version)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		module.CurrentModuleValuesVersionID = mvv.ID
	} else if valuesRaw != nil {
		mvv, err := createModuleValuesRaw(m.Config(), module, request.ValuesRaw, prevMVV.Version)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		module.CurrentModuleValuesVersionID = mvv.ID
	}

	if request.EnvVars != nil {
		prevMEV, err := m.Repo().ModuleEnvVars().ReadModuleEnvVarsVersionByID(module.ID, module.CurrentModuleEnvVarsVersionID)

		var prevMEVVersion uint

		if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		} else if prevMEV != nil {
			prevMEVVersion = prevMEV.Version
		} else {
			prevMEVVersion = 0
		}

		mev, err := models.NewModuleEnvVarsVersion(module.ID, prevMEVVersion, request.EnvVars)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		mev, err = m.Repo().ModuleEnvVars().CreateModuleEnvVarsVersion(mev)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		module.CurrentModuleEnvVarsVersionID = mev.ID
	}

	module, err = m.Repo().Module().UpdateModule(module)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	m.WriteResult(w, r, module.ToAPIType())
}
