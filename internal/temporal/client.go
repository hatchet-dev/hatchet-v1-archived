package temporal

import (
	"context"
	"fmt"
	"os"

	"go.temporal.io/sdk/client"
)

const DefaultQueueName = "default"

type Client struct {
	clients map[string]client.Client

	opts *ClientOpts
}

type ClientOpts struct {
	HostPort         string
	Namespace        string
	AuthHeaderKey    string
	AuthHeaderVal    string
	DefaultQueueName string
}

func NewTemporalClient(opts *ClientOpts) (*Client, error) {
	if opts.DefaultQueueName == "" {
		opts.DefaultQueueName = DefaultQueueName
	}

	c, err := clientFromOpts(opts, opts.DefaultQueueName)

	if err != nil {
		return nil, err
	}

	clients := make(map[string]client.Client)

	clients[opts.DefaultQueueName] = c

	return &Client{clients, opts}, nil
}

func clientFromOpts(opts *ClientOpts, taskQueueName string) (client.Client, error) {
	tOpts := client.Options{
		HostPort:  opts.HostPort,
		Namespace: opts.Namespace,
		Identity:  fmt.Sprintf("%d@%s@%s", os.Getpid(), getHostName(), taskQueueName),
	}

	if opts.AuthHeaderKey != "" && opts.AuthHeaderVal != "" {
		tOpts.HeadersProvider = authHeadersProvider{
			headers: map[string]string{
				opts.AuthHeaderKey: opts.AuthHeaderVal,
			},
		}
	}

	return client.Dial(tOpts)
}

func (c *Client) GetClient(queueName string) (client.Client, error) {
	if queueName == "" {
		return c.clients[DefaultQueueName], nil
	}

	tc, exists := c.clients[queueName]

	if !exists {
		return c.newQueueClient(queueName)
	}

	return tc, nil
}

func (c *Client) newQueueClient(taskQueueName string) (client.Client, error) {
	tc, err := clientFromOpts(c.opts, taskQueueName)

	if err != nil {
		return nil, err
	}

	c.clients[taskQueueName] = tc

	return tc, nil
}

func (c *Client) Close() {
	c.Close()
}

type authHeadersProvider struct {
	headers map[string]string
}

func (a authHeadersProvider) GetHeaders(ctx context.Context) (map[string]string, error) {
	return a.headers, nil
}

func getHostName() string {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "Unknown"
	}
	return hostName
}
