package loader

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/api/serverutils/erroralerter"
	"github.com/hatchet-dev/hatchet/internal/adapter"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/hatchet-dev/hatchet/internal/logger"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm"
	"github.com/joeshaw/envdecode"
)

type EnvDecoderConf struct {
	ServerConfigFile   server.ConfigFile
	DatabaseConfigFile database.ConfigFile
	SharedConfigFile   shared.ConfigFile
}

// ServerConfigFromEnv loads the server config file from environment variables
func ServerConfigFromEnv() (*server.ConfigFile, error) {
	var envDecoderConf EnvDecoderConf = EnvDecoderConf{}

	if err := envdecode.StrictDecode(&envDecoderConf); err != nil {
		return nil, fmt.Errorf("Failed to decode server conf: %s", err)
	}

	return &envDecoderConf.ServerConfigFile, nil
}

// DatabaseConfigFromEnv loads the database config file from environment variables
func DatabaseConfigFromEnv() (*database.ConfigFile, error) {
	var envDecoderConf EnvDecoderConf = EnvDecoderConf{}

	if err := envdecode.StrictDecode(&envDecoderConf); err != nil {
		return nil, fmt.Errorf("Failed to decode database conf: %s", err)
	}

	return &envDecoderConf.DatabaseConfigFile, nil
}

// SharedConfigFromEnv loads the shared config file from environment variables
func SharedConfigFromEnv() (*shared.ConfigFile, error) {
	var envDecoderConf EnvDecoderConf = EnvDecoderConf{}

	if err := envdecode.StrictDecode(&envDecoderConf); err != nil {
		return nil, fmt.Errorf("Failed to decode server conf: %s", err)
	}

	return &envDecoderConf.SharedConfigFile, nil
}

type EnvConfigLoader struct {
	version string
}

func (e *EnvConfigLoader) loadSharedConfig() (res *shared.Config, err error) {
	sharedConfig, err := SharedConfigFromEnv()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config from env: %v", err)
	}

	l := logger.NewConsole(sharedConfig.Debug)

	errorAlerter := erroralerter.NoOpAlerter{}

	return &shared.Config{
		Logger:       *l,
		ErrorAlerter: errorAlerter,
	}, nil
}

func (e *EnvConfigLoader) LoadDatabaseConfig() (res *database.Config, err error) {
	dc, err := DatabaseConfigFromEnv()

	if err != nil {
		return nil, fmt.Errorf("could not load database config from env: %v", err)
	}

	db, err := adapter.New(dc)

	if err != nil {
		return nil, fmt.Errorf("could not load database from adapter: %v", err)
	}

	var key [32]byte

	for i, b := range []byte(dc.EncryptionKey) {
		key[i] = b
	}

	repo := gorm.NewRepository(db, &key)

	return &database.Config{
		GormDB:     db,
		Repository: repo,
	}, nil
}

func (e *EnvConfigLoader) LoadServerConfig() (res *server.Config, err error) {
	sharedConfig, err := e.loadSharedConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load shared config: %v", err)
	}

	dbConfig, err := e.LoadDatabaseConfig()

	if err != nil {
		return nil, fmt.Errorf("could not load database config: %v", err)
	}

	sc, err := ServerConfigFromEnv()

	if err != nil {
		return nil, fmt.Errorf("could not load server config from env: %v", err)
	}

	authConfig := server.AuthConfig{
		BasicAuthEnabled: sc.BasicAuthEnabled,
	}

	serverRuntimeConfig := server.ServerRuntimeConfig{
		Port: sc.Port,
	}

	return &server.Config{
		DB:                  *dbConfig,
		Config:              *sharedConfig,
		AuthConfig:          authConfig,
		ServerRuntimeConfig: serverRuntimeConfig,
	}, nil
}
