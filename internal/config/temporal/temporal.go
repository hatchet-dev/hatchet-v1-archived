package temporal

type TemporalConfigFile struct {
	TemporalAddress      string `env:"TEMPORAL_ADDRESS,default=127.0.0.1"`
	TemporalFrontendPort int64  `env:"TEMPORAL_FRONTEND_PORT,default=7233"`

	TemporalBroadcastAddress string `env:"TEMPORAL_BROADCAST_ADDRESS,default=127.0.0.1"`

	TemporalUIEnabled bool   `env:"TEMPORAL_UI_ENABLED,default=true"`
	TemporalUIAddress string `env:"TEMPORAL_UI_ADDRESS,default=127.0.0.1"`
	TemporalUIPort    int64  `env:"TEMPORAL_UI_PORT,default=8233"`

	TemporalPProfPort int64 `env:"TEMPORAL_PPROF_PORT,default=9500"`

	TemporalMetricsAddress string `env:"TEMPORAL_METRICS_ADDRESS,default=127.0.0.1"`
	TemporalMetricsPort    int64  `env:"TEMPORAL_METRICS_PORT,default=10001"`

	TemporalLogLevel string `env:"TEMPORAL_LOG_LEVEL,default=info"`

	TemporalSQLLitePath string `env:"TEMPORAL_SQL_LITE_PATH,default=/hatchet/temporal.db"`

	TemporalNamespaces []string `env:"TEMPORAL_NAMESPACES,default=default"`
}

type Config struct {
	ConfigFile *TemporalConfigFile
}
