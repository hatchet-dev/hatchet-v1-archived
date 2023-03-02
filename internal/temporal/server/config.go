package server

import (
	"fmt"
	"time"

	"go.temporal.io/server/common/cluster"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/metrics"
	"go.temporal.io/server/common/persistence/sql/sqlplugin/sqlite"
	schemasqlite "go.temporal.io/server/schema/sqlite"

	temporalconfig "github.com/hatchet-dev/hatchet/internal/config/temporal"
)

const (
	PersistenceStoreName = "sqlite-default"
)

func GetTemporalServerConfig(configfile *temporalconfig.TemporalConfigFile) (*config.Config, error) {
	baseConfig := &config.Config{}

	sqlConfig, err := getSQLConfig(configfile)

	if err != nil {
		return nil, err
	}

	baseConfig.Global.Membership = config.Membership{
		MaxJoinDuration:  30 * time.Second,
		BroadcastAddress: configfile.TemporalBroadcastAddress,
	}

	baseConfig.Global.Metrics = &metrics.Config{
		Prometheus: &metrics.PrometheusConfig{
			ListenAddress: fmt.Sprintf("%s:%d", configfile.TemporalMetricsAddress, configfile.TemporalMetricsPort),
			HandlerPath:   "/metrics",
		},
	}

	baseConfig.Global.PProf = config.PProf{Port: int(configfile.TemporalPProfPort)}

	baseConfig.Persistence = config.Persistence{
		DefaultStore:     PersistenceStoreName,
		VisibilityStore:  PersistenceStoreName,
		NumHistoryShards: 1,
		DataStores: map[string]config.DataStore{
			PersistenceStoreName: {SQL: sqlConfig},
		},
	}

	baseConfig.ClusterMetadata = &cluster.Config{
		EnableGlobalNamespace:    false,
		FailoverVersionIncrement: 10,
		MasterClusterName:        "active",
		CurrentClusterName:       "active",
		ClusterInformation: map[string]cluster.ClusterInformation{
			"active": {
				Enabled:                true,
				InitialFailoverVersion: 1,
				RPCAddress:             fmt.Sprintf("%s:%d", configfile.TemporalBroadcastAddress, configfile.TemporalFrontendPort),
			},
		},
	}

	baseConfig.DCRedirectionPolicy = config.DCRedirectionPolicy{
		Policy: "noop",
	}

	baseConfig.Services = map[string]config.Service{
		"frontend": getService(configfile, 0),
		"history":  getService(configfile, 1),
		"matching": getService(configfile, 2),
		"worker":   getService(configfile, 3),
	}

	baseConfig.Archival = config.Archival{
		History: config.HistoryArchival{
			State:      "disabled",
			EnableRead: false,
			Provider:   nil,
		},
		Visibility: config.VisibilityArchival{
			State:      "disabled",
			EnableRead: false,
			Provider:   nil,
		},
	}

	baseConfig.PublicClient = config.PublicClient{
		HostPort: fmt.Sprintf("%s:%d", configfile.TemporalBroadcastAddress, configfile.TemporalFrontendPort),
	}

	baseConfig.NamespaceDefaults = config.NamespaceDefaults{
		Archival: config.ArchivalNamespaceDefaults{
			History: config.HistoryArchivalNamespaceDefaults{
				State: "disabled",
			},
			Visibility: config.VisibilityArchivalNamespaceDefaults{
				State: "disabled",
			},
		},
	}

	return baseConfig, nil
}

func getSQLConfig(configfile *temporalconfig.TemporalConfigFile) (*config.SQL, error) {
	sqliteConfig := &config.SQL{
		PluginName:        sqlite.PluginName,
		ConnectAttributes: make(map[string]string),
		DatabaseName:      configfile.TemporalSQLLitePath,
	}

	sqliteConfig.ConnectAttributes["mode"] = "rwc"

	// no-err setup schema
	schemasqlite.SetupSchema(sqliteConfig)

	// Pre-create namespaces
	var namespaces []*schemasqlite.NamespaceConfig

	for _, ns := range configfile.TemporalNamespaces {
		namespaces = append(namespaces, schemasqlite.NewNamespaceConfig("active", ns, false))
	}

	if err := schemasqlite.CreateNamespaces(sqliteConfig, namespaces...); err != nil {
		return nil, fmt.Errorf("error creating namespaces: %w", err)
	}

	return sqliteConfig, nil
}

func getService(configfile *temporalconfig.TemporalConfigFile, frontendPortOffset int) config.Service {
	return config.Service{
		RPC: config.RPC{
			GRPCPort:        int(configfile.TemporalFrontendPort) + frontendPortOffset,
			MembershipPort:  int(configfile.TemporalFrontendPort) + 100 + frontendPortOffset,
			BindOnLocalHost: false,
			BindOnIP:        configfile.TemporalAddress,
		},
	}
}
