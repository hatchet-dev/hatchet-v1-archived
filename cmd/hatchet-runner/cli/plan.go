package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/runner/action"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use: "plan",
	Run: func(cmd *cobra.Command, args []string) {
		err := runPlan()

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running plan:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}

func runPlan() error {
	configLoader := &loader.ConfigLoader{}
	rc, err := configLoader.LoadRunnerConfig()

	if err != nil {
		return err
	}

	err = downloadGithubRepoContents(rc)

	if err != nil {
		return action.ErrorHandler(rc, "core", fmt.Sprintf("Could not download Github repository contents: %s", err.Error()))
	}

	stdoutWriter, stderrWriter, err := action.GetWriters(rc)

	if err != nil {
		return err
	}

	a := action.NewRunnerAction(&action.RunnerActionOpts{
		Config:       rc,
		StdoutWriter: stdoutWriter,
		StderrWriter: stderrWriter,
		ErrHandler:   action.ErrorHandler,
		ReportKind:   "core",
		RequireInit:  true,
	})

	zipOut, prettyOut, jsonOut, err := a.Plan(nil)

	if err != nil {
		return err
	}

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

	return action.SuccessHandler(rc, "core", "")
}
