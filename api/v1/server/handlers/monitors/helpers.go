package monitors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

func getMonitorModulesFromRequest(config *server.Config, team *models.Team, modules []string) ([]models.Module, apierrors.RequestError) {
	var targetModuleIDs []string

	if modules == nil || len(modules) == 0 {
		targetModuleIDs = make([]string, 0)
	} else {
		targetModuleIDs = modules
	}

	var modErr error
	var targetModules []models.Module = make([]models.Module, 0)

	for _, modID := range targetModuleIDs {
		// ensure all modules are in the team
		mod, err := config.DB.Repository.Module().ReadModuleByID(team.ID, modID)

		if err != nil {
			if errors.Is(err, repository.RepositoryErrorNotFound) {
				return nil, apierrors.NewErrPassThroughToClient(types.APIError{
					Code:        types.ErrCodeBadRequest,
					Description: fmt.Sprintf("Could not find module with id %s", modID),
				}, http.StatusBadRequest)
			} else {
				modErr = multierror.Append(modErr, err)
			}
		}

		targetModules = append(targetModules, *mod)
	}

	if modErr != nil {
		return nil, apierrors.NewErrInternal(modErr)
	}

	return targetModules, nil
}
