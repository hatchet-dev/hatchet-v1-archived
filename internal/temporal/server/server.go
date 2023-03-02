package server

import (
	"fmt"
	goLog "log"

	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/temporal"

	temporalconfig "github.com/hatchet-dev/hatchet/internal/config/temporal"
)

func NewTemporalServer(configfile *temporalconfig.TemporalConfigFile) (temporal.Server, error) {
	logger := log.NewZapLogger(log.BuildZapLogger(log.Config{
		Stdout:     true,
		Level:      configfile.TemporalLogLevel,
		OutputFile: "",
	}))

	cfg, err := GetTemporalServerConfig(configfile)

	if err != nil {
		return nil, err
	}

	authorizer, err := authorization.GetAuthorizerFromConfig(
		&cfg.Global.Authorization,
	)

	if err != nil {
		goLog.Fatal(fmt.Sprintf("Unable to instantiate authorizer. Error: %v", err))
	}

	claimMapper, err := authorization.GetClaimMapperFromConfig(&cfg.Global.Authorization, logger)

	if err != nil {
		goLog.Fatal(fmt.Sprintf("Unable to instantiate claim mapper: %v.", err))
	}

	return temporal.NewServer(
		temporal.ForServices(temporal.Services),
		temporal.WithConfig(cfg),
		temporal.WithLogger(logger),
		temporal.InterruptOn(temporal.InterruptCh()),
		temporal.WithAuthorizer(authorizer),
		temporal.WithClaimMapper(func(cfg *config.Config) authorization.ClaimMapper {
			return claimMapper
		}),
	)
}
