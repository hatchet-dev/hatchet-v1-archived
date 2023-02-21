package modulerunner

import (
	"context"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/config/worker"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
)

type ModuleRunner struct {
	conf *worker.RunnerConfig
}

func NewModuleRunner(config *worker.RunnerConfig) *ModuleRunner {
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

type MonitorInput struct {
	ModuleMonitorID string
	Kind            models.ModuleMonitorKind
	Opts            *provisioner.ProvisionOpts
}

type MonitorOutput struct {
	Status      models.MonitorResultStatus
	Severity    models.MonitorResultSeverity
	Title       string
	Description string
}

func (mr *ModuleRunner) Monitor(ctx context.Context, input MonitorInput) (string, error) {
	// call provisioner
	var err error

	switch input.Kind {
	case models.MonitorKindState:
		err = mr.conf.DefaultProvisioner.RunStateMonitor(input.Opts, input.ModuleMonitorID, nil)
	case models.MonitorKindPlan:
		err = mr.conf.DefaultProvisioner.RunPlanMonitor(input.Opts, input.ModuleMonitorID, nil)
	default:
		return "", fmt.Errorf("not a supported monitor type")
	}

	if err != nil {
		return "", err
	}

	// TODO: return MonitorOutput here!
	return "monitor_successful", nil
}
