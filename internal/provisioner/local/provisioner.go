package local

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hatchet-dev/hatchet/internal/provisioner"
)

type LocalProvisioner struct {
}

func NewLocalProvisioner() *LocalProvisioner {
	return &LocalProvisioner{}
}

func (l *LocalProvisioner) RunPlan(opts *provisioner.ProvisionOpts) error {
	go func() {
		cmdProv := exec.Command("./bin/hatchet-runner", "plan")
		cmdProv.Stdout = os.Stdout
		cmdProv.Stderr = os.Stderr

		env := opts.Env
		env = append(env, cmdProv.Environ()...)

		env = append(env, "PATH=/usr/local/bin:/usr/bin:/bin")

		cmdProv.Env = env

		err := cmdProv.Run()

		if err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func (l *LocalProvisioner) RunApply(opts *provisioner.ProvisionOpts) error {
	go func() {
		cmdProv := exec.Command("./bin/hatchet-runner", "apply")
		cmdProv.Stdout = os.Stdout
		cmdProv.Stderr = os.Stderr
		env := opts.Env
		env = append(env, cmdProv.Environ()...)

		env = append(env, "PATH=/usr/local/bin:/usr/bin:/bin")

		cmdProv.Env = env

		err := cmdProv.Run()

		if err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func (l *LocalProvisioner) RunDestroy(opts *provisioner.ProvisionOpts) error {
	return nil
}