package action

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/antihax/optional"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
)

type actionHandler func(config *runner.Config, description string) error

type RunnerAction struct {
	writer     io.Writer
	errHandler actionHandler
}

func NewRunnerAction(writer io.Writer, errHandler actionHandler) *RunnerAction {
	return &RunnerAction{writer, errHandler}
}

func (r *RunnerAction) Apply(
	config *runner.Config,
	vals map[string]interface{},
) ([]byte, error) {
	if !commandExists("terraform") {
		return nil, fmt.Errorf("terraform cli command does not exist")
	}

	var planPath string

	// download plan, if github commit sha is passed in
	if config.ConfigFile.GithubSHA != "" {
		resp, _, err := config.FileClient.GetPlanByCommitSHA(
			config.ConfigFile.TeamID,
			config.ConfigFile.ModuleID,
			config.ConfigFile.ModuleRunID,
		)

		if resp != nil {
			defer resp.Close()
		}

		if err != nil {
			return nil, err
		}

		fileBytes, err := ioutil.ReadAll(resp)

		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(filepath.Join(config.TerraformConf.TFDir, "./plan.tfplan"), fileBytes, 0666)

		if err != nil {
			return nil, err
		}

		planPath = "./plan.tfplan"
	}

	err := r.downloadModuleValues(config, "./tfvars.json")

	if err != nil {
		return nil, err
	}

	// re initialize
	err = r.reInit(config)

	if err != nil {
		return nil, err
	}

	err = r.apply(config, planPath, "./tfvars.json")

	if err != nil {
		return nil, err
	}

	// get the output
	return r.output(config)
}

func (r *RunnerAction) Plan(
	config *runner.Config,
	vals map[string]interface{},
) ([]byte, []byte, []byte, error) {
	if !commandExists("terraform") {
		return nil, nil, nil, fmt.Errorf("terraform cli command does not exist")
	}

	err := r.downloadModuleValues(config, "./tfvars.json")

	if err != nil {
		return nil, nil, nil, err
	}

	// re initialize
	err = r.reInit(config)

	if err != nil {
		return nil, nil, nil, r.errHandler(config, fmt.Sprintf("Failed while reinitializing the Terraform backend: %s", err.Error()))
	}

	err = r.plan(config, "./tfvars.json")

	if err != nil {
		return nil, nil, nil, r.errHandler(config, fmt.Sprintf("Failed while running plan: %s", err.Error()))
	}

	zipOut, err := r.getPlanZIP(config)

	if err != nil {
		return nil, nil, nil, r.errHandler(config, fmt.Sprintf("Failed while getting zip output: %s", err.Error()))
	}

	prettyOut, err := r.showPretty(config)

	if err != nil {
		return nil, nil, nil, r.errHandler(config, fmt.Sprintf("Failed while generating prettified output: %s", err.Error()))
	}

	jsonOut, err := r.showJSON(config)

	if err != nil {
		return nil, nil, nil, r.errHandler(config, fmt.Sprintf("Failed while generating JSON output: %s", err.Error()))
	}

	return zipOut, prettyOut, jsonOut, nil
}

func (r *RunnerAction) downloadModuleValues(config *runner.Config, relPath string) error {
	// download values
	vals, _, err := config.APIClient.ModulesApi.GetModuleValues(
		context.Background(),
		config.ConfigFile.TeamID,
		config.ConfigFile.ModuleID,
		&swagger.ModulesApiGetModuleValuesOpts{
			GithubSha: optional.NewString(config.ConfigFile.GithubSHA),
		},
	)

	if err != nil {
		return err
	}

	fileBytes, err := json.Marshal(vals)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(config.TerraformConf.TFDir, relPath), fileBytes, 0666)

	return err
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
	planPath string,
	valsFilePath string,
) error {
	args := []string{"apply", "-json", "-auto-approve"}

	if planPath != "" {
		args = append(args, fmt.Sprintf("%s", planPath))
	}

	if valsFilePath != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", valsFilePath))
	}

	fmt.Printf("running apply with args: [%v]", args)

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
	valsFilePath string,
) error {
	args := []string{"plan", "-out=./plan.tfplan"}

	if valsFilePath != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", valsFilePath))
	}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = config.TerraformConf.TFDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := r.setBackendEnv(config, cmd)

	if err != nil {
		return err
	}

	return cmd.Run()
}

func (r *RunnerAction) getPlanZIP(
	config *runner.Config,
) ([]byte, error) {
	path := filepath.Join(config.TerraformConf.TFDir, "./plan.tfplan")
	return ioutil.ReadFile(path)
}

func (r *RunnerAction) showPretty(
	config *runner.Config,
) ([]byte, error) {
	args := []string{"show", "-no-color", "./plan.tfplan"}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = config.TerraformConf.TFDir
	cmd.Stderr = os.Stderr

	err := r.setBackendEnv(config, cmd)

	if err != nil {
		return nil, err
	}

	return cmd.Output()
}

func (r *RunnerAction) showJSON(
	config *runner.Config,
) ([]byte, error) {
	args := []string{"show", "-json", "./plan.tfplan"}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = config.TerraformConf.TFDir
	cmd.Stderr = os.Stderr

	err := r.setBackendEnv(config, cmd)

	if err != nil {
		return nil, err
	}

	return cmd.Output()
}

func (r *RunnerAction) output(config *runner.Config) ([]byte, error) {
	cmd := exec.Command("terraform", "output", "-json")
	cmd.Dir = config.TerraformConf.TFDir
	cmd.Stderr = os.Stderr

	err := r.setBackendEnv(config, cmd)

	if err != nil {
		return nil, err
	}

	return cmd.Output()
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
