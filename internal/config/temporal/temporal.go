package temporal

import (
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/spf13/viper"
)

type TemporalConfigFile struct {
	TemporalPublicURL          string   `mapstructure:"publicURL" json:"publicURL,omitempty" default:"http://127.0.0.1:7233"`
	TemporalAddress            string   `mapstructure:"address" json:"address,omitempty" default:"127.0.0.1"`
	TemporalBroadcastAddress   string   `mapstructure:"broadcastAddress" json:"broadcastAddress,omitempty" default:"127.0.0.1"`
	TemporalPProfPort          int64    `mapstructure:"pprofPort" json:"pprofPort,omitempty" default:"9500"`
	TemporalMetricsAddress     string   `mapstructure:"metricsAddress" json:"metricsAddress,omitempty" default:"127.0.0.1"`
	TemporalMetricsPort        int64    `mapstructure:"metricsPort" json:"metricsPort,omitempty" default:"10001"`
	TemporalLogLevel           string   `mapstructure:"logLevel" json:"logLevel,omitempty" default:"info"`
	TemporalSQLLitePath        string   `mapstructure:"sqlLitePath" json:"sqlLitePath,omitempty" default:"/hatchet/temporal.db"`
	TemporalNamespaces         []string `mapstructure:"namespaces" json:"namespaces,omitempty" default:"[\"default\"]"`
	TemporalInternalNamespace  string   `mapstructure:"internalNamespace" json:"internalNamespace,omitempty" default:"hatchet-internal"`
	TemporalInternalSigningKey string   `mapstructure:"internalSigningKey" json:"internalSigningKey,omitempty"`

	Token shared.ConfigFileAuthToken `mapstructure:"token" json:"token,omitempty"`

	Frontend  TemporalConfigFileFrontend  `mapstructure:"frontend" json:"frontend,omitempty"`
	Worker    TemporalConfigFileWorker    `mapstructure:"worker" json:"worker,omitempty"`
	Internode TemporalConfigFileInternode `mapstructure:"internode" json:"internode,omitempty"`
	UI        TemporalConfigFileUI        `mapstructure:"ui" json:"ui,omitempty"`
}

type TemporalConfigFileFrontend struct {
	TemporalFrontendPort          int64  `mapstructure:"port" json:"port,omitempty" default:"7233"`
	TemporalFrontendTLSServerName string `mapstructure:"tlsServerName" json:"tlsServerName,omitempty"`
	TemporalFrontendTLSCertFile   string `mapstructure:"tlsCertFile" json:"tlsCertFile,omitempty"`
	TemporalFrontendTLSKeyFile    string `mapstructure:"tlsKeyFile" json:"tlsKeyFile,omitempty"`
	TemporalFrontendTLSRootCAFile string `mapstructure:"tlsRootCAFile" json:"tlsRootCAFile,omitempty"`
}

type TemporalConfigFileWorker struct {
	TemporalWorkerTLSServerName string `mapstructure:"tlsServerName" json:"tlsServerName,omitempty"`
	TemporalWorkerTLSCertFile   string `mapstructure:"tlsCertFile" json:"tlsCertFile,omitempty"`
	TemporalWorkerTLSKeyFile    string `mapstructure:"tlsKeyFile" json:"tlsKeyFile,omitempty"`
	TemporalWorkerTLSRootCAFile string `mapstructure:"tlsRootCAFile" json:"tlsRootCAFile,omitempty"`
}

type TemporalConfigFileInternode struct {
	TemporalInternodeTLSServerName string `mapstructure:"tlsServerName" json:"tlsServerName,omitempty"`
	TemporalInternodeTLSCertFile   string `mapstructure:"tlsCertFile" json:"tlsCertFile,omitempty"`
	TemporalInternodeTLSKeyFile    string `mapstructure:"tlsKeyFile" json:"tlsKeyFile,omitempty"`
	TemporalInternodeTLSRootCAFile string `mapstructure:"tlsRootCAFile" json:"tlsRootCAFile,omitempty"`
}

type TemporalConfigFileUI struct {
	TemporalUIEnabled       bool   `mapstructure:"enabled" json:"enabled,omitempty" default:"true"`
	TemporalUIAddress       string `mapstructure:"uiAddress" json:"uiAddress,omitempty" default:"127.0.0.1"`
	TemporalUIPort          int64  `mapstructure:"port" json:"port,omitempty" default:"8233"`
	TemporalUITLSServerName string `mapstructure:"tlsServerName" json:"tlsServerName,omitempty"`
	TemporalUITLSCertFile   string `mapstructure:"tlsCertFile" json:"tlsCertFile,omitempty"`
	TemporalUITLSKeyFile    string `mapstructure:"tlsKeyFile" json:"tlsKeyFile,omitempty"`
	TemporalUITLSRootCAFile string `mapstructure:"tlsRootCAFile" json:"tlsRootCAFile,omitempty"`
}

const DefaultInternalNamespace string = "hatchet-internal"

type InternalAuthConfig struct {
	InternalNamespace  string
	InternalSigningKey []byte
	InternalTokenOpts  token.TokenOpts
}

type Config struct {
	DB database.Config

	ConfigFile *TemporalConfigFile

	InternalAuthConfig *InternalAuthConfig
}

func BindAllEnv(v *viper.Viper) {
	v.BindEnv("publicURL", "TEMPORAL_PUBLIC_URL")
	v.BindEnv("address", "TEMPORAL_ADDRESS")
	v.BindEnv("broadcastAddress", "TEMPORAL_BROADCAST_ADDRESS")
	v.BindEnv("pprofPort", "TEMPORAL_PPROF_PORT")
	v.BindEnv("metricsAddress", "TEMPORAL_METRICS_ADDRESS")
	v.BindEnv("metricsPort", "TEMPORAL_METRICS_PORT")
	v.BindEnv("logLevel", "TEMPORAL_LOG_LEVEL")
	v.BindEnv("sqlLitePath", "TEMPORAL_SQLITE_PATH")
	v.BindEnv("namespaces", "TEMPORAL_NAMESPACES")
	v.BindEnv("internalNamespace", "TEMPORAL_INTERNAL_NAMESPACE")
	v.BindEnv("internalSigningKey", "TEMPORAL_INTERNAL_SIGNING_KEY")

	v.BindEnv("token.issuer", "TEMPORAL_TOKEN_ISSUER")
	v.BindEnv("token.audience", "TEMPORAL_TOKEN_AUDIENCE")

	v.BindEnv("frontend.port", "TEMPORAL_FRONTEND_PORT")
	v.BindEnv("frontend.tlsServerName", "TEMPORAL_FRONTEND_TLS_SERVER_NAME")
	v.BindEnv("frontend.tlsCertFile", "TEMPORAL_FRONTEND_TLS_CERT_FILE")
	v.BindEnv("frontend.tlsKeyFile", "TEMPORAL_FRONTEND_TLS_KEY_FILE")
	v.BindEnv("frontend.tlsRootCAFile", "TEMPORAL_FRONTEND_TLS_ROOT_CA_FILE")

	v.BindEnv("worker.tlsServerName", "TEMPORAL_WORKER_TLS_SERVER_NAME")
	v.BindEnv("worker.tlsCertFile", "TEMPORAL_WORKER_TLS_CERT_FILE")
	v.BindEnv("worker.tlsKeyFile", "TEMPORAL_WORKER_TLS_KEY_FILE")
	v.BindEnv("worker.tlsRootCAFile", "TEMPORAL_WORKER_TLS_ROOT_CA_FILE")

	v.BindEnv("internode.tlsServerName", "TEMPORAL_INTERNODE_TLS_SERVER_NAME")
	v.BindEnv("internode.tlsCertFile", "TEMPORAL_INTERNODE_TLS_CERT_FILE")
	v.BindEnv("internode.tlsKeyFile", "TEMPORAL_INTERNODE_TLS_KEY_FILE")
	v.BindEnv("internode.tlsRootCAFile", "TEMPORAL_INTERNODE_TLS_ROOT_CA_FILE")

	v.BindEnv("ui.enabled", "TEMPORAL_UI_ENABLED")
	v.BindEnv("ui.uiAddress", "TEMPORAL_UI_ADDRESS")
	v.BindEnv("ui.port", "TEMPORAL_UI_PORT")
	v.BindEnv("ui.tlsServerName", "TEMPORAL_UI_TLS_SERVER_NAME")
	v.BindEnv("ui.tlsCertFile", "TEMPORAL_UI_TLS_CERT_FILE")
	v.BindEnv("ui.tlsKeyFile", "TEMPORAL_UI_TLS_KEY_FILE")
	v.BindEnv("ui.tlsRootCAFile", "TEMPORAL_UI_TLS_ROOT_CA_FILE")
}
