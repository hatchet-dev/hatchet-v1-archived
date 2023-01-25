package server

import (
	"strings"

	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/auth/cookie"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
	"github.com/hatchet-dev/hatchet/internal/integrations/oauth/github"
	"github.com/hatchet-dev/hatchet/internal/notifier"
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

	// RestrictedEmailDomains sets the restricted email domains for the instance.
	RestrictedEmailDomains []string `env:"RESTRICTED_EMAIL_DOMAINS"`

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

	// Notification options

	// Sendgrid notifier options
	SendgridAPIKey                string `env:"SENDGRID_API_KEY"`
	SendgridPWResetTemplateID     string `env:"SENDGRID_PW_RESET_TEMPLATE_ID"`
	SendgridVerifyEmailTemplateID string `env:"SENDGRID_VERIFY_EMAIL_TEMPLATE_ID"`
	SendgridInviteLinkTemplateID  string `env:"SENDGRID_INVITE_LINK_TEMPLATE_ID"`
	SendgridSenderEmail           string `env:"SENDGRID_SENDER_EMAIL"`

	// Github App options
	GithubAppClientID      string `env:"GITHUB_APP_CLIENT_ID"`
	GithubAppClientSecret  string `env:"GITHUB_APP_CLIENT_SECRET"`
	GithubAppName          string `env:"GITHUB_APP_NAME"`
	GithubAppWebhookSecret string `env:"GITHUB_APP_WEBHOOK_SECRET"`
	GithubAppID            string `env:"GITHUB_APP_ID"`
	GithubAppSecretPath    string `env:"GITHUB_APP_SECRET_PATH"`

	// S3 file storage options
	S3StateAWSAccessKeyID string `env:"S3_STATE_AWS_ACCESS_KEY_ID"`
	S3StateAWSSecretKey   string `env:"S3_STATE_AWS_SECRET_KEY"`
	S3StateAWSRegion      string `env:"S3_STATE_AWS_REGION"`
	S3StateBucketName     string `env:"S3_STATE_BUCKET_NAME"`
	S3StateEncryptionKey  string `env:"S3_STATE_ENCRYPTION_KEY,default=__random_strong_encryption_key__"`
}

type AuthConfig struct {
	RequireEmailVerification bool
	BasicAuthEnabled         bool
	RestrictedEmailDomains   []string
}

func (a *AuthConfig) IsEmailAllowed(email string) bool {
	if len(a.RestrictedEmailDomains) == 0 {
		return true
	}

	targetComponents := strings.Split(email, "@")
	targetDomain := targetComponents[1]

	for _, domain := range a.RestrictedEmailDomains {
		if domain == targetDomain {
			return true
		}
	}

	return false
}

type ServerRuntimeConfig struct {
	ServerURL  string
	Port       int
	CookieName string
}

type Config struct {
	shared.Config

	DB database.Config

	AuthConfig AuthConfig

	ServerRuntimeConfig ServerRuntimeConfig

	UserSessionStore *cookie.UserSessionStore

	TokenOpts *token.TokenOpts

	UserNotifier notifier.UserNotifier

	GithubApp *github.GithubAppConf

	DefaultFileStore filestorage.FileStorageManager
}

func (c *Config) ToAPIServerMetadataType() *types.APIServerMetadata {
	return &types.APIServerMetadata{
		Auth: &types.APIServerMetadataAuth{
			RequireEmailVerification: c.AuthConfig.RequireEmailVerification,
		},
	}
}
