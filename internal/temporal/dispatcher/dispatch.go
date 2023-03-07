package dispatcher

import (
	"context"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/temporal"
	"github.com/hatchet-dev/hatchet/internal/temporal/enums"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/logflusher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/monitordispatcher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/notifier"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/queuechecker"
	"go.temporal.io/sdk/client"

	hatchetenums "github.com/hatchet-dev/hatchet/internal/temporal/enums"
)

func DispatchModuleRunQueueChecker(c *temporal.Client, input *modulequeuechecker.CheckQueueInput) error {
	tc, err := c.GetClient(enums.BackgroundQueueName)

	if err != nil {
		return err
	}

	options := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s-%s", input.TeamID, input.ModuleID),
		TaskQueue: enums.BackgroundQueueName,
	}

	_, err = tc.ExecuteWorkflow(context.Background(), options, enums.WorkflowTypeNameCheckModuleQueue, input)

	return err
}

func DispatchBackgroundTasks(c *temporal.Client) error {
	tc, err := c.GetClient(enums.BackgroundQueueName)

	if err != nil {
		return err
	}

	logFlusherInput := logflusher.FlushLogsInput{}

	logFlusherOptions := client.StartWorkflowOptions{
		ID:           enums.BackgroundLogFlushID,
		TaskQueue:    enums.BackgroundQueueName,
		CronSchedule: "* * * * *",
	}

	_, err = tc.ExecuteWorkflow(context.Background(), logFlusherOptions, enums.WorkflowTypeNameLogFlush, logFlusherInput)

	if err != nil {
		return err
	}

	notifierInput := notifier.NotifierInput{}

	notifierOptions := client.StartWorkflowOptions{
		ID:           enums.BackgroundNotifierID,
		TaskQueue:    enums.BackgroundQueueName,
		CronSchedule: "* * * * *",
	}

	_, err = tc.ExecuteWorkflow(context.Background(), notifierOptions, enums.WorkflowTypeNameNotifier, notifierInput)

	if err != nil {
		return err
	}

	queueCheckerInput := queuechecker.GlobalQueueInput{}

	queueCheckerOptions := client.StartWorkflowOptions{
		ID:           enums.BackgroundQueueCheckerID,
		TaskQueue:    enums.BackgroundQueueName,
		CronSchedule: "* * * * *",
	}

	_, err = tc.ExecuteWorkflow(context.Background(), queueCheckerOptions, enums.WorkflowTypeNameCheckAllQueues, queueCheckerInput)

	return err
}

type CronMonitorOpts struct {
	Monitor   *models.ModuleMonitor
	Team      *models.Team
	Module    *models.Module
	DB        *database.Config
	TokenOpts token.TokenOpts
	ServerURL string
}

func DispatchCronMonitor(c *temporal.Client, teamID, monitorID, cronSchedule string) error {
	// TODO: this queue name should align with the team id
	tc, err := c.GetClient(enums.BackgroundQueueName)

	if err != nil {
		return err
	}

	monitorInput := monitordispatcher.MonitorDispatcherInput{
		TeamID:    teamID,
		MonitorID: monitorID,
	}

	runMonitorOptions := client.StartWorkflowOptions{
		ID:           fmt.Sprintf("%s/%s", teamID, monitorID),
		TaskQueue:    enums.BackgroundQueueName,
		CronSchedule: cronSchedule,
	}

	_, err = tc.ExecuteWorkflow(context.Background(), runMonitorOptions, hatchetenums.WorkflowTypeNameDispatchMonitors, monitorInput)

	return err
}

func UpdateCronMonitor(c *temporal.Client, teamID, monitorID, cronSchedule string) error {
	// TODO: this queue name should align with the team id
	tc, err := c.GetClient(enums.BackgroundQueueName)

	if err != nil {
		return err
	}

	monitorInput := monitordispatcher.MonitorDispatcherInput{
		TeamID:    teamID,
		MonitorID: monitorID,
	}

	workflowID := fmt.Sprintf("%s-%s", teamID, monitorID)

	runMonitorOptions := client.StartWorkflowOptions{
		ID:           workflowID,
		TaskQueue:    enums.BackgroundQueueName,
		CronSchedule: cronSchedule,
	}

	// delete/terminate the first workflow
	tc.TerminateWorkflow(context.Background(), workflowID, "", "Terminated due to cron schedule update")

	_, err = tc.ExecuteWorkflow(context.Background(), runMonitorOptions, hatchetenums.WorkflowTypeNameDispatchMonitors, monitorInput)

	return err
}

func DeleteCronMonitor(c *temporal.Client, teamID, monitorID string) error {
	// TODO: this queue name should align with the team id
	tc, err := c.GetClient(enums.BackgroundQueueName)

	if err != nil {
		return err
	}

	workflowID := fmt.Sprintf("%s-%s", teamID, monitorID)

	return tc.TerminateWorkflow(context.Background(), workflowID, "", "Terminated due to monitor deletion")
}
