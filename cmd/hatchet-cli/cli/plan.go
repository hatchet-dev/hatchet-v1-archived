package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/runner/action"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "run a Hatchet plan",
	Run: func(cmd *cobra.Command, args []string) {
		preflight()

		err := runPlan()

		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "could not run plan: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	planCmd.PersistentFlags().StringVar(
		&planVarFilePath,
		"var-file",
		"",
		"The var file.",
	)
}

var planVarFilePath string

func runPlan() error {
	path, err := os.Getwd()

	if err != nil {
		return err
	}

	// get the module id based on the file path
	mods, _, err := config.APIClient.ModulesApi.ListModules(
		context.Background(),
		config.ConfigFile.TeamID,
		&swagger.ModulesApiListModulesOpts{},
	)

	if err != nil {
		return err
	}

	var matchedMod *swagger.Module

	for _, mod := range mods.Rows {
		if mod.Deployment.Path == path {
			matchedMod = &mod
			break
		}
	}

	run, a, rc, err := getAction(matchedMod.Id, "plan", path)

	if err != nil {
		return err
	}

	succeeded, err := runAllMonitors(run.Monitors, "before_plan", action.MonitorBeforePlan, a, rc)

	if !succeeded {
		return err
	}

	// pass in module variables
	variables, err := loadVarFile(planVarFilePath)

	if err != nil {
		return err
	}

	zipOut, prettyOut, jsonOut, err := a.Plan(variables)

	if err != nil {
		return err
	}

	fmt.Print(string(prettyOut))

	_, err = rc.FileClient.UploadPlanZIPFile(
		rc.ConfigFile.Resources.TeamID,
		rc.ConfigFile.Resources.ModuleID,
		rc.ConfigFile.Resources.ModuleRunID,
		zipOut,
	)

	if err != nil {
		action.ErrorHandler(rc, "core", fmt.Sprintf("Could not upload plan file to server"))

		return err
	}

	_, err = rc.APIClient.ModulesApi.CreateTerraformPlan(
		context.Background(),
		swagger.CreateTerraformPlanRequest{
			PlanJson:   string(jsonOut),
			PlanPretty: string(prettyOut),
		},
		rc.ConfigFile.Resources.TeamID,
		rc.ConfigFile.Resources.ModuleID,
		rc.ConfigFile.Resources.ModuleRunID,
	)

	if err != nil {
		action.ErrorHandler(rc, "core", fmt.Sprintf("Could not create terraform plan file on server"))

		return err
	}

	err = action.SuccessHandler(rc, "core", "")

	if err != nil {
		return err
	}

	succeeded, err = runAllMonitors(run.Monitors, "after_plan", action.MonitorAfterPlan, a, rc)

	if !succeeded {
		return err
	}

	return nil
}

func getHostName() string {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "Unknown"
	}
	return hostName
}
