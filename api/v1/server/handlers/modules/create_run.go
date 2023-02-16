package modules

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
)

type RunCreateHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewRunCreateHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &RunCreateHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *RunCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)

	run := &models.ModuleRun{
		ModuleID:    module.ID,
		Status:      models.ModuleRunStatusQueued,
		Kind:        models.ModuleRunKindPlan,
		LogLocation: m.Config().DefaultLogStore.GetID(),
	}

	desc, err := generateRunDescription(m.Config(), module, run, models.ModuleRunStatusInProgress)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	run.StatusDescription = desc

	run, err = m.Repo().Module().CreateModuleRun(run)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	err = m.Config().ModuleRunQueueManager.Enqueue(module, run)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	err = dispatcher.DispatchModuleRunQueueChecker(m.Config().TemporalClient.GetClient(), &modulequeuechecker.CheckQueueInput{
		TeamID:   module.TeamID,
		ModuleID: module.ID,
	})

	// err = m.Config().DefaultProvisioner.RunPlan(opts)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
