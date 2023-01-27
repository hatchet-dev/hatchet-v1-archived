package action

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/hatchet-dev/hatchet/internal/config/runner"
)

type RunnerAction struct {
	writer io.Writer
}

func NewRunnerAction(writer io.Writer) *RunnerAction {
	return &RunnerAction{writer}
}

func (r *RunnerAction) Apply(
	config *runner.Config,
	vals map[string]interface{},
) ([]byte, error) {
	if !commandExists("terraform") {
		return nil, fmt.Errorf("terraform cli command does not exist")
	}

	// re initialize
	err := r.reInit(config)

	if err != nil {
		return nil, err
	}

	err = r.apply(config)

	if err != nil {
		return nil, err
	}

	// get the output
	return r.output(config)
}

func (r *RunnerAction) Plan(
	config *runner.Config,
	vals map[string]interface{},
) ([]byte, error) {
	if !commandExists("terraform") {
		return nil, fmt.Errorf("terraform cli command does not exist")
	}

	// re initialize
	err := r.reInit(config)

	if err != nil {
		return nil, err
	}

	return r.plan(config)
}

func (r *RunnerAction) reInit(config *runner.Config) error {
	cmd := exec.Command("terraform", "init", "-reconfigure", "-upgrade")
	cmd.Dir = config.TerraformConf.TFDir
	cmd.Stdout = r.writer
	cmd.Stderr = r.writer
	cmd.Stdin = strings.NewReader("yes\n")

	err := r.setBackendEnv(config, cmd)

	if err != nil {
		return err
	}

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func (r *RunnerAction) setBackendEnv(config *runner.Config, cmd *exec.Cmd) error {
	rc := config.ConfigFile

	tfStateAddress := fmt.Sprintf("%s/api/v1/teams/%s/modules/%s/runs/%s/tfstate",
		rc.APIServerAddress,
		rc.TeamID,
		rc.ModuleID,
		rc.ModuleRunID)

	cmd.Env = append(cmd.Environ(), []string{
		fmt.Sprintf("TF_HTTP_USERNAME=mrt"),
		fmt.Sprintf("TF_HTTP_PASSWORD=%s", rc.APIToken),
		fmt.Sprintf("TF_HTTP_ADDRESS=%s", tfStateAddress),
		fmt.Sprintf("TF_HTTP_LOCK_ADDRESS=%s", tfStateAddress),
		fmt.Sprintf("TF_HTTP_UNLOCK_ADDRESS=%s", tfStateAddress),
	}...)

	return nil
}

func (r *RunnerAction) apply(
	config *runner.Config,
) error {
	args := []string{"apply", "-json", "-auto-approve"}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = config.TerraformConf.TFDir

	// writer := io.MultiWriter(streamer, os.Stdout, os.Stderr)
	cmd.Stdout = r.writer
	cmd.Stderr = r.writer

	err := r.setBackendEnv(config, cmd)

	if err != nil {
		return err
	}

	return cmd.Run()
}

func (r *RunnerAction) plan(
	config *runner.Config,
) ([]byte, error) {
	args := []string{"plan"}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = config.TerraformConf.TFDir

	err := r.setBackendEnv(config, cmd)

	if err != nil {
		return nil, err
	}

	return cmd.Output()
}

func (r *RunnerAction) output(config *runner.Config) ([]byte, error) {
	cmd := exec.Command("terraform", "output", "-json")
	cmd.Dir = config.TerraformConf.TFDir

	return cmd.Output()
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
