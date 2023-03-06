package server

import (
	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/temporal"

	temporalconfig "github.com/hatchet-dev/hatchet/internal/config/temporal"
	"github.com/hatchet-dev/hatchet/internal/temporal/server/authorizer"
)

func NewTemporalServer(tconfig *temporalconfig.Config) (temporal.Server, error) {
	configfile := tconfig.ConfigFile

	logger := log.NewZapLogger(log.BuildZapLogger(log.Config{
		Stdout:     true,
		Level:      configfile.TemporalLogLevel,
		OutputFile: "",
	}))

	cfg, err := GetTemporalServerConfig(tconfig)

	if err != nil {
		return nil, err
	}

	authorizerAndClaimMapper := authorizer.NewHatchetAuthorizer(tconfig, &cfg.Global.Authorization, logger)

	return temporal.NewServer(
		temporal.ForServices(temporal.Services),
		temporal.WithConfig(cfg),
		temporal.WithLogger(logger),
		temporal.InterruptOn(temporal.InterruptCh()),
		temporal.WithAuthorizer(authorizerAndClaimMapper),
		temporal.WithClaimMapper(func(cfg *config.Config) authorization.ClaimMapper {
			return authorizerAndClaimMapper
		}),
	)
}
