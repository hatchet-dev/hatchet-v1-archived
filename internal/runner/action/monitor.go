package action

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/hatchet-dev/hatchet/api/v1/client/swagger"
	"github.com/hatchet-dev/hatchet/internal/config/runner"
)

func RunMonitorFunc(f MonitorFunc, action *RunnerAction, rc *runner.Config) (string, error) {
	policyBytes, err := DownloadMonitorPolicy(rc)

	if err != nil {
		return "failed", err
	}

	res, err := f(action, policyBytes)

	if err != nil {
		return "failed", err
	}

	res.MonitorID = rc.ConfigFile.Resources.ModuleMonitorID

	_, err = rc.APIClient.ModulesApi.CreateMonitorResult(
		context.Background(),
		swagger.CreateMonitorResultRequest{
			MonitorId:       rc.ConfigFile.Resources.ModuleMonitorID,
			FailureMessages: res.FailureMessages,
			SuccessMessage:  res.SuccessMessage,
			Severity:        res.Severity,
			Status:          res.Status,
			Title:           res.Title,
		},
		rc.ConfigFile.Resources.TeamID,
		rc.ConfigFile.Resources.ModuleID,
		rc.ConfigFile.Resources.ModuleRunID,
	)

	if err != nil {
		ErrorHandler(rc, "monitor", fmt.Sprintf("Could not report monitor result to server"))

		return "failed", err
	}

	if res.Status == "failed" {
		ErrorHandler(rc, "monitor", "")

		return "failed", err
	}

	return "succeeded", SuccessHandler(rc, "monitor", "")
}

func DownloadMonitorPolicy(config *runner.Config) ([]byte, error) {
	resp, apiErr, err := config.FileClient.GetMonitorPolicy(config.ConfigFile.Resources.TeamID, config.ConfigFile.Resources.ModuleMonitorID)

	if resp != nil {
		defer resp.Close()
	}

	if err != nil {
		return nil, err
	}

	if apiErr != nil {
		return nil, fmt.Errorf(apiErr.Description)
	}

	if resp == nil {
		return nil, fmt.Errorf("empty body from server")
	}

	return ioutil.ReadAll(resp)
}
