package action

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/antihax/optional"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
	"github.com/hatchet-dev/hatchet/internal/opa"
)

const hatchetVarFile = ".hatchet.tfvars.json"

type actionHandler func(config *runner.Config, reportKind, description string) error

type RunnerAction struct {
	config       *runner.Config
	stdoutWriter io.Writer
	stderrWriter io.Writer
	errHandler   actionHandler
	reportKind   string
	requireInit  bool
}

type RunnerActionOpts struct {
	Config       *runner.Config
	StdoutWriter io.Writer
	StderrWriter io.Writer
	ErrHandler   actionHandler
	ReportKind   string

	// Whether a `terraform init` should be run before plan/apply operations
	RequireInit bool
}

func NewRunnerAction(opts *RunnerActionOpts) *RunnerAction {
	return &RunnerAction{
		config:       opts.Config,
		stdoutWriter: opts.StdoutWriter,
		stderrWriter: opts.StderrWriter,
		errHandler:   opts.ErrHandler,
		reportKind:   opts.ReportKind,
		requireInit:  opts.RequireInit,
	}
}

func (r *RunnerAction) Apply(
	variables map[string]interface{},
) ([]byte, error) {
	if !commandExists("terraform") {
		return nil, fmt.Errorf("terraform cli command does not exist")
	}

	var planPath string

	// download plan, if github commit sha is passed in
	if r.config.ConfigFile.Github.GithubSHA != "" {
		planPath = "./plan.tfplan"

		err := r.downloadPlanToFile(planPath)

		if err != nil {
			return nil, err
		}
	}

	if variables == nil {
		err := r.downloadModuleValuesToFile(hatchetVarFile)

		if err != nil {
			r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not download module values"))

			return nil, err
		}
	} else {
		err := r.copyVarsToFile(variables, hatchetVarFile)

		if err != nil {
			r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not copy module variables: %s", err.Error()))
			return nil, err
		}
	}

	// re initialize
	if r.requireInit {
		err := r.reInit()

		if err != nil {
			return nil, r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not initialize Terraform backend: %s", err.Error()))
		}
	}

	err := r.apply(planPath, hatchetVarFile)

	if err != nil {
		return nil, r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not apply Terraform changes: %s", err.Error()))
	}

	// get the output
	return r.output()
}

func (r *RunnerAction) Destroy(
	variables map[string]interface{},
) error {
	if !commandExists("terraform") {
		return fmt.Errorf("terraform cli command does not exist")
	}

	var planPath string

	// download plan, if github commit sha is passed in
	if r.config.ConfigFile.Github.GithubSHA != "" {
		planPath = "./plan.tfplan"

		err := r.downloadPlanToFile(planPath)

		if err != nil {
			return err
		}
	}

	if variables == nil {
		err := r.downloadModuleValuesToFile(hatchetVarFile)

		if err != nil {
			r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not download module values"))

			return err
		}
	} else {
		err := r.copyVarsToFile(variables, hatchetVarFile)

		if err != nil {
			r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not copy module variables: %s", err.Error()))
			return err
		}
	}

	// re initialize if required
	if r.requireInit {
		err := r.reInit()

		if err != nil {
			return r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not initialize Terraform backend: %s", err.Error()))
		}
	}

	err := r.destroy(planPath, hatchetVarFile)

	if err != nil {
		return r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not apply Terraform changes: %s", err.Error()))
	}

	// get the output
	return nil
}

func (r *RunnerAction) Plan(variables map[string]interface{}) ([]byte, []byte, []byte, error) {
	if !commandExists("terraform") {
		return nil, nil, nil, fmt.Errorf("terraform cli command does not exist")
	}

	if variables == nil {
		err := r.downloadModuleValuesToFile(hatchetVarFile)

		if err != nil {
			r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not download module values from server"))
			return nil, nil, nil, err
		}
	} else {
		err := r.copyVarsToFile(variables, hatchetVarFile)

		if err != nil {
			r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not copy module variables: %s", err.Error()))
			return nil, nil, nil, err
		}
	}

	// re initialize if required
	if r.requireInit {
		err := r.reInit()

		if err != nil {
			return nil, nil, nil, r.errHandler(r.config, r.reportKind, fmt.Sprintf("Failed while reinitializing the Terraform backend: %s", err.Error()))
		}
	}

	err := r.plan(hatchetVarFile)

	if err != nil {
		return nil, nil, nil, r.errHandler(r.config, r.reportKind, fmt.Sprintf("Failed while running plan: %s", err.Error()))
	}

	zipOut, err := r.getPlanZIP()

	if err != nil {
		return nil, nil, nil, r.errHandler(r.config, r.reportKind, fmt.Sprintf("Failed while getting zip output: %s", err.Error()))
	}

	prettyOut, err := r.showPretty()

	if err != nil {
		return nil, nil, nil, r.errHandler(r.config, r.reportKind, fmt.Sprintf("Failed while generating prettified output: %s", err.Error()))
	}

	jsonOut, err := r.showJSON()

	if err != nil {
		return nil, nil, nil, r.errHandler(r.config, r.reportKind, fmt.Sprintf("Failed while generating JSON output: %s", err.Error()))
	}

	return zipOut, prettyOut, jsonOut, nil
}

func (r *RunnerAction) Init() error {
	if !commandExists("terraform") {
		return fmt.Errorf("terraform cli command does not exist")
	}

	err := r.reInit()

	if err != nil {
		return fmt.Errorf("Failed while reinitializing the Terraform backend: %s", err.Error())
	}

	return nil
}

type MonitorFunc func(
	r *RunnerAction,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error)

func MonitorState(
	r *RunnerAction,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(policyBytes, r.populateVariables, r.populatePlan, r.populateState)
}

func MonitorPlan(
	r *RunnerAction,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(policyBytes, r.populateVariables, r.populatePlan, r.populateState)
}

func MonitorBeforePlan(
	r *RunnerAction,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(policyBytes, r.populateVariables, r.populateState)
}

func MonitorAfterPlan(
	r *RunnerAction,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(policyBytes, r.populateVariables, r.populatePlan, r.populateState)
}

func MonitorBeforeApply(
	r *RunnerAction,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(policyBytes, r.populateVariables, r.populatePlan, r.populateState)
}

func MonitorAfterApply(
	r *RunnerAction,
	policyBytes []byte,
) (*types.CreateMonitorResultRequest, error) {
	return r.monitor(policyBytes, r.populateVariables, r.populatePlan, r.populateState)
}

type populatorFunc func(
	input map[string]interface{},
) error

func (r *RunnerAction) monitor(
	policyBytes []byte,
	populators ...populatorFunc,
) (*types.CreateMonitorResultRequest, error) {
	if !commandExists("terraform") {
		return nil, fmt.Errorf("terraform cli command does not exist")
	}

	if r.requireInit {
		err := r.reInit()

		if err != nil {
			return nil, r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not initialize Terraform backend: %s", err.Error()))
		}
	}

	input := make(map[string]interface{})

	for _, f := range populators {
		err := f(input)

		if err != nil {
			return nil, err
		}
	}

	opaQuery, err := opa.LoadQueryFromBytes(opa.PACKAGE_HATCHET_MODULE, policyBytes)

	if err != nil {
		return nil, r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not load OPA query: %s", err.Error()))
	}

	return opa.RunMonitorQuery(opaQuery, input)
}

func (r *RunnerAction) populateState(
	input map[string]interface{},
) error {
	stateBytes, err := r.showStateJSON()

	if err != nil {
		return r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not get Terraform state bytes: %s", err.Error()))
	}

	state := make(map[string]interface{})

	err = json.Unmarshal(stateBytes, &state)

	if err != nil {
		return r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not unmarshal Terraform state to json: %s", err.Error()))
	}

	input["state"] = state

	return nil
}

func (r *RunnerAction) populatePlan(
	input map[string]interface{},
) error {
	// if there's a github SHA that we can retrieve the plan from, download the plan to a file
	if r.config.ConfigFile.Github.GithubSHA != "" {
		planPath := "./plan.tfplan"

		err := r.downloadPlanToFile(planPath)

		if err != nil {
			return err
		}
	} else {
		err := r.plan(hatchetVarFile)

		if err != nil {
			return r.errHandler(r.config, r.reportKind, fmt.Sprintf("Failed while running plan for monitor: %s", err.Error()))
		}
	}

	planBytes, err := r.showJSON()

	if err != nil {
		return r.errHandler(r.config, r.reportKind, fmt.Sprintf("Failed while generating JSON plan output: %s", err.Error()))
	}

	plan := make(map[string]interface{})

	err = json.Unmarshal(planBytes, &plan)

	if err != nil {
		return r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not unmarshal Terraform plan to json: %s", err.Error()))
	}

	input["plan"] = plan

	return nil
}

func (r *RunnerAction) populateVariables(
	input map[string]interface{},
) error {
	vars, err := r.downloadModuleValues()

	if err != nil {
		r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not download module values from server"))
		return err
	}

	input["variables"] = vars

	return nil
}

func (r *RunnerAction) downloadPlanToFile(planPath string) error {
	resp, _, err := r.config.FileClient.GetPlanByCommitSHA(
		r.config.ConfigFile.Resources.TeamID,
		r.config.ConfigFile.Resources.ModuleID,
		r.config.ConfigFile.Resources.ModuleRunID,
	)

	if resp != nil {
		defer resp.Close()
	}

	if err != nil {
		r.errHandler(r.config, r.reportKind, fmt.Sprintf("Could not get plan to apply"))

		return err
	}

	fileBytes, err := ioutil.ReadAll(resp)

	if err != nil {
		r.errHandler(r.config, r.reportKind, "")

		return err
	}

	err = ioutil.WriteFile(filepath.Join(r.config.TerraformConf.TFDir, planPath), fileBytes, 0666)

	if err != nil {
		r.errHandler(r.config, r.reportKind, "")

		return err
	}

	return nil
}

func (r *RunnerAction) downloadModuleValuesToFile(relPath string) error {
	// download values
	vars, err := r.downloadModuleValues()

	if err != nil {
		return err
	}

	return r.copyVarsToFile(vars, relPath)
}

func (r *RunnerAction) copyVarsToFile(vars map[string]interface{}, targetPath string) error {
	fileBytes, err := json.Marshal(vars)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(r.config.TerraformConf.TFDir, targetPath), fileBytes, 0666)

	return err
}

func (r *RunnerAction) downloadModuleValues() (map[string]interface{}, error) {
	if r.config.ConfigFile != nil && r.config.ConfigFile.Github.GithubSHA != "" {
		vals, _, err := r.config.APIClient.ModulesApi.GetCurrentModuleValues(
			context.Background(),
			r.config.ConfigFile.Resources.TeamID,
			r.config.ConfigFile.Resources.ModuleID,
			&swagger.ModulesApiGetCurrentModuleValuesOpts{
				GithubSha: optional.NewString(r.config.ConfigFile.Github.GithubSHA),
			},
		)

		if err != nil {
			return nil, err
		}

		return vals, nil
	}

	return map[string]interface{}{}, nil
}

func (r *RunnerAction) reInit() error {
	cmd := exec.Command("terraform", "init", "-reconfigure")
	cmd.Dir = r.config.TerraformConf.TFDir
	cmd.Stdout = r.stdoutWriter
	cmd.Stderr = r.stderrWriter
	cmd.Stdin = strings.NewReader("yes\n")

	err := r.setBackendEnv(cmd)

	if err != nil {
		return err
	}

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func (r *RunnerAction) setBackendEnv(cmd *exec.Cmd) error {
	rc := r.config.ConfigFile

	tfStateAddress := fmt.Sprintf("%s/api/v1/teams/%s/modules/%s/runs/%s/tfstate",
		rc.API.APIServerAddress,
		rc.Resources.TeamID,
		rc.Resources.ModuleID,
		rc.Resources.ModuleRunID)

	cmd.Env = append(cmd.Environ(), []string{
		fmt.Sprintf("TF_LOG=JSON"),
		fmt.Sprintf("TF_HTTP_USERNAME=mrt"),
		fmt.Sprintf("TF_HTTP_PASSWORD=%s", rc.API.APIToken),
		fmt.Sprintf("TF_HTTP_ADDRESS=%s", tfStateAddress),
		fmt.Sprintf("TF_HTTP_LOCK_ADDRESS=%s", tfStateAddress),
		fmt.Sprintf("TF_HTTP_UNLOCK_ADDRESS=%s", tfStateAddress),
	}...)

	return nil
}

func (r *RunnerAction) apply(
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
	cmd.Dir = r.config.TerraformConf.TFDir

	cmd.Stdout = r.stdoutWriter
	cmd.Stderr = r.stderrWriter

	err := r.setBackendEnv(cmd)

	if err != nil {
		return err
	}

	return cmd.Run()
}

func (r *RunnerAction) destroy(
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
	cmd.Dir = r.config.TerraformConf.TFDir

	cmd.Stdout = r.stdoutWriter
	cmd.Stderr = r.stderrWriter

	err := r.setBackendEnv(cmd)

	if err != nil {
		return err
	}

	return cmd.Run()
}

func (r *RunnerAction) plan(
	valsFilePath string,
) error {
	args := []string{"plan", "-out=./plan.tfplan"}

	if valsFilePath != "" {
		args = append(args, fmt.Sprintf("-var-file=%s", valsFilePath))
	}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = r.config.TerraformConf.TFDir
	cmd.Stdout = r.stdoutWriter
	cmd.Stderr = r.stderrWriter

	err := r.setBackendEnv(cmd)

	if err != nil {
		return err
	}

	return cmd.Run()
}

func (r *RunnerAction) getPlanZIP() ([]byte, error) {
	path := filepath.Join(r.config.TerraformConf.TFDir, "./plan.tfplan")
	return ioutil.ReadFile(path)
}

func (r *RunnerAction) showPretty() ([]byte, error) {
	args := []string{"show", "-no-color", "./plan.tfplan"}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = r.config.TerraformConf.TFDir
	cmd.Stderr = r.stderrWriter

	err := r.setBackendEnv(cmd)

	if err != nil {
		return nil, err
	}

	return cmd.Output()
}

func (r *RunnerAction) showStateJSON() ([]byte, error) {
	args := []string{"show", "-json"}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = r.config.TerraformConf.TFDir
	cmd.Stderr = r.stderrWriter

	err := r.setBackendEnv(cmd)

	if err != nil {
		return nil, err
	}

	return cmd.Output()
}

func (r *RunnerAction) showJSON() ([]byte, error) {
	args := []string{"show", "-json", "./plan.tfplan"}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = r.config.TerraformConf.TFDir
	cmd.Stderr = r.stderrWriter

	err := r.setBackendEnv(cmd)

	if err != nil {
		return nil, err
	}

	return cmd.Output()
}

func (r *RunnerAction) output() ([]byte, error) {
	cmd := exec.Command("terraform", "output", "-json")
	cmd.Dir = r.config.TerraformConf.TFDir
	cmd.Stderr = r.stderrWriter

	err := r.setBackendEnv(cmd)

	if err != nil {
		return nil, err
	}

	return cmd.Output()
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
