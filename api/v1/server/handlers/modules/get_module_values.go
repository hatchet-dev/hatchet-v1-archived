package modules

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type ModuleValuesGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleValuesGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleValuesGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleValuesGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mev, _ := r.Context().Value(types.ModuleValuesScope).(*models.ModuleValuesVersion)

	var mv *models.ModuleValues
	var err error

	if mev.Kind == models.ModuleValuesVersionKindDatabase {
		mv, err = m.Repo().ModuleValues().ReadModuleValuesByVersionID(mev.ID)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

	}

	res, err := mev.ToAPIType(mv)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	m.WriteResult(w, r, res)
}
