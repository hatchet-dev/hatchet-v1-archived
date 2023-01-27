package grpc

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/api/v1/server/pb"
)

func (s *ProvisionerServer) GetLog(moduleRun *pb.ModuleRun, server pb.Provisioner_GetLogServer) error {
	res, ok, err := verifyModuleRunToken(s.config, server.Context())

	if !ok {
		if err != nil {
			s.config.Logger.Error().Msgf("%s", err.Error())
		}

		return fmt.Errorf("unauthorized")
	}

	return s.config.DefaultLogStore.StreamLogs(server.Context(), getLogPathFromResult(res), &pbWriteCloser{server})
}

type pbWriteCloser struct {
	server pb.Provisioner_GetLogServer
}

func (w *pbWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), w.server.Send(&pb.LogString{
		Log: string(p),
	})
}

func (w *pbWriteCloser) Close() error {
	return nil
}
