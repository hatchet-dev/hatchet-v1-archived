package worker

import (
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
	"github.com/hatchet-dev/hatchet/internal/temporal"
)

type ConfigFile struct {
	// Provisioner config options
	ProvisionerRunnerMethod string `env:"PROVISIONER_RUNNER_METHOD,default=local"`
	RunnerGRPCServerAddress string `env:"RUNNER_GRPC_SERVER_ADDRESS,default=http://localhost:8080"`

	// Temporal config options
	TemporalEnabled       bool   `env:"TEMPORAL_ENABLED,default=true"`
	TemporalRunWorkers    bool   `env:"TEMPORAL_RUN_WORKERS,default=true"`
	TemporalHostPort      string `env:"TEMPORAL_HOST_PORT,default=127.0.0.1:7233"`
	TemporalNamespace     string `env:"TEMPORAL_NAMESPACE,default=default"`
	TemporalAuthHeaderKey string `env:"TEMPORAL_AUTH_HEADER_KEY"`
	TemporalAuthHeaderVal string `env:"TEMPORAL_AUTH_HEADER_VAL"`
}

type Config struct {
	shared.Config

	DB database.Config

	DefaultProvisioner provisioner.Provisioner

	TemporalClient *temporal.Client
}
