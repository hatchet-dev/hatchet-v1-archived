package cli

import (
	"os"

	"github.com/fatih/color"
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
		return err
	}

	writer, err := getWriter(rc)

	if err != nil {
		return err
	}

	a := action.NewRunnerAction(writer)

	_, err = a.Plan(rc, map[string]interface{}{})

	if err != nil {
		return err
	}

	return nil
}
