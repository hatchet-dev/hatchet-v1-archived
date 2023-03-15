package modules

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
	"github.com/hatchet-dev/hatchet/internal/provisioner/provisionerutils"
)

type ModuleRunGetLocalTokenHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewModuleRunGetLocalTokenHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ModuleRunGetLocalTokenHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (m *ModuleRunGetLocalTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(types.UserScope).(*models.User)
	team, _ := r.Context().Value(types.TeamScope).(*models.Team)
	module, _ := r.Context().Value(types.ModuleScope).(*models.Module)
	run, _ := r.Context().Value(types.ModuleRunScope).(*models.ModuleRun)

	if module.DeploymentMechanism != models.DeploymentMechanismLocal {
		m.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(types.APIError{
			Code:        types.ErrCodeBadRequest,
			Description: fmt.Sprintf("cannot get deployment token for modules that aren't locally deployed"),
		}, http.StatusBadRequest))

		return
	}

	if module.DeploymentConfig.UserID != user.ID {
		m.HandleAPIError(w, r, apierrors.NewErrForbidden(
			fmt.Errorf("deployment config user id %s does not match user id %s", module.DeploymentConfig.UserID, user.ID),
		))

		return
	}

	envOpts, err := provisionerutils.GetProvisionerEnvOpts(
		team,
		module,
		run,
		m.Config().DB,
		*m.Config().TokenOpts,
		m.Config().ServerRuntimeConfig.ServerURL,
		m.Config().ServerRuntimeConfig.BroadcastGRPCAddress,
	)

	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	rawTok, err := provisioner.GetRunnerToken(envOpts)
	if err != nil {
		m.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	m.WriteResult(w, r, &types.GetModuleRunTokenResponse{
		Token: rawTok,
	})
}
