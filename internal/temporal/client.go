package temporal

import (
	"context"

	"go.temporal.io/sdk/client"
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

type authHeadersProvider struct {
	headers map[string]string
}

func (a authHeadersProvider) GetHeaders(ctx context.Context) (map[string]string, error) {
	return a.headers, nil
}
