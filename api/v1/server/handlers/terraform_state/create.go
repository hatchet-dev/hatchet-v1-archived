package terraform_state

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type TerraformStateCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTerraformStateCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TerraformStateCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TerraformStateCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	if run.LockID != "" {
		// make sure the lock request is correct
		req := &types.CreateTerraformStateRequest{}

		if ok := t.DecodeAndValidateQueryOnly(w, r, req); !ok {
			return
		}

		if req.ID != run.LockID {
			w.WriteHeader(http.StatusLocked)
			return
		}
	}

	// TODO(abelanger5): pull path logic into separate method
	path := fmt.Sprintf("%s/%s/%s", team.ID, module.ID, run.ID)

	// read state file
	fileBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	err = t.Config().DefaultFileStore.WriteFile(path, fileBytes, true)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	return
}
