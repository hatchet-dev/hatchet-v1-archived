package cli

import (
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/runner/action"

	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use: "monitor",
}

var monitorStateCmd = &cobra.Command{
	Use: "state",
	Run: func(cmd *cobra.Command, args []string) {
		err := runMonitorFunc(action.MonitorState)

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [monitor state]:", err.Error())
			os.Exit(1)
		}
	},
}

var monitorPlanCmd = &cobra.Command{
	Use: "plan",
	Run: func(cmd *cobra.Command, args []string) {
		err := runMonitorFunc(action.MonitorPlan)

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [monitor plan]:", err.Error())
			os.Exit(1)
		}
	},
}

var monitorBeforePlanCmd = &cobra.Command{
	Use: "before-plan",
	Run: func(cmd *cobra.Command, args []string) {
		err := runMonitorFunc(action.MonitorBeforePlan)

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [monitor before-plan]:", err.Error())
			os.Exit(1)
		}
	},
}

var monitorAfterPlanCmd = &cobra.Command{
	Use: "after-plan",
	Run: func(cmd *cobra.Command, args []string) {
		err := runMonitorFunc(action.MonitorAfterPlan)

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [monitor after-plan]:", err.Error())
			os.Exit(1)
		}
	},
}

var monitorBeforeApplyCmd = &cobra.Command{
	Use: "before-apply",
	Run: func(cmd *cobra.Command, args []string) {
		err := runMonitorFunc(action.MonitorBeforeApply)

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [monitor before-apply]:", err.Error())
			os.Exit(1)
		}
	},
}

var monitorAfterApplyCmd = &cobra.Command{
	Use: "after-apply",
	Run: func(cmd *cobra.Command, args []string) {
		err := runMonitorFunc(action.MonitorAfterApply)

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [monitor before-apply]:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
	monitorCmd.AddCommand(monitorStateCmd)
	monitorCmd.AddCommand(monitorPlanCmd)
	monitorCmd.AddCommand(monitorBeforePlanCmd)
	monitorCmd.AddCommand(monitorAfterPlanCmd)
	monitorCmd.AddCommand(monitorBeforeApplyCmd)
	monitorCmd.AddCommand(monitorAfterApplyCmd)
}

func runMonitorFunc(f action.MonitorFunc) error {
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

	a := action.NewRunnerAction(&action.RunnerActionOpts{
		Config:       rc,
		StdoutWriter: stdoutWriter,
		StderrWriter: stderrWriter,
		ErrHandler:   action.ErrorHandler,
		ReportKind:   "monitor",
		RequireInit:  true,
	})

	_, err = action.RunMonitorFunc(f, a, rc)

	return err
}
