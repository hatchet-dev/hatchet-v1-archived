package server

import (
	"github.com/hatchet-dev/hatchet/internal/auth/cookie"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
)

type ConfigFile struct {
	// General server config options

	// Port is the port that the core server listens on
	Port int `env:"SERVER_PORT,default=8080"`

	// Authn and authz options

	// BasedAuthEnabled controls whether email and password-based login is enabled for this
	// Hatchet instances
	BasicAuthEnabled bool `env:"BASIC_AUTH_ENABLED,default=true"`

	CookieName          string   `env:"COOKIE_NAME,default=hatchet"`
	CookieDomain        string   `env:"COOKIE_DOMAIN"`
	CookieSecrets       []string `env:"COOKIE_SECRETS,default=random_hash_key_;random_block_key"`
	CookieAllowInsecure bool     `env:"COOKIE_INSECURE,default=false"`
}

type AuthConfig struct {
	BasicAuthEnabled bool
}

type ServerRuntimeConfig struct {
	Port int
}

type Config struct {
	shared.Config

	DB database.Config

	AuthConfig AuthConfig

	ServerRuntimeConfig ServerRuntimeConfig

	UserSessionStore *cookie.UserSessionStore
}
