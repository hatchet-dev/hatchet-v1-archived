package temporal

import (
	"context"
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/logflusher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/queuechecker"
	"go.temporal.io/sdk/client"
)

const (
	BackgroundQueueName         string = "background"
	ModuleRunSchedulerQueueName string = "module_run_scheduler_queue"
	ModuleRunQueueName          string = "module_run_queue"
)

const (
	BackgroundLogFlushID string = "log_flusher"
	// TODO: module run scheduler should in background queue?
	BackgroundQueueCheckerID       string = "queue_checker"
	BackgroundModuleRunSchedulerID string = "module_run_scheduler"
)

const (
	WorkflowTypeNameLogFlush string = "FlushLogs"
	// TODO: consolidate language here
	WorkflowTypeNameCheckModuleQueue string = "ScheduleFromQueue"
	WorkflowTypeNameCheckAllQueues   string = "CheckQueues"
)

type Client struct {
	tc client.Client
}

type ClientOpts struct {
	HostPort      string
	Namespace     string
	AuthHeaderKey string
	AuthHeaderVal string
}

func NewTemporalClient(opts *ClientOpts) (*Client, error) {
	tOpts := client.Options{
		HostPort:  opts.HostPort,
		Namespace: opts.Namespace,
	}

	if opts.AuthHeaderKey != "" && opts.AuthHeaderVal != "" {
		tOpts.HeadersProvider = authHeadersProvider{
			headers: map[string]string{
				opts.AuthHeaderKey: opts.AuthHeaderVal,
			},
		}
	}

	c, err := client.Dial(tOpts)

	if err != nil {
		return nil, err
	}

	return &Client{c}, nil
}

func (c *Client) GetClient() client.Client {
	return c.tc
}

func (c *Client) Close() {
	c.Close()
}

func (c *Client) StartBackgroundTasks() error {
	logFlusherInput := logflusher.FlushLogsInput{}

	logFlusherOptions := client.StartWorkflowOptions{
		ID:           BackgroundLogFlushID,
		TaskQueue:    BackgroundQueueName,
		CronSchedule: "* * * * *",
	}

	_, err := c.tc.ExecuteWorkflow(context.Background(), logFlusherOptions, WorkflowTypeNameLogFlush, logFlusherInput)

	if err != nil {
		return err
	}

	queueCheckerInput := queuechecker.GlobalQueueInput{}

	queueCheckerOptions := client.StartWorkflowOptions{
		ID:           BackgroundQueueCheckerID,
		TaskQueue:    BackgroundQueueName,
		CronSchedule: "* * * * *",
	}

	_, err = c.tc.ExecuteWorkflow(context.Background(), queueCheckerOptions, WorkflowTypeNameCheckAllQueues, queueCheckerInput)

	return err
}

func (c *Client) TriggerModuleRunQueueChecker(input *modulequeuechecker.CheckQueueInput) error {
	options := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("%s-%s", input.TeamID, input.ModuleID),
		TaskQueue: ModuleRunSchedulerQueueName,
	}

	_, err := c.tc.ExecuteWorkflow(context.Background(), options, WorkflowTypeNameCheckModuleQueue, input)

	return err
}

type authHeadersProvider struct {
	headers map[string]string
}

func (a authHeadersProvider) GetHeaders(ctx context.Context) (map[string]string, error) {
	return a.headers, nil
}
