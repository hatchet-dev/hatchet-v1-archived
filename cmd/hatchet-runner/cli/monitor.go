package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/runner/action"

	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use: "monitor",
	Run: func(cmd *cobra.Command, args []string) {
		err := runMonitor()

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running monitor:", err.Error())
			os.Exit(1)
		}
	},
}

var policyFilePath string

func init() {
	rootCmd.AddCommand(monitorCmd)

	monitorCmd.PersistentFlags().StringVarP(
		&policyFilePath,
		"policy-file",
		"",
		"",
		"path to the OPA policy file to run the command against",
	)

	monitorCmd.MarkPersistentFlagRequired("policy-file")
}

func runMonitor() error {
	// load policy file bytes
	policyBytes, err := ioutil.ReadFile(policyFilePath)

	if err != nil {
		return err
	}

	configLoader := &loader.EnvConfigLoader{}
	rc, err := configLoader.LoadRunnerConfigFromEnv()

	if err != nil {
		return err
	}

	writer, err := getWriter(rc)

	if err != nil {
		return err
	}

	action := action.NewRunnerAction(writer, errorHandler)

	res, err := action.MonitorState(rc, policyBytes)

	if err != nil {
		return err
	}

	res.MonitorID = rc.ConfigFile.ModuleMonitorID

	_, err = rc.APIClient.ModulesApi.CreateMonitorResult(
		context.Background(),
		swagger.CreateMonitorResultRequest{
			MonitorID:       rc.ConfigFile.ModuleMonitorID,
			FailureMessages: res.FailureMessages,
			SuccessMessage:  res.SuccessMessage,
			Severity:        res.Severity,
			Status:          res.Status,
			Title:           res.Title,
		},
		rc.ConfigFile.TeamID,
		rc.ConfigFile.ModuleID,
		rc.ConfigFile.ModuleRunID,
	)

	if err != nil {
		errorHandler(rc, fmt.Sprintf("Could not report monitor result to server"))

		return err
	}

	return successHandler(rc, "")
}
