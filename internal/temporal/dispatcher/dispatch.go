package dispatcher

import (
	"context"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/temporal"
	"github.com/hatchet-dev/hatchet/internal/temporal/enums"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/logflusher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/queuechecker"
	"go.temporal.io/sdk/client"
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

	queueCheckerInput := queuechecker.GlobalQueueInput{}

	queueCheckerOptions := client.StartWorkflowOptions{
		ID:           enums.BackgroundQueueCheckerID,
		TaskQueue:    enums.BackgroundQueueName,
		CronSchedule: "* * * * *",
	}

	_, err = tc.ExecuteWorkflow(context.Background(), queueCheckerOptions, enums.WorkflowTypeNameCheckAllQueues, queueCheckerInput)

	return err
}
