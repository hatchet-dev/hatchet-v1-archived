package action

import (
	"context"
	"fmt"

	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
)

func ErrorHandler(config *runner.Config, reportKind, description string) error {
	_, _, err := config.APIClient.ModulesApi.FinalizeModuleRun(
		context.Background(),
		swagger.FinalizeModuleRunRequest{
			Status:      string(types.ModuleRunStatusFailed),
			Description: description,
			ReportKind:  reportKind,
		},
		config.ConfigFile.Resources.TeamID,
		config.ConfigFile.Resources.ModuleID,
		config.ConfigFile.Resources.ModuleRunID,
	)

	if err != nil {
		return fmt.Errorf("Error reporting error: %s. Original error: %s", err.Error(), description)
	}

	return fmt.Errorf(description)
}

func SuccessHandler(config *runner.Config, reportKind, description string) error {
	_, _, err := config.APIClient.ModulesApi.FinalizeModuleRun(
		context.Background(),
		swagger.FinalizeModuleRunRequest{
			Status:      string(types.ModuleRunStatusCompleted),
			Description: description,
			ReportKind:  reportKind,
		},
		config.ConfigFile.Resources.TeamID,
		config.ConfigFile.Resources.ModuleID,
		config.ConfigFile.Resources.ModuleRunID,
	)

	if err != nil {
		return fmt.Errorf("Error reporting success: %s. Success message: %s", err.Error(), description)
	}

	return nil
}
