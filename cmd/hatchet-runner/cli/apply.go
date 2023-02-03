package cli

import (
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/runner/action"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use: "apply",
	Run: func(cmd *cobra.Command, args []string) {
		err := runApply()

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running apply:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func runApply() error {
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

	action := action.NewRunnerAction(writer, errorHandler)

	_, err = action.Apply(rc, map[string]interface{}{})

	if err != nil {
		return err
	}

	return successHandler(rc, "apply ran successfully")
}
