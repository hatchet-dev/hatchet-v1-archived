package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
	"github.com/hatchet-dev/hatchet/internal/runner/action"

	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use: "monitor",
}

var monitorStateCmd = &cobra.Command{
	Use: "state",
	Run: func(cmd *cobra.Command, args []string) {
		err := runMonitorState()

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
		err := runMonitorPlan()

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
		err := runMonitorBeforePlan()

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [monitor before-plan]:", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
	monitorCmd.AddCommand(monitorStateCmd)
	monitorCmd.AddCommand(monitorPlanCmd)
	monitorCmd.AddCommand(monitorBeforePlanCmd)
}

func runMonitorState() error {
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

	policyBytes, err := downloadMonitorPolicy(rc)

	if err != nil {
		return err
	}

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

	if res.Status == "failed" {
		err := fmt.Errorf("Monitor failed") // TODO: better error message
		errorHandler(rc, err.Error())

		return err
	}

	return nil
}

func runMonitorPlan() error {
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

	policyBytes, err := downloadMonitorPolicy(rc)

	if err != nil {
		return err
	}

	res, err := action.MonitorPlan(rc, policyBytes)

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

	if res.Status == "failed" {
		err := fmt.Errorf("Monitor failed") // TODO: better error message
		errorHandler(rc, err.Error())

		return err
	}

	return nil
}

func runMonitorBeforePlan() error {
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

	policyBytes, err := downloadMonitorPolicy(rc)

	if err != nil {
		return err
	}

	res, err := action.MonitorBeforePlan(rc, policyBytes)

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

	if res.Status == "failed" {
		err := fmt.Errorf("Monitor failed") // TODO: better error message
		errorHandler(rc, err.Error())

		return err
	}

	return nil
}

func downloadMonitorPolicy(config *runner.Config) ([]byte, error) {
	resp, _, err := config.FileClient.GetMonitorPolicy(config.ConfigFile.TeamID, config.ConfigFile.ModuleMonitorID)

	if resp != nil {
		defer resp.Close()
	}

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp)
}
