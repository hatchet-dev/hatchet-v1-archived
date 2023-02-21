package provisionerutils

import (
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
)

func GetProvisionerEnvOpts(team *models.Team, mod *models.Module, run *models.ModuleRun, dbConfig database.Config, tokenOpts token.TokenOpts, serverURL string) (*provisioner.GetEnvOpts, error) {
	envVars := make(map[string]string)
	var err error

	envVarVersion := &mod.CurrentModuleEnvVarsVersion

	if envVarVersion == nil || envVarVersion.ID == "" {
		envVarVersion, err = dbConfig.Repository.ModuleEnvVars().ReadModuleEnvVarsVersionByID(mod.ID, mod.CurrentModuleEnvVarsVersionID)

		if err != nil {
			return nil, err
		}
	}

	if envVarVersion != nil {
		envVars, err = envVarVersion.GetEnvVars(dbConfig.GetEncryptionKey())

		if err != nil {
			return nil, err
		}
	}

	return &provisioner.GetEnvOpts{
		Team:       team,
		Module:     mod,
		ModuleRun:  run,
		EnvVars:    envVars,
		TokenOpts:  tokenOpts,
		Repository: dbConfig.Repository,
		ServerURL:  serverURL,
	}, nil
}
