package modules

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
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
	}

	run, err := m.Repo().Module().CreateModuleRun(run)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	mrt, err := models.NewModuleRunTokenFromRunID(team.ServiceAccountRunnerID, run.ID)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	_, err = token.GenerateTokenFromMRT(mrt, m.Config().TokenOpts)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// write the user to the db
	mrt, err = m.Repo().Module().CreateModuleRunToken(mrt)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
