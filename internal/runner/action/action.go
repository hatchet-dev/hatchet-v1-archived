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
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
	"github.com/hatchet-dev/hatchet/internal/opa"
)

type actionHandler func(config *runner.Config, reportKind, description string) error

type RunnerAction struct {
	writer     io.Writer
	errHandler actionHandler
	reportKind string
}

func NewRunnerAction(writer io.Writer, errHandler actionHandler, reportKind string) *RunnerAction {
	return &RunnerAction{writer, errHandler, reportKind}
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
		planPath = "./plan.tfplan"

		err := r.downloadPlanToFile(config, planPath)

		if err != nil {
			return nil, err
		}
	}

	err := r.downloadModuleValuesToFile(config, "./tfvars.json")

	if err != nil {
		r.errHandler(config, r.reportKind, fmt.Sprintf("Could not download module values"))

		return nil, err
	}

	// re initialize
	err = r.reInit(config)

	if err != nil {
		return nil, r.errHandler(config, r.reportKind, fmt.Sprintf("Could not initialize Terraform backend: %s", err.Error()))
	}

	err = r.apply(config, planPath, "./tfvars.json")

	if err != nil {
		return nil, r.errHandler(config, r.reportKind, fmt.Sprintf("Could not apply Terraform changes: %s", err.Error()))
	}

	// get the output
	return r.output(config)
}

func (r *RunnerAction) Destroy(
	config *runner.Config,
	vals map[string]interface{},
) error {
	if !commandExists("terraform") {
		return fmt.Errorf("terraform cli command does not exist")
	}

	var planPath string

	// download plan, if github commit sha is passed in
	if config.ConfigFile.GithubSHA != "" {
		planPath = "./plan.tfplan"

		err := r.downloadPlanToFile(config, planPath)

		if err != nil {
			return err
		}
	}

	err := r.downloadModuleValuesToFile(config, "./tfvars.json")

	if err != nil {
		r.errHandler(config, r.reportKind, fmt.Sprintf("Could not download module values"))

		return err
	}

	// re initialize
	err = r.reInit(config)

	if err != nil {
		return r.errHandler(config, r.reportKind, fmt.Sprintf("Could not initialize Terraform backend: %s", err.Error()))
	}

	err = r.destroy(config, planPath, "./tfvars.json")

	if err != nil {
		return r.errHandler(config, r.reportKind, fmt.Sprintf("Could not apply Terraform changes: %s", err.Error()))
	}

	// get the output
	return nil
}

func (r *RunnerAction) Plan(
	config *runner.Config,
	vals map[string]interface{},
) ([]byte, []byte, []byte, error) {
	if !commandExists("terraform") {
		return nil, nil, nil, fmt.Errorf("terraform cli command does not exist")
	}

	err := r.downloadModuleValuesToFile(config, "./tfvars.json")

	if err != nil {
		r.errHandler(config, r.reportKind, fmt.Sprintf("Could not download module values from server"))
		return nil, nil, nil, err
	}

	// re initialize
	err = r.reInit(config)

	if err != nil {
		return nil, nil, nil, r.errHandler(config, r.reportKind, fmt.Sprintf("Failed while reinitializing the Terraform backend: %s", err.Error()))
	}

	err = r.plan(config, "./tfvars.json")

	if err != nil {
		return nil, nil, nil, r.errHandler(config, r.reportKind, fmt.Sprintf("Failed while running plan: %s", err.Error()))
	}

	zipOut, err := r.getPlanZIP(config)

	if err != nil {
		return nil, nil, nil, r.errHandler(config, r.reportKind, fmt.Sprintf("Failed while getting zip output: %s", err.Error()))
	}

	prettyOut, err := r.showPretty(config)

	if err != nil {
		return nil, nil, nil, r.errHandler(config, r.reportKind, fmt.Sprintf("Failed while generating prettified output: %s", err.Error()))
	}

	jsonOut, err := r.showJSON(config)

	if err != nil {
		return nil, nil, nil, r.errHandler(config, r.reportKind, fmt.Sprintf("Failed while generating JSON output: %s", err.Error()))
	}

	return zipOut, prettyOut, jsonOut, nil
}

type MonitorFunc func(
	r *RunnerAction,
	config *runner.Config,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error)

func MonitorState(
	r *RunnerAction,
	config *runner.Config,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(config, policyBytes, r.populateVariables, r.populatePlan, r.populateState)
}

func MonitorPlan(
	r *RunnerAction,
	config *runner.Config,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(config, policyBytes, r.populateVariables, r.populatePlan, r.populateState)
}

func MonitorBeforePlan(
	r *RunnerAction,
	config *runner.Config,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(config, policyBytes, r.populateVariables, r.populateState)
}

func MonitorAfterPlan(
	r *RunnerAction,
	config *runner.Config,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(config, policyBytes, r.populateVariables, r.populatePlan, r.populateState)
}

func MonitorBeforeApply(
	r *RunnerAction,
	config *runner.Config,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(config, policyBytes, r.populateVariables, r.populatePlan, r.populateState)
}

func MonitorAfterApply(
	r *RunnerAction,
	config *runner.Config,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(config, policyBytes, r.populateVariables, r.populatePlan, r.populateState)
}

type populatorFunc func(
	config *runner.Config,
	input map[string]interface{},
) error

func (r *RunnerAction) monitor(
	config *runner.Config,
	policyBytes []byte,
	populators ...populatorFunc,
) (*types.CreateMonitorResultRequest, error) {
	if !commandExists("terraform") {
		return nil, fmt.Errorf("terraform cli command does not exist")
	}

	err := r.reInit(config)

	if err != nil {
		return nil, r.errHandler(config, r.reportKind, fmt.Sprintf("Could not initialize Terraform backend: %s", err.Error()))
	}

	input := make(map[string]interface{})

	for _, f := range populators {
		err = f(config, input)

		if err != nil {
			return nil, err
		}
	}

	opaQuery, err := opa.LoadQueryFromBytes(opa.PACKAGE_HATCHET_MODULE, policyBytes)

	if err != nil {
		return nil, r.errHandler(config, r.reportKind, fmt.Sprintf("Could not load OPA query: %s", err.Error()))
	}

	return opa.RunMonitorQuery(opaQuery, input)
}

func (r *RunnerAction) populateState(
	config *runner.Config,
	input map[string]interface{},
) error {
	stateBytes, err := r.showStateJSON(config)

	if err != nil {
		return r.errHandler(config, r.reportKind, fmt.Sprintf("Could not get Terraform state bytes: %s", err.Error()))
	}

	state := make(map[string]interface{})

	err = json.Unmarshal(stateBytes, &state)

	if err != nil {
		return r.errHandler(config, r.reportKind, fmt.Sprintf("Could not unmarshal Terraform state to json: %s", err.Error()))
	}

	input["state"] = state

	return nil
}

func (r *RunnerAction) populatePlan(
	config *runner.Config,
	input map[string]interface{},
) error {
	// if there's a github SHA that we can retrieve the plan from, download the plan to a file
	if config.ConfigFile.GithubSHA != "" {
		planPath := "./plan.tfplan"

		err := r.downloadPlanToFile(config, planPath)

		if err != nil {
			return err
		}
	} else {
		err := r.plan(config, "./tfvars.json")

		if err != nil {
			return r.errHandler(config, r.reportKind, fmt.Sprintf("Failed while running plan for monitor: %s", err.Error()))
		}
	}

	planBytes, err := r.showJSON(config)

	if err != nil {
		return r.errHandler(config, r.reportKind, fmt.Sprintf("Failed while generating JSON plan output: %s", err.Error()))
	}

	plan := make(map[string]interface{})

	err = json.Unmarshal(planBytes, &plan)

	if err != nil {
		return r.errHandler(config, r.reportKind, fmt.Sprintf("Could not unmarshal Terraform plan to json: %s", err.Error()))
	}

	input["plan"] = plan

	return nil
}

func (r *RunnerAction) populateVariables(
	config *runner.Config,
	input map[string]interface{},
) error {
	vars, err := r.getModuleValues(config)

	if err != nil {
		r.errHandler(config, r.reportKind, fmt.Sprintf("Could not download module values from server"))
		return err
	}

	input["variables"] = vars

	return nil
}

func (r *RunnerAction) downloadPlanToFile(config *runner.Config, planPath string) error {
	resp, _, err := config.FileClient.GetPlanByCommitSHA(
		config.ConfigFile.TeamID,
		config.ConfigFile.ModuleID,
		config.ConfigFile.ModuleRunID,
	)

	if resp != nil {
		defer resp.Close()
	}

	if err != nil {
		r.errHandler(config, r.reportKind, fmt.Sprintf("Could not get plan to apply"))

		return err
	}

	fileBytes, err := ioutil.ReadAll(resp)

	if err != nil {
		r.errHandler(config, r.reportKind, "")

		return err
	}

	err = ioutil.WriteFile(filepath.Join(config.TerraformConf.TFDir, planPath), fileBytes, 0666)

	if err != nil {
		r.errHandler(config, r.reportKind, "")

		return err
	}

	return nil
}

func (r *RunnerAction) downloadModuleValuesToFile(config *runner.Config, relPath string) error {
	// download values
	vals, err := r.getModuleValues(config)

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

func (r *RunnerAction) getModuleValues(config *runner.Config) (map[string]interface{}, error) {
	vals, _, err := config.APIClient.ModulesApi.GetCurrentModuleValues(
		context.Background(),
		config.ConfigFile.TeamID,
		config.ConfigFile.ModuleID,
		&swagger.ModulesApiGetCurrentModuleValuesOpts{
			GithubSha: optional.NewString(config.ConfigFile.GithubSHA),
		},
	)

	if err != nil {
		return nil, err
	}

	return vals, nil
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
		fmt.Sprintf("TF_LOG=JSON"),
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

	// var-file option can only be set when there is not a planned run
	if valsFilePath != "" && planPath == "" {
		args = append(args, fmt.Sprintf("-var-file=%s", valsFilePath))
	}

	if planPath != "" {
		args = append(args, fmt.Sprintf("%s", planPath))
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

func (r *RunnerAction) destroy(
	config *runner.Config,
	planPath string,
	valsFilePath string,
) error {
	args := []string{"destroy", "-json", "-auto-approve"}

	// var-file option can only be set when there is not a planned run
	if valsFilePath != "" && planPath == "" {
		args = append(args, fmt.Sprintf("-var-file=%s", valsFilePath))
	}

	if planPath != "" {
		args = append(args, fmt.Sprintf("%s", planPath))
	}

	fmt.Printf("running destroy with args: [%v]", args)

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
	cmd.Stdout = r.writer
	cmd.Stderr = r.writer

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

func (r *RunnerAction) showStateJSON(
	config *runner.Config,
) ([]byte, error) {
	args := []string{"show", "-json"}

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
