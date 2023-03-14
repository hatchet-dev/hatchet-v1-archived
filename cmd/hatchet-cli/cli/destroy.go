package cli

import (
	"context"
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/runner/action"
	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "run a Hatchet destroy",
	Run: func(cmd *cobra.Command, args []string) {
		preflight()

		err := runDestroy()

		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "could not run destroy: %v\n", err)
			os.Exit(1)
		}
	},
}

var destroyVarFilePath string

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.PersistentFlags().StringVar(
		&destroyVarFilePath,
		"var-file",
		"",
		"The var file.",
	)
}

func runDestroy() error {
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

	run, a, rc, err := getAction(matchedMod.Id, "destroy", path)

	if err != nil {
		return err
	}

	succeeded, err := runAllMonitors(run.Monitors, "before_destroy", action.MonitorBeforeDestroy, a, rc)

	if !succeeded {
		return err
	}

	// pass in module variables
	variables, err := loadVarFile(destroyVarFilePath)

	if err != nil {
		return err
	}

	err = a.Destroy(variables)

	if err != nil {
		return err
	}

	err = action.SuccessHandler(rc, "core", "")

	if err != nil {
		return err
	}

	succeeded, err = runAllMonitors(run.Monitors, "after_destroy", action.MonitorAfterDestroy, a, rc)

	if !succeeded {
		return err
	}

	return nil
}
