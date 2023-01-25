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

type TerraformStateLockHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTerraformStateLockHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TerraformStateLockHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TerraformStateLockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	if run.LockID != "" {
		w.WriteHeader(http.StatusLocked)

		res := run.ToTerraformLockType()
		t.WriteResult(w, r, res)

		return
	}

	req := &types.LockTerraformStateRequest{}

	if ok := t.DecodeAndValidate(w, r, req); !ok {
		return
	}

	run.LockID = req.ID
	run.LockInfo = req.Info
	run.LockCreated = req.Created
	run.LockOperation = req.Operation
	run.LockPath = req.Path
	run.LockWho = req.Who

	run, err := t.Repo().Module().UpdateModuleRun(run)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}
}
