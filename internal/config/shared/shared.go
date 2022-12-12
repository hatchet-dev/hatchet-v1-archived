package shared

import (
	"github.com/hatchet-dev/hatchet/api/serverutils/erroralerter"
	"github.com/hatchet-dev/hatchet/internal/logger"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type ConfigFile struct {
	// Debug is whether to print out debug lines
	Debug bool `env:"DEBUG,default=false"`
}

type Config struct {
	Logger       logger.Logger
	ErrorAlerter erroralerter.Alerter
	Repository   repository.Repository
}
