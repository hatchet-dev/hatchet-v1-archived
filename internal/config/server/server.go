package server

import (
	"strings"

	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/auth/cookie"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
	"github.com/hatchet-dev/hatchet/internal/notifier"
	"github.com/hatchet-dev/hatchet/internal/queuemanager"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	Runtime ConfigFileRuntime `mapstructure:"runtime" json:"runtime,omitempty"`

	Auth shared.ConfigFileAuth `mapstructure:"auth" json:"auth,omitempty"`

	Notification ConfigFileNotification `mapstructure:"notification" json:"notification,omitempty"`

	VCS ConfigFileVCS `mapstructure:"vcs" json:"vcs,omitempty"`

	FileStore shared.FileStorageConfigFile `mapstructure:"fileStore" json:"fileStore,omitempty"`

	LogStore shared.LogStoreConfigFile `mapstructure:"logStore" json:"logStore,omitempty"`
}

// General server runtime options
type ConfigFileRuntime struct {
	// Port is the port that the core server listens on
	Port int `mapstructure:"port" json:"port,omitempty" default:"8080"`

	// ServerURL is the full server URL of the instance, INCLUDING protocol.
	// We include the protocol as several auth implementations depend on it, like
	// JWT token and cookies.
	ServerURL string `mapstructure:"url" validator:"url" json:"url,omitempty" default:"http://localhost:8080"`

	// BroadcastGRPCAddress is the endpoint for the grpc server to be used by clients
	BroadcastGRPCAddress string `mapstructure:"broadcastGRPCAddress" json:"broadcastGRPCAddress,omitempty" validator:"url" default:"http://localhost:8080"`

	RunBackgroundWorker  bool   `mapstructure:"runBackgroundWorker" json:"runBackgroundWorker,omitempty" default:"false"`
	RunRunnerWorker      bool   `mapstructure:"runRunnerWorker" json:"runRunnerWorker,omitempty" default:"false"`
	RunTemporalServer    bool   `mapstructure:"runTemporalServer" json:"runTemporalServer,omitempty" default:"false"`
	RunStaticFileServer  bool   `mapstructure:"runStaticFileServer" json:"runStaticFileServer,omitempty" default:"false"`
	StaticFileServerPath string `mapstructure:"staticFilePath" json:"staticFilePath,omitempty"`

	PermittedModuleDeploymentMechanisms []string `mapstructure:"permittedModuleDeploymentMechanisms" json:"permittedModuleDeploymentMechanisms,omitempty" default:"[\"github\",\"api\",\"local\"]"`
}

type ConfigFileNotification struct {
	Sendgrid ConfigFileNotificationSendgrid `mapstructure:"sendgrid" json:"sendgrid,omitempty"`
}

type ConfigFileNotificationSendgrid struct {
	SendgridAPIKey                string `mapstructure:"apiKey" json:"apiKey,omitempty"`
	SendgridPWResetTemplateID     string `mapstructure:"pwResetTemplateID" json:"pwResetTemplateID,omitempty"`
	SendgridVerifyEmailTemplateID string `mapstructure:"verifyEmailTemplateID" json:"verifyEmailTemplateID,omitempty"`
	SendgridInviteLinkTemplateID  string `mapstructure:"inviteLinkTemplateID" json:"inviteLinkTemplateID,omitempty"`
	SendgridSenderEmail           string `mapstructure:"senderEmail" json:"senderEmail,omitempty" validator:"email"`
}

type ConfigFileVCS struct {
	Github ConfigFileGithub `mapstructure:"github" json:"github,omitempty"`
}

type ConfigFileGithub struct {
	Enabled                bool   `mapstructure:"enabled" json:"enabled"`
	GithubAppClientID      string `mapstructure:"appClientID" json:"appClientID,omitempty"`
	GithubAppClientSecret  string `mapstructure:"appClientSecret" json:"appClientSecret,omitempty"`
	GithubAppName          string `mapstructure:"appName" json:"appName,omitempty"`
	GithubAppWebhookSecret string `mapstructure:"appWebhookSecret" json:"appWebhookSecret,omitempty"`
	GithubAppID            string `mapstructure:"appID" json:"appID,omitempty"`
	GithubAppSecretPath    string `mapstructure:"appSecretPath" json:"appSecretPath,omitempty"`
}

func BindAllEnv(v *viper.Viper) {
	v.BindEnv("runtime.url", "SERVER_RUNTIME_URL")
	v.BindEnv("runtime.port", "SERVER_RUNTIME_PORT")
	v.BindEnv("runtime.broadcastGRPCAddress", "SERVER_RUNTIME_BROADCAST_GRPC_ADDRESS")
	v.BindEnv("runtime.runBackgroundWorker", "SERVER_RUNTIME_RUN_BACKGROUND_WORKER")
	v.BindEnv("runtime.runRunnerWorker", "SERVER_RUNTIME_RUN_RUNNER_WORKER")
	v.BindEnv("runtime.runTemporalServer", "SERVER_RUNTIME_RUN_TEMPORAL_WORKER")
	v.BindEnv("runtime.runStaticFileServer", "SERVER_RUNTIME_RUN_STATIC_FILE_SERVER")
	v.BindEnv("runtime.staticFilePath", "SERVER_RUNTIME_STATIC_FILE_PATH")

	v.BindEnv("auth.restrictedEmailDomains", "SERVER_AUTH_RESTRICTED_EMAIL_DOMAINS")
	v.BindEnv("auth.basicAuthEnabled", "SERVER_AUTH_BASIC_AUTH_ENABLED")
	v.BindEnv("auth.cookie.name", "SERVER_AUTH_COOKIE_NAME")
	v.BindEnv("auth.cookie.domain", "SERVER_AUTH_COOKIE_DOMAIN")
	v.BindEnv("auth.cookie.secrets", "SERVER_AUTH_COOKIE_SECRETS")
	v.BindEnv("auth.cookie.insecure", "SERVER_AUTH_COOKIE_INSECURE")
	v.BindEnv("auth.token.issuer", "SERVER_AUTH_TOKEN_ISSUER")
	v.BindEnv("auth.token.audience", "SERVER_AUTH_TOKEN_AUDIENCE")

	v.BindEnv("notification.sendgrid.apiKey", "SERVER_NOTIFICATION_SENDGRID_API_KEY")
	v.BindEnv("notification.sendgrid.pwResetTemplateID", "SERVER_NOTIFICATION_SENDGRID_PW_RESET_TEMPLATE_ID")
	v.BindEnv("notification.sendgrid.verifyEmailTemplateID", "SERVER_NOTIFICATION_SENDGRID_VERIFY_EMAIL_TEMPLATE_ID")
	v.BindEnv("notification.sendgrid.inviteLinkTemplateID", "SERVER_NOTIFICATION_SENDGRID_INVITE_LINK_TEMPLATE_ID")
	v.BindEnv("notification.sendgrid.senderEmail", "SERVER_NOTIFICATION_SENDGRID_SENDER_EMAIL")

	v.BindEnv("vcs.kind", "SERVER_VCS_KIND")
	v.BindEnv("vcs.github.enabled", "SERVER_VCS_GITHUB_ENABLED")
	v.BindEnv("vcs.github.appClientID", "SERVER_VCS_GITHUB_APP_CLIENT_ID")
	v.BindEnv("vcs.github.appClientSecret", "SERVER_VCS_GITHUB_APP_CLIENT_SECRET")
	v.BindEnv("vcs.github.appName", "SERVER_VCS_GITHUB_APP_NAME")
	v.BindEnv("vcs.github.appWebhookSecret", "SERVER_VCS_GITHUB_APP_WEBHOOK_SECRET")
	v.BindEnv("vcs.github.appID", "SERVER_VCS_GITHUB_APP_ID")
	v.BindEnv("vcs.github.appSecretPath", "SERVER_VCS_GITHUB_APP_SECRET_PATH")

	v.BindEnv("fileStore.kind", "SERVER_FILESTORE_KIND")
	v.BindEnv("fileStore.s3.accessKeyID", "SERVER_FILESTORE_S3_ACCESS_KEY_ID")
	v.BindEnv("fileStore.s3.secretKey", "SERVER_FILESTORE_S3_SECRET_KEY")
	v.BindEnv("fileStore.s3.region", "SERVER_FILESTORE_S3_REGION")
	v.BindEnv("fileStore.s3.bucketName", "SERVER_FILESTORE_S3_BUCKET_NAME")
	v.BindEnv("fileStore.s3.encryptionKey", "SERVER_FILESTORE_S3_ENCRYPTION_KEY")
	v.BindEnv("fileStore.local.directory", "SERVER_FILESTORE_LOCAL_DIRECTORY")
	v.BindEnv("fileStore.local.encryptionKey", "SERVER_FILESTORE_LOCAL_ENCRYPTION_KEY")

	v.BindEnv("logStore.kind", "SERVER_LOGSTORE_KIND")
	v.BindEnv("logStore.redis.host", "SERVER_LOGSTORE_REDIS_HOST")
	v.BindEnv("logStore.redis.port", "SERVER_LOGSTORE_REDIS_PORT")
	v.BindEnv("logStore.redis.user", "SERVER_LOGSTORE_REDIS_USER")
	v.BindEnv("logStore.redis.password", "SERVER_LOGSTORE_REDIS_PASSWORD")
	v.BindEnv("logStore.redis.db", "SERVER_LOGSTORE_REDIS_DB")
	v.BindEnv("logStore.file.directory", "SERVER_LOGSTORE_FILE_DIRECTORY")
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
	Version string

	ServerURL            string
	Port                 int
	CookieName           string
	BroadcastGRPCAddress string
	RunBackgroundWorker  bool
	RunRunnerWorker      bool
	RunTemporalServer    bool
	RunStaticFileServer  bool
	StaticFileServerPath string

	PermittedModuleDeploymentMechanisms []string
}

type Config struct {
	shared.Config

	DB database.Config

	AuthConfig AuthConfig

	ServerRuntimeConfig ServerRuntimeConfig

	UserSessionStore *cookie.UserSessionStore

	TokenOpts *token.TokenOpts

	UserNotifier notifier.UserNotifier

	VCSProviders map[vcs.VCSRepositoryKind]vcs.VCSProvider

	DefaultFileStore filestorage.FileStorageManager

	DefaultLogStore logstorage.LogStorageBackend

	ModuleRunQueueManager queuemanager.ModuleRunQueueManager
}

func (c *Config) ToAPIServerMetadataType() *types.APIServerMetadata {
	return &types.APIServerMetadata{
		Auth: &types.APIServerMetadataAuth{
			RequireEmailVerification: c.AuthConfig.RequireEmailVerification,
		},
		Integrations: &types.APIServerMetadataIntegrations{
			GithubApp: c.VCSProviders[vcs.VCSRepositoryKindGithub] != nil,
			Email:     c.UserNotifier.GetID() != "noop",
		},
	}
}
