package terraform_state

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type TerraformStateUnlockHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTerraformStateUnlockHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TerraformStateUnlockHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TerraformStateUnlockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	run.LockID = ""

	run, err := t.Repo().Module().UpdateModuleRun(run)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

}
