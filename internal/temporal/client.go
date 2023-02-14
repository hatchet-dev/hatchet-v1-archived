package temporal

import (
	"context"

	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/logflusher"
	"go.temporal.io/sdk/client"
)

const (
	BackgroundQueueName string = "background"
)

const (
	BackgroundLogFlushID string = "log_flusher"
)

const (
	WorkflowTypeNameLogFlush string = "FlushLogs"
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
	input := logflusher.FlushLogsInput{}

	options := client.StartWorkflowOptions{
		ID:           BackgroundLogFlushID,
		TaskQueue:    BackgroundQueueName,
		CronSchedule: "* * * * *",
	}

	_, err := c.tc.ExecuteWorkflow(context.Background(), options, WorkflowTypeNameLogFlush, input)

	return err
}

type authHeadersProvider struct {
	headers map[string]string
}

func (a authHeadersProvider) GetHeaders(ctx context.Context) (map[string]string, error) {
	return a.headers, nil
}
