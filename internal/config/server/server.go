package server

import (
	"github.com/hatchet-dev/hatchet/internal/auth/cookie"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
)

type ConfigFile struct {
	// General server config options

	// Port is the port that the core server listens on
	Port int `env:"SERVER_PORT,default=8080"`

	// ServerURL is the full server URL of the instance, INCLUDING protocol.
	// We include the protocol as several auth implementations depend on it, like
	// JWT token and cookies.
	ServerURL string `env:"SERVER_URL,default=https://hatchet.run"`

	// Authn and authz options

	// BasedAuthEnabled controls whether email and password-based login is enabled for this
	// Hatchet instances
	BasicAuthEnabled bool `env:"BASIC_AUTH_ENABLED,default=true"`

	CookieName          string   `env:"COOKIE_NAME,default=hatchet"`
	CookieDomain        string   `env:"COOKIE_DOMAIN"`
	CookieSecrets       []string `env:"COOKIE_SECRETS,default=random_hash_key_;random_block_key"`
	CookieAllowInsecure bool     `env:"COOKIE_INSECURE,default=false"`

	// TokenIssuerURL is the endpoint of the issuer, typically equivalent to the server URL.
	// This field should INCLUDE the protocol.
	// If this is not set, it is set to the SERVER_URL variable.
	TokenIssuerURL string `env:"TOKEN_ISSUER_URL"`

	// TokenAudience is the set of audiences for the JWT token issuer, typically equivalent to the server URL.
	// This field should INCLUDE the protocol.
	// If this is not set, it is set to the SERVER_URL variable.
	TokenAudience []string `env:"TOKEN_AUDIENCE"`
}

type AuthConfig struct {
	BasicAuthEnabled bool
}

type ServerRuntimeConfig struct {
	ServerURL string
	Port      int
}

type Config struct {
	shared.Config

	DB database.Config

	AuthConfig AuthConfig

	ServerRuntimeConfig ServerRuntimeConfig

	UserSessionStore *cookie.UserSessionStore

	TokenOpts *token.TokenOpts
}
