package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "run a Hatchet apply",
	Run: func(cmd *cobra.Command, args []string) {
		preflight()

		err := runApply()

		if err != nil {
			color.New(color.FgRed).Fprintf(os.Stderr, "could not run apply: %v\n", err)
			os.Exit(1)
		}
	},
}

var applyVarFilePath string

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.PersistentFlags().StringVar(
		&applyVarFilePath,
		"var-file",
		"",
		"The var file.",
	)
}

func runApply() error {
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

	a, rc, err := getAction(matchedMod.Id, "apply", path)

	if err != nil {
		return err
	}

	// pass in module variables
	variables, err := loadVarFile(applyVarFilePath)

	if err != nil {
		return err
	}

	out, err := a.Apply(variables)

	if err != nil {
		return err
	}

	fmt.Print(string(out))

	return successHandler(rc, "core", "")
}
