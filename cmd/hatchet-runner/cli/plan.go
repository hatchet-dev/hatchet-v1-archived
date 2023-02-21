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
	configLoader := &loader.EnvConfigLoader{}
	rc, err := configLoader.LoadRunnerConfigFromEnv()

	if err != nil {
		return err
	}

	err = downloadGithubRepoContents(rc)

	if err != nil {
		return errorHandler(rc, fmt.Sprintf("Could not download Github repository contents: %s", err.Error()))
	}

	writer, err := getWriter(rc)

	if err != nil {
		errorHandler(rc, fmt.Sprintf("Could not upload plan file to server"))

		return err
	}

	a := action.NewRunnerAction(writer, errorHandler)

	zipOut, prettyOut, jsonOut, err := a.Plan(rc, map[string]interface{}{})

	if err != nil {
		return err
	}

	_, err = rc.FileClient.UploadPlanZIPFile(
		rc.ConfigFile.TeamID,
		rc.ConfigFile.ModuleID,
		rc.ConfigFile.ModuleRunID,
		zipOut,
	)

	if err != nil {
		errorHandler(rc, fmt.Sprintf("Could not upload plan file to server"))

		return err
	}

	_, err = rc.APIClient.ModulesApi.CreateTerraformPlan(
		context.Background(),
		swagger.CreateTerraformPlanRequest{
			PlanJson:   string(jsonOut),
			PlanPretty: string(prettyOut),
		},
		rc.ConfigFile.TeamID,
		rc.ConfigFile.ModuleID,
		rc.ConfigFile.ModuleRunID,
	)

	if err != nil {
		errorHandler(rc, fmt.Sprintf("Could not create terraform plan file on server"))

		return err
	}

	return successHandler(rc, "")
}