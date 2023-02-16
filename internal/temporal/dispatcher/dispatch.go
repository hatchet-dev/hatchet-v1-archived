package dispatcher

import (
	"context"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/temporal/enums"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/logflusher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/queuechecker"
	"go.temporal.io/sdk/client"
)

func DispatchModuleRunQueueChecker(c client.Client, input *modulequeuechecker.CheckQueueInput) error {
	options := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s-%s", input.TeamID, input.ModuleID),
		TaskQueue: enums.BackgroundQueueName,
	}

	// TODO: queue name shouldn't always be background
	_, err := c.ExecuteWorkflow(context.Background(), options, enums.WorkflowTypeNameCheckModuleQueue, input)

	return err
}

func DispatchBackgroundTasks(c client.Client) error {
	logFlusherInput := logflusher.FlushLogsInput{}

	logFlusherOptions := client.StartWorkflowOptions{
		ID:           enums.BackgroundLogFlushID,
		TaskQueue:    enums.BackgroundQueueName,
		CronSchedule: "* * * * *",
	}

	_, err := c.ExecuteWorkflow(context.Background(), logFlusherOptions, enums.WorkflowTypeNameLogFlush, logFlusherInput)

	if err != nil {
		return err
	}

	queueCheckerInput := queuechecker.GlobalQueueInput{}

	queueCheckerOptions := client.StartWorkflowOptions{
		ID:           enums.BackgroundQueueCheckerID,
		TaskQueue:    enums.BackgroundQueueName,
		CronSchedule: "* * * * *",
	}

	_, err = c.ExecuteWorkflow(context.Background(), queueCheckerOptions, enums.WorkflowTypeNameCheckAllQueues, queueCheckerInput)

	return err
}
