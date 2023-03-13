package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	gojson "encoding/json"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/hatchet-dev/hatchet/internal/runner/action"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"github.com/zclconf/go-cty/cty/json"
)

func getAction(modID, kind, path string) (*action.RunnerAction, *runner.Config, error) {
	// create a new module run via the API
	run, _, err := config.APIClient.ModulesApi.CreateModuleRun(
		context.Background(),
		swagger.CreateModuleRunRequest{
			Kind:     kind,
			Hostname: getHostName(),
		},
		config.ConfigFile.TeamID,
		modID,
	)

	if err != nil {
		return nil, nil, err
	}

	// get a module run token via the API
	tok, _, err := config.APIClient.ModulesApi.GetModuleRunLocalToken(
		context.Background(),
		config.ConfigFile.TeamID,
		modID,
		run.Id,
	)

	if err != nil {
		return nil, nil, err
	}

	// using this token, we generate a runner config
	sharedCF := &shared.ConfigFile{
		Debug: true,
	}

	runnerCF := &runner.ConfigFile{
		Resources: runner.ConfigFileResources{
			TeamID:      config.ConfigFile.TeamID,
			ModuleID:    modID,
			ModuleRunID: run.Id,
		},
		GRPC: runner.ConfigFileGRPC{
			GRPCServerAddress: config.ConfigFile.Address,
			GRPCToken:         tok.Token,
		},
		API: runner.ConfigFileAPI{
			APIToken:         tok.Token,
			APIServerAddress: config.ConfigFile.Address,
		},
		Terraform: runner.ConfigFileTerraform{
			TFDir: filepath.Dir(path),
		},
	}

	sc, err := loader.GetSharedConfigFromConfigFile(sharedCF)

	if err != nil {
		return nil, nil, err
	}

	rc, err := loader.GetRunnerConfigFromConfigFile(runnerCF, sc)

	if err != nil {
		errorHandler(rc, "core", fmt.Sprintf("Could not get runner configuration: %s", err.Error()))

		return nil, nil, err
	}

	stdoutWriter, stderrWriter, err := action.GetWriters(rc)

	if err != nil {
		errorHandler(rc, "core", fmt.Sprintf("Could not get writers for plan"))

		return nil, nil, err
	}

	a := action.NewRunnerAction(&action.RunnerActionOpts{
		Config:       rc,
		StdoutWriter: stdoutWriter,
		StderrWriter: stderrWriter,
		ErrHandler:   errorHandler,
		ReportKind:   "core",
		RequireInit:  false,
	})

	return a, rc, nil
}

func errorHandler(config *runner.Config, reportKind, description string) error {
	_, _, err := config.APIClient.ModulesApi.FinalizeModuleRun(
		context.Background(),
		swagger.FinalizeModuleRunRequest{
			Status:      string(types.ModuleRunStatusFailed),
			Description: description,
			ReportKind:  reportKind,
		},
		config.ConfigFile.Resources.TeamID,
		config.ConfigFile.Resources.ModuleID,
		config.ConfigFile.Resources.ModuleRunID,
	)

	if err != nil {
		return fmt.Errorf("Error reporting error: %s. Original error: %s", err.Error(), description)
	}

	return fmt.Errorf(description)
}

func successHandler(config *runner.Config, reportKind, description string) error {
	_, _, err := config.APIClient.ModulesApi.FinalizeModuleRun(
		context.Background(),
		swagger.FinalizeModuleRunRequest{
			Status:      string(types.ModuleRunStatusCompleted),
			Description: description,
			ReportKind:  reportKind,
		},
		config.ConfigFile.Resources.TeamID,
		config.ConfigFile.Resources.ModuleID,
		config.ConfigFile.Resources.ModuleRunID,
	)

	if err != nil {
		return fmt.Errorf("Error reporting success: %s. Success message: %s", err.Error(), description)
	}

	return nil
}

func preflight() {
	if config.ConfigFile.OrganizationID == "" {
		color.New(color.FgRed).Fprintf(os.Stderr, "team id must be set: run [hatchet config set-organization]\n")
		os.Exit(1)
	}

	if config.ConfigFile.TeamID == "" {
		color.New(color.FgRed).Fprintf(os.Stderr, "team id must be set: run [hatchet config set-team]\n")
		os.Exit(1)
	}

	if config.ConfigFile.APIToken == "" {
		color.New(color.FgRed).Fprintf(os.Stderr, "api token must be set: run [hatchet config set-api-token]\n")
		os.Exit(1)
	}
}

func loadVarFile(varFilePath string) (map[string]interface{}, error) {
	if varFilePath == "" {
		// look for terraform.tfvars or terraform.tfvars.json
		if fileExists("terraform.tfvars") {
			varFilePath = "terraform.tfvars"
		} else if fileExists("terraform.tfvars.json") {
			varFilePath = "terraform.tfvars.json"
		} else {
			return map[string]interface{}{}, nil
		}
	}

	p := hclparse.NewParser()

	src, err := ioutil.ReadFile(varFilePath)

	if err != nil {
		return nil, fmt.Errorf("could not read file at path %s: %w", varFilePath, err)
	}

	var file *hcl.File
	var diags hcl.Diagnostics
	switch {
	case strings.HasSuffix(varFilePath, ".json"):
		file, diags = p.ParseJSON(src, varFilePath)
	default:
		file, diags = p.ParseHCL(src, varFilePath)
	}

	if file == nil {
		return nil, diags
	}

	ctyTarget := make(map[string]cty.Value)

	attrs, diags := file.Body.JustAttributes()

	for name, attr := range attrs {
		val, valDiags := attr.Expr.Value(nil)
		diags = append(diags, valDiags...)
		ctyTarget[name] = val
	}

	if diags.HasErrors() {
		return nil, diags
	}

	jsonTarget, err := gocty.ToCtyValue(ctyTarget, cty.Map(cty.DynamicPseudoType))

	sjv := json.SimpleJSONValue{
		Value: jsonTarget,
	}

	jsonBytes, err := sjv.MarshalJSON()

	fmt.Println("JSON BYTES", string(jsonBytes), err)

	target := make(map[string]interface{})

	err = gojson.Unmarshal(jsonBytes, &target)

	if err != nil {
		return nil, err
	}

	return target, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
