package modules

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/valuesstorage"
	"github.com/hatchet-dev/hatchet/internal/integrations/valuesstorage/db"
	githubv "github.com/hatchet-dev/hatchet/internal/integrations/valuesstorage/github"

	"github.com/hatchet-dev/hatchet/internal/integrations/git/github"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type ModuleValuesCurrentGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleValuesCurrentGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleValuesCurrentGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleValuesCurrentGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)

	req := &types.GetModuleValuesRequest{}

	if ok := m.DecodeAndValidate(w, r, req); !ok {
		return
	}

	sha := req.GithubSHA

	mvv, err := m.Repo().ModuleValues().ReadModuleValuesVersionByID(module.ID, module.CurrentModuleValuesVersionID)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	var valuesManager valuesstorage.ValuesStorageManager

	switch mvv.Kind {
	case models.ModuleValuesVersionKindDatabase:
		valuesManager = db.NewDatabaseValuesStore(m.Repo())
	case models.ModuleValuesVersionKindGithub:
		gai, err := m.Repo().GithubAppInstallation().ReadGithubAppInstallationByID(mvv.GithubAppInstallationID)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		githubClient, err := github.GetGithubAppClientFromGAI(m.Config(), gai)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		valuesManager = githubv.NewGithubValuesStore(githubClient, sha)
	default:
		m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: fmt.Sprintf("Module %s does not have an attached module values version object", module.ID),
		}, http.StatusBadRequest))

		return
	}

	vals, err := valuesManager.ReadValues(mvv)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	m.WriteResult(w, r, vals)
}
