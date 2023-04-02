package runner

import (
	"github.com/hatchet-dev/hatchet/api/v1/client/fileclient"
	"github.com/hatchet-dev/hatchet/api/v1/client/grpc"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	Resources ConfigFileResources `mapstructure:"resources" json:"resources,omitempty"`

	GRPC ConfigFileGRPC `mapstructure:"grpc" json:"grpc,omitempty"`

	API ConfigFileAPI `mapstructure:"api" json:"api,omitempty"`

	VCS ConfigFileVCS `mapstructure:"vcs" json:"github,omitempty"`

	Terraform ConfigFileTerraform `mapstructure:"terraform" json:"terraform,omitempty"`
}

type ConfigFileResources struct {
	TeamID          string `mapstructure:"teamID" json:"teamID,omitempty"`
	ModuleID        string `mapstructure:"moduleID" json:"moduleID,omitempty"`
	ModuleRunID     string `mapstructure:"moduleRunID" json:"moduleRunID,omitempty"`
	ModuleMonitorID string `mapstructure:"moduleMonitorID" json:"moduleMonitorID,omitempty"`
}

type ConfigFileGRPC struct {
	// GRPC server address
	GRPCServerAddress string `mapstructure:"serverAddress" json:"serverAddress,omitempty" default:"http://localhost:8080"`

	// GRPC connection auth
	GRPCToken string `mapstructure:"token" json:"token,omitempty"`
}

type ConfigFileAPI struct {
	APIToken string `mapstructure:"token" json:"token,omitempty"`

	APIServerAddress string `mapstructure:"serverAddress" json:"serverAddress,omitempty" default:"http://localhost:8080"`
}

type ConfigFileVCS struct {
	VCSRepositoryName string `mapstructure:"repositoryName" json:"repositoryName,omitempty"`
	VCSModulePath     string `mapstructure:"modulePath" json:"modulePath,omitempty"`
	VCSSHA            string `mapstructure:"sha" json:"sha,omitempty"`
	VCSRepositoryDest string `mapstructure:"repositoryDest" json:"repositoryDest,omitempty" default:"./bin/tmp"`
}

type ConfigFileTerraform struct {
	// TFDir is a relative or absolute path to the terraform directory
	TFDir string `mapstructure:"dir" json:"dir,omitempty"`
}

type Config struct {
	shared.Config

	ConfigFile *ConfigFile

	GRPCClient *grpc.GRPCClient

	FileClient *fileclient.FileClient

	APIClient *swagger.APIClient

	GithubTarballURL string

	TerraformConf TerraformConf
}

// SetTerraformDir is used to set the terraform directory as it may change after
func (c *Config) SetTerraformDir(path string) {
	c.TerraformConf.TFDir = path
}

// TerraformConf is the configuration for Terraform params
type TerraformConf struct {
	TFDir string
}

func BindAllEnv(v *viper.Viper) {
	v.BindEnv("resources.teamID", "RUNNER_TEAM_ID")
	v.BindEnv("resources.moduleID", "RUNNER_MODULE_ID")
	v.BindEnv("resources.moduleRunID", "RUNNER_MODULE_RUN_ID")
	v.BindEnv("resources.moduleMonitorID", "RUNNER_MODULE_MONITOR_ID")

	v.BindEnv("grpc.serverAddress", "RUNNER_GRPC_SERVER_ADDRESS")
	v.BindEnv("grpc.token", "RUNNER_GRPC_TOKEN")

	v.BindEnv("api.serverAddress", "RUNNER_API_SERVER_ADDRESS")
	v.BindEnv("api.token", "RUNNER_API_TOKEN")

	v.BindEnv("vcs.repositoryName", "RUNNER_VCS_REPOSITORY_NAME")
	v.BindEnv("vcs.modulePath", "RUNNER_VCS_MODULE_PATH")
	v.BindEnv("vcs.sha", "RUNNER_VCS_SHA")
	v.BindEnv("vcs.repositoryDest", "RUNNER_VCS_REPOSITORY_DEST")

	v.BindEnv("terraform.dir", "RUNNER_TERRAFORM_DIR")
}
