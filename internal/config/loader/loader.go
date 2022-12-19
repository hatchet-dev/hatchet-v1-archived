package loader

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/api/serverutils/erroralerter"
	"github.com/hatchet-dev/hatchet/internal/adapter"
	"github.com/hatchet-dev/hatchet/internal/auth/cookie"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
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

	return e.LoadSharedConfigFromConfigFile(sharedConfig)
}

func (e *EnvConfigLoader) LoadSharedConfigFromConfigFile(sharedConfigFile *shared.ConfigFile) (res *shared.Config, err error) {
	l := logger.NewConsole(sharedConfigFile.Debug)

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

	res = &database.Config{
		GormDB:     db,
		Repository: repo,
	}

	res.SetEncryptionKey(&key)

	return res, nil
}

func (e *EnvConfigLoader) LoadServerConfigFromEnv() (res *server.Config, err error) {
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

	return e.LoadServerConfigFromConfigFile(sc, dbConfig, sharedConfig)
}

func (e *EnvConfigLoader) LoadServerConfigFromConfigFile(sc *server.ConfigFile, dbConfig *database.Config, sharedConfig *shared.Config) (res *server.Config, err error) {
	authConfig := server.AuthConfig{
		BasicAuthEnabled: sc.BasicAuthEnabled,
	}

	serverRuntimeConfig := server.ServerRuntimeConfig{
		Port:      sc.Port,
		ServerURL: sc.ServerURL,
	}

	userSessionStore, err := cookie.NewUserSessionStore(&cookie.UserSessionStoreOpts{
		SessionRepository:   dbConfig.Repository.UserSession(),
		CookieSecrets:       sc.CookieSecrets,
		CookieAllowInsecure: sc.CookieAllowInsecure,
		CookieDomain:        sc.CookieDomain,
		CookieName:          sc.CookieName,
	})

	if err != nil {
		return nil, fmt.Errorf("could not initialize session store: %v", err)
	}

	tokenOpts := &token.TokenOpts{
		Issuer:   sc.TokenIssuerURL,
		Audience: sc.TokenAudience,
	}

	if sc.TokenIssuerURL == "" {
		tokenOpts.Issuer = sc.ServerURL
	}

	if len(sc.TokenAudience) == 0 {
		tokenOpts.Audience = []string{sc.ServerURL}
	}

	return &server.Config{
		DB:                  *dbConfig,
		Config:              *sharedConfig,
		AuthConfig:          authConfig,
		ServerRuntimeConfig: serverRuntimeConfig,
		UserSessionStore:    userSessionStore,
		TokenOpts:           tokenOpts,
	}, nil
}
