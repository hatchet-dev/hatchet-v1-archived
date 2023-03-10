package temporal

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"go.temporal.io/sdk/client"
)

const DefaultQueueName = "default"

type Client struct {
	clients map[string]client.Client

	mu sync.Mutex

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

	clients := make(map[string]client.Client)

	res := &Client{clients, sync.Mutex{}, opts}

	err := res.eventualClientFromOpts(opts, opts.DefaultQueueName, 5)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) eventualClientFromOpts(opts *ClientOpts, taskQueueName string, maxRetries uint) error {
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
			return fmt.Errorf("unable to load CA cert from file: %v", err)
		}

		if !caPool.AppendCertsFromPEM(caBytes) {
			return errors.New("unknown failure constructing cert pool for ca")
		}

		tlsConfig.RootCAs = caPool

		tOpts.ConnectionOptions.TLS = tlsConfig
	}

	getter := func() {
		var err error
		var tClient client.Client

		for i := 0; i < int(maxRetries); i++ {
			tClient, err = client.Dial(tOpts)

			if err == nil {
				c.mu.Lock()
				c.clients[taskQueueName] = tClient
				c.mu.Unlock()
				break
			} else {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("could not create temporal client for queue %s: %s. Retrying (attempt %d of %d)...\n", taskQueueName, err.Error(), i+1, maxRetries))
				time.Sleep(5 * time.Second)
			}
		}

		if err != nil {
			// TODO: use shared logger here
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Fatal: could not create temporal client for queue %s: %s\n", taskQueueName, err.Error()))
		}
	}

	if maxRetries == 1 {
		getter()
	} else {
		go getter()
	}

	return nil
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
	err := c.eventualClientFromOpts(c.opts, taskQueueName, 1)

	if err != nil {
		return nil, err
	}

	return c.clients[taskQueueName], nil
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
