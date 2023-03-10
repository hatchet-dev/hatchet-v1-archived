package shared

import (
	"github.com/hatchet-dev/hatchet/api/serverutils/erroralerter"
	"github.com/hatchet-dev/hatchet/internal/logger"
	"github.com/hatchet-dev/hatchet/internal/temporal"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	Debug bool `mapstructure:"debug" default:"false"`

	Temporal ConfigFileTemporal `mapstructure:"temporal"`
}

type ConfigFileTemporal struct {
	Client ConfigFileTemporalClient `mapstructure:"client"`
}

type ConfigFileTemporalClient struct {
	// Temporal config options
	TemporalBroadcastAddress string `mapstructure:"broadcastAddress" default:"http://127.0.0.1:7233"`
	TemporalEnabled          bool   `mapstructure:"enabled" default:"true"`
	TemporalHostPort         string `mapstructure:"hostPort" default:"127.0.0.1:7233"`
	TemporalNamespace        string `mapstructure:"namespace" default:"default"`
	TemporalBearerToken      string `mapstructure:"bearerToken"`

	// TLS options
	TemporalClientTLSRootCAFile string `mapstructure:"tlsRootCAFile"`
	TemporalClientTLSCertFile   string `mapstructure:"tlsCertFile"`
	TemporalClientTLSKeyFile    string `mapstructure:"tlsKeyFile"`
	TemporalTLSServerName       string `mapstructure:"tlsServerName"`
}

type LogStoreConfigFile struct {
	LogStorageKind string `mapstructure:"kind" default:"file"`

	Redis LogStoreConfigFileRedis `mapstructure:"redis"`

	File LogStoreConfigFileDirectory `mapstructure:"file"`
}

type LogStoreConfigFileRedis struct {
	RedisHost     string `mapstructure:"host" default:"redis"`
	RedisPort     string `mapstructure:"port" default:"6379"`
	RedisUsername string `mapstructure:"user"`
	RedisPassword string `mapstructure:"password"`
	RedisDB       int    `mapstructure:"db" default:"0"`
}

type LogStoreConfigFileDirectory struct {
	FileDirectory string `mapstructure:"directory" default:"./tmp/logs"`
}

type FileStorageConfigFile struct {
	FileStorageKind string `mapstructure:"kind" default:"local" validator:"oneof=s3 local"`

	Local FileStorageConfigFileLocal `mapstructure:"local"`

	S3 FileStorageConfigFileS3 `mapstructure:"s3"`
}

type FileStorageConfigFileLocal struct {
	FileDirectory     string `mapstructure:"directory" default:"./tmp/files"`
	FileEncryptionKey string `mapstructure:"encryptionKey" default:"__random_strong_encryption_key__"`
}

type FileStorageConfigFileS3 struct {
	S3StateAWSAccessKeyID string `mapstructure:"accessKeyID"`
	S3StateAWSSecretKey   string `mapstructure:"secretKey"`
	S3StateAWSRegion      string `mapstructure:"region"`
	S3StateBucketName     string `mapstructure:"bucketName"`
	S3StateEncryptionKey  string `mapstructure:"encryptionKey" default:"__random_strong_encryption_key__"`
}

type ConfigFileAuth struct {
	// RestrictedEmailDomains sets the restricted email domains for the instance.
	RestrictedEmailDomains []string `mapstructure:"restrictedEmailDomains"`

	// BasedAuthEnabled controls whether email and password-based login is enabled for this
	// Hatchet instances
	BasicAuthEnabled bool `mapstructure:"basicAuthEnabled" default:"true"`

	// Configuration options for the cookie
	Cookie ConfigFileAuthCookie `mapstructure:"cookie"`

	// Configuration options for the token
	Token ConfigFileAuthToken `mapstructure:"token"`
}

type ConfigFileAuthCookie struct {
	Name     string   `mapstructure:"name" default:"hatchet"`
	Domain   string   `mapstructure:"domain"`
	Secrets  []string `mapstructure:"secrets" default:"[\"random_hash_key_\",\"random_block_key\"]"`
	Insecure bool     `mapstructure:"insecure" default:"false"`
}

type ConfigFileAuthToken struct {
	// TokenIssuerURL is the endpoint of the issuer, typically equivalent to the server URL.
	// This field should INCLUDE the protocol.
	// If this is not set, it is set to the SERVER_URL variable.
	TokenIssuerURL string `mapstructure:"issuer" validator:"url"`

	// TokenAudience is the set of audiences for the JWT token issuer, typically equivalent to the server URL.
	// This field should INCLUDE the protocol.
	// If this is not set, it is set to the SERVER_URL variable.
	TokenAudience []string `mapstructure:"audience" validator:"url"`
}

type Config struct {
	Logger       logger.Logger
	ErrorAlerter erroralerter.Alerter

	TemporalClient *temporal.Client
}

func BindAllEnv(v *viper.Viper) {
	v.BindEnv("debug", "DEBUG")
	v.BindEnv("temporal.client.broadcastAddress", "TEMPORAL_CLIENT_BROADCAST_ADDRESS")
	v.BindEnv("temporal.client.enabled", "TEMPORAL_CLIENT_ENABLED")
	v.BindEnv("temporal.client.hostPort", "TEMPORAL_CLIENT_HOST_PORT")
	v.BindEnv("temporal.client.namespace", "TEMPORAL_CLIENT_NAMESPACE")
	v.BindEnv("temporal.client.bearerToken", "TEMPORAL_CLIENT_BEARER_TOKEN")
	v.BindEnv("temporal.client.tlsRootCAFile", "TEMPORAL_CLIENT_TLS_ROOT_CA_FILE")
	v.BindEnv("temporal.client.tlsCertFile", "TEMPORAL_CLIENT_TLS_CERT_FILE")
	v.BindEnv("temporal.client.tlsKeyFile", "TEMPORAL_CLIENT_TLS_KEY_FILE")
	v.BindEnv("temporal.client.tlsServerName", "TEMPORAL_CLIENT_TLS_SERVER_NAME")
}
