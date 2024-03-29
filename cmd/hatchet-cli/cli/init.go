package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/runner/action"
	"github.com/spf13/cobra"
)

// initCmd initializes a module from the current directory
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes a module from the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		preflight()

		err := runInit()

		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "could not run init: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit() error {
	path, err := os.Getwd()

	if err != nil {
		return err
	}

	base := filepath.Base(path)

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

	variables, err := loadVarFile(applyVarFilePath)

	if err != nil {
		return err
	}

	var modID string

	if matchedMod == nil {
		mod, _, err := config.APIClient.ModulesApi.CreateModule(
			context.Background(),
			swagger.CreateModuleRequest{
				Name: base,
				Local: &swagger.CreateModuleRequestLocal{
					LocalPath: path,
				},
				ValuesRaw: variables,
			},
			config.ConfigFile.TeamID,
		)

		if err != nil {
			return fmt.Errorf("could not create module: %w", err)
		}

		color.New(color.FgGreen).Fprintf(os.Stdout, "successfully created module %s with id %s\n", mod.Name, mod.Id)

		modID = mod.Id
	} else {
		color.New(color.FgGreen).Fprintf(os.Stdout, "found existing module %s with id %s\n", matchedMod.Name, matchedMod.Id)

		modID = matchedMod.Id
	}

	_, a, rc, err := getAction(modID, "init", path)

	if err != nil {
		return err
	}

	err = a.Init()

	if err != nil {
		return err
	}

	return action.SuccessHandler(rc, "core", "")
}
