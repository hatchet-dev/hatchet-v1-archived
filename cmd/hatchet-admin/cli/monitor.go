package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"

	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Command used to manage default organization/team monitors",
}

var monitorCreateCmd = &cobra.Command{
	Use:   "create-default",
	Short: "Creates a new default monitor for a given team",
	Run: func(cmd *cobra.Command, args []string) {
		err := runCreateDefaultMonitor()

		if err != nil {
			red := color.New(color.FgRed)
			red.Println("Error running [monitor create-default]:", err.Error())
			os.Exit(1)
		}
	},
}

type monitorCreateRequest struct {
	TeamID       string `form:"required,uuid"`
	Kind         string `form:"required,oneof=plan state before_plan after_plan before_apply after_apply before_destroy after_destroy"`
	CronSchedule string `form:"required"`
	PolicyBytes  []byte `form:"required"`
	Name         string `form:"required"`
}

var createRequest = monitorCreateRequest{}

var monitorCreatePolicyFilePath string

func init() {
	rootCmd.AddCommand(monitorCmd)
	monitorCmd.AddCommand(monitorCreateCmd)

	monitorCreateCmd.PersistentFlags().StringVar(
		&createRequest.CronSchedule,
		"schedule",
		"",
		"The cron schedule for the monitor.",
	)

	monitorCreateCmd.MarkPersistentFlagRequired("schedule")

	monitorCreateCmd.PersistentFlags().StringVar(
		&monitorCreatePolicyFilePath,
		"policy-file",
		"",
		"The path to the policy file.",
	)

	monitorCreateCmd.MarkPersistentFlagRequired("policy-file")

	monitorCreateCmd.PersistentFlags().StringVar(
		&createRequest.TeamID,
		"team",
		"",
		"The team ID to assign this monitor to.",
	)

	monitorCreateCmd.MarkPersistentFlagRequired("team")

	monitorCreateCmd.PersistentFlags().StringVar(
		&createRequest.Kind,
		"kind",
		"",
		"The kind of monitor: options are [plan state before_plan after_plan before_apply after_apply before_destroy after_destroy].",
	)

	monitorCreateCmd.MarkPersistentFlagRequired("kind")

	monitorCreateCmd.PersistentFlags().StringVar(
		&createRequest.Name,
		"name",
		"",
		"The name for the monitor.",
	)

	monitorCreateCmd.MarkPersistentFlagRequired("name")
}

func runCreateDefaultMonitor() error {
	// read the policy bytes
	policyBytes, err := ioutil.ReadFile(monitorCreatePolicyFilePath)

	if err != nil {
		return fmt.Errorf("could not read policy file: %s", err.Error())
	}

	createRequest.PolicyBytes = policyBytes

	// validate the monitor
	v := handlerutils.NewDefaultValidator()

	err = v.Validate(createRequest)

	if err != nil {
		return err
	}

	monitor := &models.ModuleMonitor{
		TeamID:       createRequest.TeamID,
		Kind:         models.ModuleMonitorKind(createRequest.Kind),
		DisplayName:  createRequest.Name,
		CronSchedule: createRequest.CronSchedule,
		IsDefault:    true,
		CurrentMonitorPolicyBytesVersion: models.MonitorPolicyBytesVersion{
			Version:     1,
			PolicyBytes: policyBytes,
		},
	}

	monitor, err = sc.DB.Repository.ModuleMonitor().CreateModuleMonitor(monitor)

	if err != nil {
		return err
	}

	err = dispatcher.DispatchCronMonitor(sc.TemporalClient, createRequest.TeamID, monitor.ID, createRequest.CronSchedule)

	if err != nil {
		return err
	}

	return err
}
