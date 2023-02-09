package provisionerutils

import (
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
)

func GetProvisionerOpts(team *models.Team, mod *models.Module, run *models.ModuleRun, config *server.Config) (*provisioner.ProvisionOpts, error) {
	envVars := make(map[string]string)
	var err error

	envVarVersion := &mod.CurrentModuleEnvVarsVersion

	if envVarVersion == nil || envVarVersion.ID == "" {
		envVarVersion, err = config.DB.Repository.ModuleEnvVars().ReadModuleEnvVarsVersionByID(mod.ID, mod.CurrentModuleEnvVarsVersionID)

		if err != nil {
			return nil, err
		}
	}

	if envVarVersion != nil {
		envVars, err = envVarVersion.GetEnvVars(config.DB.GetEncryptionKey())

		if err != nil {
			return nil, err
		}
	}

	return &provisioner.ProvisionOpts{
		Team:       team,
		Module:     mod,
		ModuleRun:  run,
		EnvVars:    envVars,
		TokenOpts:  *config.TokenOpts,
		Repository: config.DB.Repository,
		ServerURL:  config.ServerRuntimeConfig.ServerURL,
	}, nil
}
