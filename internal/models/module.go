package models

import "github.com/hatchet-dev/hatchet/api/v1/types"

type DeploymentMechanism string

const (
	DeploymentMechanismGithub DeploymentMechanism = "github"
	DeploymentMechanismAPI    DeploymentMechanism = "api"
)

type Module struct {
	Base

	TeamID string
	Team   Team `gorm:"foreignKey:TeamID"`

	Name string

	DeploymentMechanism DeploymentMechanism

	DeploymentConfig ModuleDeploymentConfig

	// TODO(abelanger5): mechanism for values/secrets

	Runs []ModuleRun
}

func (m *Module) ToAPIType() *types.Module {
	return &types.Module{
		APIResourceMeta:  m.ToAPITypeMetadata(),
		Name:             m.Name,
		DeploymentConfig: *m.DeploymentConfig.ToAPIType(),
	}
}

type ModuleDeploymentConfig struct {
	Base

	ModuleID string

	ModulePath string

	GithubRepoName   string
	GithubRepoOwner  string
	GithubRepoBranch string

	GithubAppInstallationID string
	GithubAppInstallation   GithubAppInstallation `gorm:"foreignKey:GithubAppInstallationID"`
}

func (m *ModuleDeploymentConfig) ToAPIType() *types.ModuleDeploymentConfig {
	return &types.ModuleDeploymentConfig{
		Path:                    m.ModulePath,
		GithubRepoName:          m.GithubRepoName,
		GithubRepoOwner:         m.GithubRepoOwner,
		GithubRepoBranch:        m.GithubRepoBranch,
		GithubAppInstallationID: m.GithubAppInstallationID,
	}
}

type ModuleRun struct {
	Base

	ModuleID string
}
