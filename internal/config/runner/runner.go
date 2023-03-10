package runner

import (
	"github.com/hatchet-dev/hatchet/api/v1/client/fileclient"
	"github.com/hatchet-dev/hatchet/api/v1/client/grpc"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	Resources ConfigFileResources `mapstructure:"resources"`

	GRPC ConfigFileGRPC `mapstructure:"grpc"`

	API ConfigFileAPI `mapstructure:"api"`

	Github ConfigFileGithub `mapstructure:"github"`

	Terraform ConfigFileTerraform `mapstructure:"terraform"`
}

type ConfigFileResources struct {
	TeamID          string `mapstructure:"teamID"`
	ModuleID        string `mapstructure:"moduleID"`
	ModuleRunID     string `mapstructure:"moduleRunID"`
	ModuleMonitorID string `mapstructure:"moduleMonitorID"`
}

type ConfigFileGRPC struct {
	// GRPC server address
	GRPCServerAddress string `mapstructure:"serverAddress" default:"http://localhost:8080"`

	// GRPC connection auth
	GRPCToken string `mapstructure:"token"`
}

type ConfigFileAPI struct {
	APIToken string `mapstructure:"token"`

	APIServerAddress string `mapstructure:"serverAddress" default:"http://localhost:8080"`
}

type ConfigFileGithub struct {
	GithubRepositoryName string `mapstructure:"repositoryName"`
	GithubModulePath     string `mapstructure:"modulePath"`
	GithubSHA            string `mapstructure:"sha"`
	GithubRepositoryDest string `mapstructure:"repositoryDest" default:"./bin/tmp"`
}

type ConfigFileTerraform struct {
	// TFDir is a relative or absolute path to the terraform directory
	TFDir string `mapstructure:"dir"`
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

	v.BindEnv("github.repositoryName", "RUNNER_GITHUB_REPOSITORY_NAME")
	v.BindEnv("github.modulePath", "RUNNER_GITHUB_MODULE_PATH")
	v.BindEnv("github.sha", "RUNNER_GITHUB_SHA")
	v.BindEnv("github.repositoryDest", "RUNNER_GITHUB_REPOSITORY_DEST")

	v.BindEnv("terraform.dir", "RUNNER_TERRAFORM_DIR")
}
