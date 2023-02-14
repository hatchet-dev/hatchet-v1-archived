package grpcstreamer

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hatchet-dev/hatchet/api/v1/client/grpc"
	"github.com/hatchet-dev/hatchet/api/v1/server/pb"

	"github.com/hatchet-dev/hatchet/internal/runner/types"
)

type GRPCStreamer struct {
	stream      pb.Provisioner_StoreLogClient
	cancel      context.CancelFunc
	ctx         context.Context
	client      *grpc.GRPCClient
	workspaceID string
}

func NewGRPCStreamer(client *grpc.GRPCClient) (*GRPCStreamer, error) {
	ctx, cancel := client.NewGRPCContext()

	res := &GRPCStreamer{
		cancel: cancel,
		ctx:    ctx,
		client: client,
	}

	err := res.setStream()

	if err != nil {
		cancel()
		return nil, err
	}

	return res, nil
}

func (g *GRPCStreamer) setStream() (err error) {
	stream, err := g.client.GRPCClient.StoreLog(g.ctx)

	if err != nil {
		return err
	}

	g.stream = stream

	return nil
}

func (g *GRPCStreamer) Write(p []byte) (int, error) {
	for _, line := range bytes.Split(p, []byte("\n")) {
		var tfLog *types.TFLogLine = &types.TFLogLine{}

		if len(line) == 0 {
			continue
		}

		err := json.Unmarshal(line, tfLog)

		if err != nil {
			continue
		}

		err = g.stream.Send(tfLog.ToPBType())

		if err != nil {
			continue
		}
	}

	return len(p), nil
}
