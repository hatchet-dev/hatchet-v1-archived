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
	vcsvalues "github.com/hatchet-dev/hatchet/internal/integrations/valuesstorage/vcs"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"

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
	case models.ModuleValuesVersionKindVCS:
		vcsRepo, err := vcs.GetVCSRepositoryFromModule(m.Config().VCSProviders, module)

		if err != nil {
			m.HandleAPIErrorNoWrite(w, r, apierrors.NewErrInternal(err))

			return
		}

		valuesManager = vcsvalues.NewGithubValuesStore(vcsRepo, sha)
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
