package logs

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type LogGetHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewLogGetHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &LogGetHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *LogGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	res := &types.GetLogsResponse{}

	// case on the log storage location
	if run.LogLocation == models.LogLocationFileStorage {
		fileBytes, err := m.Config().DefaultFileStore.ReadFile(filestorage.GetModuleRunLogsPath(module.TeamID, module.ID, run.ID), true)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		logs := make([]string, 0)

		err = json.Unmarshal(fileBytes, &logs)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		res.Logs = logs
	} else {
		logsReader := logstorage.NewLogStrArr()

		err := m.Config().DefaultLogStore.ReadLogs(context.Background(), &logstorage.LogGetOpts{
			Path:  logstorage.GetLogStoragePath(module.TeamID, module.ID, run.ID),
			Count: 0,
		}, logsReader)

		if err != nil {
			m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		res.Logs = logsReader.GetLogs()
	}

	m.WriteResult(w, r, res)
}
