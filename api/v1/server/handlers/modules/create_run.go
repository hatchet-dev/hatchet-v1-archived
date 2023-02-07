package modules

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
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
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)

	run := &models.ModuleRun{
		ModuleID: module.ID,
		Status:   models.ModuleRunStatusInProgress,
		Kind:     models.ModuleRunKindPlan,
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

	// TODO(abelanger5): queue, don't run plan
	err = m.Config().DefaultProvisioner.RunPlan(&provisioner.ProvisionOpts{
		Team:       team,
		Module:     module,
		ModuleRun:  run,
		TokenOpts:  *m.Config().TokenOpts,
		Repository: m.Repo(),
		ServerURL:  m.Config().ServerRuntimeConfig.ServerURL,
	})

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
