package shared

import (
	"github.com/hatchet-dev/hatchet/api/serverutils/erroralerter"
	"github.com/hatchet-dev/hatchet/internal/logger"
	"github.com/hatchet-dev/hatchet/internal/temporal"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	Debug bool `mapstructure:"debug" json:"debug,omitempty" default:"false"`

	Temporal ConfigFileTemporal `mapstructure:"temporal" json:"temporal,omitempty"`
}

type ConfigFileTemporal struct {
	Client ConfigFileTemporalClient `mapstructure:"client" json:"client,omitempty"`
}

type ConfigFileTemporalClient struct {
	// Temporal config options
	TemporalBroadcastAddress string `mapstructure:"broadcastAddress" json:"broadcastAddress,omitempty" default:"http://127.0.0.1:7233"`
	TemporalEnabled          bool   `mapstructure:"enabled" json:"enabled,omitempty" default:"true"`
	TemporalHostPort         string `mapstructure:"hostPort" json:"hostPort,omitempty" default:"127.0.0.1:7233"`
	TemporalNamespace        string `mapstructure:"namespace" json:"namespace,omitempty" default:"default"`
	TemporalBearerToken      string `mapstructure:"bearerToken" json:"bearerToken,omitempty"`

	// TLS options
	TemporalClientTLSRootCAFile string `mapstructure:"tlsRootCAFile" json:"tlsRootCAFile,omitempty"`
	TemporalClientTLSCertFile   string `mapstructure:"tlsCertFile" json:"tlsCertFile,omitempty"`
	TemporalClientTLSKeyFile    string `mapstructure:"tlsKeyFile" json:"tlsKeyFile,omitempty"`
	TemporalTLSServerName       string `mapstructure:"tlsServerName" json:"tlsServerName,omitempty"`
}

type LogStoreConfigFile struct {
	LogStorageKind string `mapstructure:"kind" json:"kind,omitempty" default:"file"`

	Redis LogStoreConfigFileRedis `mapstructure:"redis" json:"redis,omitempty"`

	File LogStoreConfigFileDirectory `mapstructure:"file" json:"file,omitempty"`
}

type LogStoreConfigFileRedis struct {
	RedisHost     string `mapstructure:"host" json:"host,omitempty" default:"redis"`
	RedisPort     string `mapstructure:"port" json:"port,omitempty" default:"6379"`
	RedisUsername string `mapstructure:"user" json:"user,omitempty"`
	RedisPassword string `mapstructure:"password" json:"password,omitempty"`
	RedisDB       int    `mapstructure:"db" json:"db,omitempty" default:"0"`
}

type LogStoreConfigFileDirectory struct {
	FileDirectory string `mapstructure:"directory" json:"directory,omitempty" default:"./tmp/logs"`
}

type FileStorageConfigFile struct {
	FileStorageKind string `mapstructure:"kind" json:"kind,omitempty" default:"local" validator:"oneof=s3 local"`

	Local FileStorageConfigFileLocal `mapstructure:"local" json:"local,omitempty"`

	S3 FileStorageConfigFileS3 `mapstructure:"s3" json:"s3,omitempty"`
}

type FileStorageConfigFileLocal struct {
	FileDirectory     string `mapstructure:"directory" json:"directory,omitempty" default:"./tmp/files"`
	FileEncryptionKey string `mapstructure:"encryptionKey" json:"encryptionKey,omitempty" default:"__random_strong_encryption_key__"`
}

type FileStorageConfigFileS3 struct {
	S3StateAWSAccessKeyID string `mapstructure:"accessKeyID" json:"accessKeyID,omitempty"`
	S3StateAWSSecretKey   string `mapstructure:"secretKey" json:"secretKey,omitempty"`
	S3StateAWSRegion      string `mapstructure:"region" json:"region,omitempty"`
	S3StateBucketName     string `mapstructure:"bucketName" json:"bucketName,omitempty"`
	S3StateEncryptionKey  string `mapstructure:"encryptionKey" json:"encryptionKey,omitempty" default:"__random_strong_encryption_key__"`
}

type ConfigFileAuth struct {
	// RestrictedEmailDomains sets the restricted email domains for the instance.
	RestrictedEmailDomains []string `mapstructure:"restrictedEmailDomains" json:"restrictedEmailDomains,omitempty"`

	// BasedAuthEnabled controls whether email and password-based login is enabled for this
	// Hatchet instances
	BasicAuthEnabled bool `mapstructure:"basicAuthEnabled" json:"basicAuthEnabled,omitempty" default:"true"`

	// Configuration options for the cookie
	Cookie ConfigFileAuthCookie `mapstructure:"cookie" json:"cookie,omitempty"`

	// Configuration options for the token
	Token ConfigFileAuthToken `mapstructure:"token" json:"token,omitempty"`
}

type ConfigFileAuthCookie struct {
	Name     string   `mapstructure:"name" json:"name,omitempty" default:"hatchet"`
	Domain   string   `mapstructure:"domain" json:"domain,omitempty"`
	Secrets  []string `mapstructure:"secrets" json:"secrets,omitempty" default:"[\"random_hash_key_\",\"random_block_key\"]"`
	Insecure bool     `mapstructure:"insecure" json:"insecure,omitempty" default:"false"`
}

type ConfigFileAuthToken struct {
	// TokenIssuerURL is the endpoint of the issuer, typically equivalent to the server URL.
	// This field should INCLUDE the protocol.
	// If this is not set, it is set to the SERVER_URL variable.
	TokenIssuerURL string `mapstructure:"issuer" json:"issuer,omitempty" validator:"url"`

	// TokenAudience is the set of audiences for the JWT token issuer, typically equivalent to the server URL.
	// This field should INCLUDE the protocol.
	// If this is not set, it is set to the SERVER_URL variable.
	TokenAudience []string `mapstructure:"audience" json:"audience,omitempty" validator:"url"`
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
