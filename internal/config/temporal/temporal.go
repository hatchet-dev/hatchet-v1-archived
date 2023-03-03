package temporal

import "github.com/hatchet-dev/hatchet/internal/auth/token"

type TemporalConfigFile struct {
	TemporalPublicURL string `env:"TEMPORAL_PUBLIC_URL,default=http:127.0.0.1:7233"`

	TemporalAddress               string `env:"TEMPORAL_ADDRESS,default=127.0.0.1"`
	TemporalFrontendPort          int64  `env:"TEMPORAL_FRONTEND_PORT,default=7233"`
	TemporalFrontendTLSServerName string `env:"TEMPORAL_FRONTEND_TLS_SERVER_NAME"`
	TemporalFrontendTLSCertFile   string `env:"TEMPORAL_FRONTEND_TLS_CERT_FILE"`
	TemporalFrontendTLSKeyFile    string `env:"TEMPORAL_FRONTEND_TLS_KEY_FILE"`
	TemporalFrontendTLSRootCAFile string `env:"TEMPORAL_FRONTEND_TLS_ROOT_CA_FILE"`

	TemporalWorkerTLSServerName string `env:"TEMPORAL_WORKER_TLS_SERVER_NAME"`
	TemporalWorkerTLSCertFile   string `env:"TEMPORAL_WORKER_TLS_CERT_FILE"`
	TemporalWorkerTLSKeyFile    string `env:"TEMPORAL_WORKER_TLS_KEY_FILE"`
	TemporalWorkerTLSRootCAFile string `env:"TEMPORAL_WORKER_TLS_ROOT_CA_FILE"`

	TemporalInternodeTLSServerName string `env:"TEMPORAL_INTERNODE_TLS_SERVER_NAME"`
	TemporalInternodeTLSCertFile   string `env:"TEMPORAL_INTERNODE_TLS_CERT_FILE"`
	TemporalInternodeTLSKeyFile    string `env:"TEMPORAL_INTERNODE_TLS_KEY_FILE"`
	TemporalInternodeTLSRootCAFile string `env:"TEMPORAL_INTERNODE_TLS_ROOT_CA_FILE"`

	TemporalBroadcastAddress string `env:"TEMPORAL_BROADCAST_ADDRESS,default=127.0.0.1"`

	TemporalUIEnabled       bool   `env:"TEMPORAL_UI_ENABLED,default=true"`
	TemporalUIAddress       string `env:"TEMPORAL_UI_ADDRESS,default=127.0.0.1"`
	TemporalUIPort          int64  `env:"TEMPORAL_UI_PORT,default=8233"`
	TemporalUITLSServerName string `env:"TEMPORAL_UI_TLS_SERVER_NAME"`
	TemporalUITLSCertFile   string `env:"TEMPORAL_UI_TLS_CERT_FILE"`
	TemporalUITLSKeyFile    string `env:"TEMPORAL_UI_TLS_KEY_FILE"`
	TemporalUITLSRootCAFile string `env:"TEMPORAL_UI_TLS_ROOT_CA_FILE"`

	TemporalPProfPort int64 `env:"TEMPORAL_PPROF_PORT,default=9500"`

	TemporalMetricsAddress string `env:"TEMPORAL_METRICS_ADDRESS,default=127.0.0.1"`
	TemporalMetricsPort    int64  `env:"TEMPORAL_METRICS_PORT,default=10001"`

	TemporalLogLevel string `env:"TEMPORAL_LOG_LEVEL,default=info"`

	TemporalSQLLitePath string `env:"TEMPORAL_SQL_LITE_PATH,default=/hatchet/temporal.db"`

	TemporalNamespaces []string `env:"TEMPORAL_NAMESPACES,default=default"`

	TemporalInternalNamespace  string `env:"TEMPORAL_INTERNAL_NAMESPACE,default=hatchet-internal"`
	TemporalInternalSigningKey string `env:"TEMPORAL_INTERNAL_SIGNING_KEY,default=__random_strong_encryption_key__"`

	// TemporalInternalTokenIssuerURL is the endpoint of the issuer, typically equivalent to the server URL.
	// This field should INCLUDE the protocol.
	// If this is not set, it is set to the TEMPORAL_PUBLIC_URL variable.
	TemporalInternalTokenIssuerURL string `env:"TOKEN_ISSUER_URL"`

	// TemporalInternalTokenAudience is the set of audiences for the JWT token issuer, typically equivalent to the server URL.
	// This field should INCLUDE the protocol.
	// If this is not set, it is set to the TEMPORAL_PUBLIC_URL variable.
	TemporalInternalTokenAudience []string `env:"TOKEN_AUDIENCE"`
}

const DefaultInternalNamespace string = "hatchet-internal"

type InternalAuthConfig struct {
	InternalNamespace  string
	InternalSigningKey []byte
	InternalTokenOpts  token.TokenOpts
}

type Config struct {
	ConfigFile *TemporalConfigFile

	InternalAuthConfig *InternalAuthConfig
}
