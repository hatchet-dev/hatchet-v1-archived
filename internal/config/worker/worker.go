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
)

type BackgroundConfigFile struct {
	ServerURL      string   `env:"SERVER_URL,default=https://hatchet.run"`
	TokenIssuerURL string   `env:"TOKEN_ISSUER_URL"`
	TokenAudience  []string `env:"TOKEN_AUDIENCE"`

	S3StateStore shared.FileStorageConfigFile

	RedisLogStore shared.RedisConfigFile

	// Notification options

	// RestrictedEmailDomains sets the restricted email domains for the instance.
	RestrictedEmailDomains []string `env:"RESTRICTED_EMAIL_DOMAINS"`

	// Sendgrid notifier options
	SendgridAPIKey             string `env:"SENDGRID_API_KEY"`
	SendgridIncidentTemplateID string `env:"SENDGRID_INCIDENT_TEMPLATE_ID"`
	SendgridSenderEmail        string `env:"SENDGRID_SENDER_EMAIL"`
}

type BackgroundConfig struct {
	shared.Config

	DB database.Config

	ServerURL string

	TokenOpts *token.TokenOpts

	DefaultFileStore filestorage.FileStorageManager

	DefaultLogStore logstorage.LogStorageBackend

	ModuleRunQueueManager queuemanager.ModuleRunQueueManager

	IncidentNotifier notifier.IncidentNotifier
}

type RunnerConfigFile struct {
	// Provisioner config options
	ProvisionerRunnerMethod string `env:"PROVISIONER_RUNNER_METHOD,default=local"`
	RunnerGRPCServerAddress string `env:"RUNNER_GRPC_SERVER_ADDRESS,default=http://localhost:8080"`
}

type RunnerConfig struct {
	shared.Config

	DefaultProvisioner provisioner.Provisioner
}
