package server

import (
	"fmt"

	temporalconfig "github.com/hatchet-dev/hatchet/internal/config/temporal"
	uiserver "github.com/temporalio/ui-server/v2/server"
	uiconfig "github.com/temporalio/ui-server/v2/server/config"
	uiserveroptions "github.com/temporalio/ui-server/v2/server/server_options"
)

func NewUIServer(configfile *temporalconfig.TemporalConfigFile) (*uiserver.Server, error) {
	cfg := &uiconfig.Config{
		Host:                configfile.UI.TemporalUIAddress,
		Port:                int(configfile.UI.TemporalUIPort),
		TemporalGRPCAddress: fmt.Sprintf("%s:%d", configfile.TemporalBroadcastAddress, configfile.Frontend.TemporalFrontendPort),
		EnableUI:            true,
		TLS: uiconfig.TLS{
			CaFile:     configfile.UI.TemporalUITLSRootCAFile,
			CertFile:   configfile.UI.TemporalUITLSCertFile,
			KeyFile:    configfile.UI.TemporalUITLSKeyFile,
			ServerName: configfile.UI.TemporalUITLSServerName,
		},
	}

	return uiserver.NewServer(uiserveroptions.WithConfigProvider(cfg)), nil
}
