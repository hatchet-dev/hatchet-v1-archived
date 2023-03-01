package runner

import (
	"github.com/hatchet-dev/hatchet/api/v1/client/fileclient"
	"github.com/hatchet-dev/hatchet/api/v1/client/grpc"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
)

type ConfigFile struct {

	// GRPC Connection Options

	// GRPC server address
	GRPCServerAddress string `env:"GRPC_SERVER_ADDRESS,default=http://localhost:8080"`

	// GRPC connection auth
	GRPCToken string `env:"GRPC_TOKEN"`

	// The IDs that this run operates on
	TeamID          string `env:"TEAM_ID"`
	ModuleID        string `env:"MODULE_ID"`
	ModuleRunID     string `env:"MODULE_RUN_ID"`
	ModuleMonitorID string `env:"MODULE_MONITOR_ID"`

	// API client options
	APIToken string `env:"API_TOKEN"`

	APIServerAddress string `env:"API_SERVER_ADDRESS,default=http://localhost:8080"`

	// Github options
	GithubRepositoryName string `env:"GITHUB_REPOSITORY_NAME"`
	GithubModulePath     string `env:"GITHUB_MODULE_PATH"`
	GithubSHA            string `env:"GITHUB_SHA"`
	GithubRepositoryDest string `env:"GITHUB_REPOSITORY_DEST,default=./bin/tmp"`

	// Terraform options
	// TFDir is a relative or absolute path to the terraform directory
	TFDir string `env:"TF_DIR,default=../../terraform"`
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
