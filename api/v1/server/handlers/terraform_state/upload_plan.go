package terraform_state

import (
	"io/ioutil"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type TerraformPlanUploadZIPHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewTerraformPlanUploadZIPHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &TerraformPlanUploadZIPHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (t *TerraformPlanUploadZIPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	file, _, err := r.FormFile("file")

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)

	zipPlanPath := GetPlanZIPPath(team.ID, module.ID, run.ID)

	err = t.Config().DefaultFileStore.WriteFile(zipPlanPath, fileBytes, true)

	if err != nil {
		t.HandleAPIError(w, r, apierrors.NewErrInternal(err))

		return
	}

	t.WriteResult(w, r, nil)
}
