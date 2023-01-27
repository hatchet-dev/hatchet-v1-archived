package grpc

import (
	"context"
	"net/url"

	"github.com/hatchet-dev/hatchet/api/v1/server/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GRPCClient struct {
	BaseURL                              string
	Token, TeamID, ModuleID, ModuleRunID string
	GRPCClient                           pb.ProvisionerClient

	conn *grpc.ClientConn
}

func NewGRPCClient(baseURL, token, teamID, moduleID, moduleRunID string) (*GRPCClient, error) {
	parsedURL, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(parsedURL.Host, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	gClient := pb.NewProvisionerClient(conn)

	client := &GRPCClient{
		BaseURL:     baseURL,
		Token:       token,
		TeamID:      teamID,
		ModuleID:    moduleID,
		ModuleRunID: moduleRunID,
		GRPCClient:  gClient,
		conn:        conn,
	}

	return client, nil
}

func (c *GRPCClient) NewGRPCContext() (context.Context, context.CancelFunc) {
	headers := map[string]string{
		"token":         c.Token,
		"team_id":       c.TeamID,
		"module_id":     c.ModuleID,
		"module_run_id": c.ModuleRunID,
	}

	header := metadata.New(headers)

	ctx := metadata.NewOutgoingContext(context.Background(), header)

	return context.WithCancel(ctx)
}

func (c *GRPCClient) CloseConnection() error {
	return c.conn.Close()
}
