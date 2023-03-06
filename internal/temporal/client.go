package temporal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
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
	BroadcastAddress string
	HostPort         string
	Namespace        string
	BearerToken      string
	DefaultQueueName string

	ClientKeyFile  string
	ClientCertFile string
	RootCAFile     string
	TLSServerName  string
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

	if opts.BearerToken != "" {
		tOpts.HeadersProvider = authHeadersProvider{
			headers: map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", opts.BearerToken),
			},
		}
	}

	if opts.ClientCertFile != "" && opts.ClientKeyFile != "" && opts.RootCAFile != "" {
		tlsConfig := &tls.Config{
			ServerName: opts.TLSServerName,
		}

		cert, err := tls.LoadX509KeyPair(
			opts.ClientCertFile,
			opts.ClientKeyFile,
		)

		tlsConfig.Certificates = []tls.Certificate{cert}

		caPool := x509.NewCertPool()
		var caBytes []byte

		caBytes, err = os.ReadFile(
			opts.RootCAFile,
		)

		if err != nil {
			return nil, fmt.Errorf("unable to load CA cert from file: %v", err)
		}

		if !caPool.AppendCertsFromPEM(caBytes) {
			return nil, errors.New("unknown failure constructing cert pool for ca")
		}

		tlsConfig.RootCAs = caPool

		tOpts.ConnectionOptions.TLS = tlsConfig
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

func (c *Client) GetBroadcastAddress() string {
	return c.opts.BroadcastAddress
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
