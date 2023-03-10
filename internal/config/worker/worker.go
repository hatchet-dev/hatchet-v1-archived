package worker

import (
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage"
	"github.com/hatchet-dev/hatchet/internal/notifier"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
	"github.com/hatchet-dev/hatchet/internal/queuemanager"
	"github.com/spf13/viper"
)

type BackgroundConfigFile struct {
	ServerURL string `mapstructure:"serverURL"`

	BroadcastGRPCAddress string `mapstructure:"broadcastGRPCAddress" validator:"url" default:"http://localhost:8080"`

	Auth shared.ConfigFileAuth `mapstructure:"auth"`

	FileStore shared.FileStorageConfigFile `mapstructure:"fileStore"`

	LogStore shared.LogStoreConfigFile `mapstructure:"logStore"`

	Notifier BackgroundConfigFileNotifier `mapstructure:"notifier"`
}

type BackgroundConfigFileNotifier struct {
	Kind string `mapstructure:"kind"`

	Sendgrid BackgroundConfigFileNotifierSendgrid `mapstructure:"sendgrid"`
}

type BackgroundConfigFileNotifierSendgrid struct {
	SendgridAPIKey             string `mapstructure:"apiKey"`
	SendgridIncidentTemplateID string `mapstructure:"incidentTemplateID"`
	SendgridSenderEmail        string `mapstructure:"senderEmail" validator:"email"`
}

type BackgroundConfig struct {
	shared.Config

	DB database.Config

	ServerURL string

	BroadcastGRPCAddress string

	TokenOpts *token.TokenOpts

	DefaultFileStore filestorage.FileStorageManager

	DefaultLogStore logstorage.LogStorageBackend

	ModuleRunQueueManager queuemanager.ModuleRunQueueManager

	IncidentNotifier notifier.IncidentNotifier
}

func BindAllBackgroundEnv(v *viper.Viper) {
	v.BindEnv("serverURL", "BACKGROUND_SERVER_URL")
	v.BindEnv("broadcastGRPCAddress", "BACKGROUND_BROADCAST_GRPC_ADDRESS")

	v.BindEnv("auth.restrictedEmailDomains", "BACKGROUND_AUTH_RESTRICTED_EMAIL_DOMAINS")
	v.BindEnv("auth.basicAuthEnabled", "BACKGROUND_AUTH_BASIC_AUTH_ENABLED")
	v.BindEnv("auth.cookie.name", "BACKGROUND_AUTH_COOKIE_NAME")
	v.BindEnv("auth.cookie.domain", "BACKGROUND_AUTH_COOKIE_DOMAIN")
	v.BindEnv("auth.cookie.secrets", "BACKGROUND_AUTH_COOKIE_SECRETS")
	v.BindEnv("auth.cookie.insecure", "BACKGROUND_AUTH_COOKIE_INSECURE")
	v.BindEnv("auth.token.issuer", "BACKGROUND_AUTH_TOKEN_ISSUER")
	v.BindEnv("auth.token.audience", "BACKGROUND_AUTH_TOKEN_AUDIENCE")

	v.BindEnv("fileStore.kind", "BACKGROUND_FILESTORE_KIND")
	v.BindEnv("fileStore.s3.accessKeyID", "BACKGROUND_FILESTORE_S3_ACCESS_KEY_ID")
	v.BindEnv("fileStore.s3.secretKey", "BACKGROUND_FILESTORE_S3_SECRET_KEY")
	v.BindEnv("fileStore.s3.region", "BACKGROUND_FILESTORE_S3_REGION")
	v.BindEnv("fileStore.s3.bucketName", "BACKGROUND_FILESTORE_S3_BUCKET_NAME")
	v.BindEnv("fileStore.s3.encryptionKey", "BACKGROUND_FILESTORE_S3_ENCRYPTION_KEY")
	v.BindEnv("fileStore.local.directory", "BACKGROUND_FILESTORE_LOCAL_DIRECTORY")
	v.BindEnv("fileStore.local.encryptionKey", "BACKGROUND_FILESTORE_LOCAL_ENCRYPTION_KEY")

	v.BindEnv("logStore.kind", "BACKGROUND_LOGSTORE_KIND")
	v.BindEnv("logStore.redis.host", "BACKGROUND_LOGSTORE_REDIS_HOST")
	v.BindEnv("logStore.redis.port", "BACKGROUND_LOGSTORE_REDIS_PORT")
	v.BindEnv("logStore.redis.user", "BACKGROUND_LOGSTORE_REDIS_USER")
	v.BindEnv("logStore.redis.password", "BACKGROUND_LOGSTORE_REDIS_PASSWORD")
	v.BindEnv("logStore.redis.db", "BACKGROUND_LOGSTORE_REDIS_DB")
	v.BindEnv("logStore.file.directory", "BACKGROUND_LOGSTORE_FILE_DIRECTORY")

	v.BindEnv("notifier.kind", "BACKGROUND_NOTIFIER_KIND")
	v.BindEnv("notifier.sendgrid.apiKey", "BACKGROUND_NOTIFIER_SENDGRID_API_KEY")
	v.BindEnv("notifier.sendgrid.incidentTemplateID", "BACKGROUND_NOTIFIER_SENDGRID_INCIDENT_TEMPLATE_ID")
	v.BindEnv("notifier.sendgrid.senderEmail", "BACKGROUND_NOTIFIER_SENDGRID_SENDER_EMAIL")
}

type RunnerConfigFile struct {
	Provisioner RunnerConfigFileProvisioner `mapstructure:"provisioner"`

	RunnerGRPCServerAddress string `mapstructure:"grpcServerAddress" default:"http://localhost:8080"`
}

type RunnerConfigFileProvisioner struct {
	Kind string `mapstructure:"kind" default:"local"`
}

type RunnerConfig struct {
	shared.Config

	DefaultProvisioner provisioner.Provisioner
}

func BindAllRunnerEnv(v *viper.Viper) {
	v.BindEnv("provisioner.kind", "RUNNER_WORKER_PROVISIONER_KIND")
	v.BindEnv("grpcServerAddress", "RUNNER_WORKER_GRPC_SERVER_ADDRESS")
}
