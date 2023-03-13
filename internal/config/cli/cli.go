package cli

import (
	"github.com/hatchet-dev/hatchet/api/v1/client/fileclient"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	Address        string `mapstructure:"address"`
	APIToken       string `mapstructure:"apiToken"`
	OrganizationID string `mapstructure:"organizationID"`
	TeamID         string `mapstructure:"teamID"`
}

type Config struct {
	shared.Config

	ConfigFile *ConfigFile

	FileClient *fileclient.FileClient

	APIClient *swagger.APIClient
}

func BindAllEnv(v *viper.Viper) {
	v.BindEnv("address", "HATCHET_ADDRESS")
	v.BindEnv("apiToken", "HATCHET_API_TOKEN")
	v.BindEnv("organizationID", "HATCHET_ORGANIZATION_ID")
	v.BindEnv("teamID", "HATCHET_TEAM_ID")
}
