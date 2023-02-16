package modulerunner

import (
	"context"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/config/worker"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
)

type ModuleRunner struct {
	conf *worker.Config
}

func NewModuleRunner(config *worker.Config) *ModuleRunner {
	return &ModuleRunner{config}
}

type RunInput struct {
	Kind models.ModuleRunKind
	Opts *provisioner.ProvisionOpts
}

func (mr *ModuleRunner) Run(ctx context.Context, input RunInput) (string, error) {
	// call provisioner
	var err error

	switch input.Kind {
	case models.ModuleRunKindApply:
		err = mr.conf.DefaultProvisioner.RunApply(input.Opts)
	case models.ModuleRunKindPlan:
		err = mr.conf.DefaultProvisioner.RunPlan(input.Opts)
	case models.ModuleRunKindDestroy:
		err = mr.conf.DefaultProvisioner.RunDestroy(input.Opts)
	default:
		return "", fmt.Errorf("not a supported run type")
	}

	if err != nil {
		return "", err
	}

	return "run_successful", nil
}
