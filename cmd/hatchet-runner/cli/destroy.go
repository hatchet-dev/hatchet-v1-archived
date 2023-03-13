package cli

import (
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/runner/action"

	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use: "destroy",
	Run: func(cmd *cobra.Command, args []string) {
		err := runApply()

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running destroy:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}

func runDestroy() error {
	configLoader := &loader.ConfigLoader{}
	rc, err := configLoader.LoadRunnerConfig()

	if err != nil {
		return err
	}

	err = downloadGithubRepoContents(rc)

	if err != nil {
		return err
	}

	stdoutWriter, stderrWriter, err := action.GetWriters(rc)

	if err != nil {
		return err
	}

	action := action.NewRunnerAction(&action.RunnerActionOpts{
		StdoutWriter: stdoutWriter,
		StderrWriter: stderrWriter,
		ErrHandler:   errorHandler,
		ReportKind:   "core",
		RequireInit:  true,
	})

	err = action.Destroy(map[string]interface{}{})

	if err != nil {
		return err
	}

	return successHandler(rc, "core", "")
}
