package local

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hatchet-dev/hatchet/internal/provisioner"
)

type LocalProvisioner struct {
	binaryPath string
}

type LocalProvisionerOpts struct {
	BinaryPath string
}

func NewLocalProvisioner(opts *LocalProvisionerOpts) *LocalProvisioner {
	return &LocalProvisioner{opts.BinaryPath}
}

func (l *LocalProvisioner) RunPlan(opts *provisioner.ProvisionOpts) error {
	runFunc := l.getRunFunc(opts, "plan")

	if opts.WaitForRunFinished {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil
}

func (l *LocalProvisioner) RunApply(opts *provisioner.ProvisionOpts) error {
	runFunc := l.getRunFunc(opts, "apply")

	if opts.WaitForRunFinished {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil
}

func (l *LocalProvisioner) RunDestroy(opts *provisioner.ProvisionOpts) error {
	runFunc := l.getRunFunc(opts, "destroy")

	if opts.WaitForRunFinished {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil
}

func (l *LocalProvisioner) RunStateMonitor(opts *provisioner.ProvisionOpts, monitorID string, policy []byte) error {
	runFunc := l.getMonitorFunc(opts, monitorID, policy, "state")

	if opts.WaitForRunFinished {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil
}

func (l *LocalProvisioner) RunPlanMonitor(opts *provisioner.ProvisionOpts, monitorID string, policy []byte) error {
	runFunc := l.getMonitorFunc(opts, monitorID, policy, "plan")

	if opts.WaitForRunFinished {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil
}

func (l *LocalProvisioner) RunBeforePlanMonitor(opts *provisioner.ProvisionOpts, monitorID string, policy []byte) error {
	runFunc := l.getMonitorFunc(opts, monitorID, policy, "before-plan")

	if opts.WaitForRunFinished {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil
}

func (l *LocalProvisioner) RunAfterPlanMonitor(opts *provisioner.ProvisionOpts, monitorID string, policy []byte) error {
	runFunc := l.getMonitorFunc(opts, monitorID, policy, "after-plan")

	if opts.WaitForRunFinished {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil
}

func (l *LocalProvisioner) RunBeforeApplyMonitor(opts *provisioner.ProvisionOpts, monitorID string, policy []byte) error {
	runFunc := l.getMonitorFunc(opts, monitorID, policy, "before-apply")

	if opts.WaitForRunFinished {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil
}

func (l *LocalProvisioner) RunAfterApplyMonitor(opts *provisioner.ProvisionOpts, monitorID string, policy []byte) error {
	runFunc := l.getMonitorFunc(opts, monitorID, policy, "after-apply")

	if opts.WaitForRunFinished {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil
}

func (l *LocalProvisioner) getRunFunc(opts *provisioner.ProvisionOpts, arg string) func() error {
	return func() error {
		cmdProv := exec.Command(l.binaryPath, arg)
		cmdProv.Stdout = os.Stdout
		cmdProv.Stderr = os.Stderr

		env := opts.Env
		env = append(env, cmdProv.Environ()...)

		env = append(env, "PATH=/usr/local/bin:/usr/bin:/bin")

		cmdProv.Env = env

		err := cmdProv.Run()

		if err != nil && !opts.WaitForRunFinished {
			fmt.Println(err)
		}

		return err
	}
}

func (l *LocalProvisioner) getMonitorFunc(opts *provisioner.ProvisionOpts, monitorID string, policy []byte, arg string) func() error {
	return func() error {
		cmdProv := exec.Command(l.binaryPath, "monitor", arg)
		cmdProv.Stdout = os.Stdout
		cmdProv.Stderr = os.Stderr
		env := opts.Env
		env = append(env, cmdProv.Environ()...)

		env = append(env, fmt.Sprintf("RUNNER_MODULE_MONITOR_ID=%s", monitorID))

		env = append(env, "PATH=/usr/local/bin:/usr/bin:/bin")

		cmdProv.Env = env

		err := cmdProv.Run()

		if err != nil && !opts.WaitForRunFinished {
			fmt.Println(err)
		}

		return err
	}
}
