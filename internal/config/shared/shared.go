package shared

import (
	"github.com/hatchet-dev/hatchet/api/serverutils/erroralerter"
	"github.com/hatchet-dev/hatchet/internal/logger"
	"github.com/hatchet-dev/hatchet/internal/temporal"
)

// ConfigFile is shared by ALL hatchet workloads
type ConfigFile struct {
	// Debug is whether to print out debug lines
	Debug bool `env:"DEBUG,default=false"`

	// Temporal config options
	TemporalBroadcastAddress string `env:"TEMPORAL_BROADCAST_ADDRESS,default=http://127.0.0.1:7233"`
	TemporalEnabled          bool   `env:"TEMPORAL_ENABLED,default=true"`
	TemporalHostPort         string `env:"TEMPORAL_HOST_PORT,default=127.0.0.1:7233"`
	TemporalNamespace        string `env:"TEMPORAL_NAMESPACE,default=default"`
	TemporalBearerToken      string `env:"TEMPORAL_BEARER_TOKEN"`

	// TLS options
	TemporalClientTLSRootCAFile string `env:"TEMPORAL_CLIENT_TLS_ROOT_CA_FILE"`
	TemporalClientTLSCertFile   string `env:"TEMPORAL_CLIENT_TLS_CERT_FILE"`
	TemporalClientTLSKeyFile    string `env:"TEMPORAL_CLIENT_TLS_KEY_FILE"`
	TemporalTLSServerName       string `env:"TEMPORAL_TLS_SERVER_NAME"`
}

type Config struct {
	Logger       logger.Logger
	ErrorAlerter erroralerter.Alerter

	TemporalClient *temporal.Client
}

// RedisConfigFile is used for initiating Redis connections
type RedisConfigFile struct {
	RedisHost     string `env:"REDIS_HOST,default=redis"`
	RedisPort     string `env:"REDIS_PORT,default=6379"`
	RedisUsername string `env:"REDIS_USER"`
	RedisPassword string `env:"REDIS_PASS"`
	RedisDB       int    `env:"REDIS_DB,default=0"`
}

// FileStorageConfigFile is used for setting up a file storage backend
type FileStorageConfigFile struct {
	S3StateAWSAccessKeyID string `env:"S3_AWS_ACCESS_KEY_ID"`
	S3StateAWSSecretKey   string `env:"S3_AWS_SECRET_KEY"`
	S3StateAWSRegion      string `env:"S3_AWS_REGION"`
	S3StateBucketName     string `env:"S3_BUCKET_NAME"`
	S3StateEncryptionKey  string `env:"S3_ENCRYPTION_KEY,default=__random_strong_encryption_key__"`
}
