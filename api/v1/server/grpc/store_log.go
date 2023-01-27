package grpc

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hatchet-dev/hatchet/api/v1/server/pb"
)

func (s *ProvisionerServer) StoreLog(stream pb.Provisioner_StoreLogServer) error {
	res, ok, err := verifyModuleRunToken(s.config, stream.Context())

	if !ok {
		if err != nil {
			s.config.Logger.Error().Msgf("%s", err.Error())
		}

		return fmt.Errorf("unauthorized")
	}

	for {
		tfLog, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.TerraformStateMeta{})
		} else if err != nil {
			return err
		}

		tfLogBytes, err := json.Marshal(tfLog)

		if err != nil {
			return err
		}

		if s.config.DefaultLogStore != nil {
			err = s.config.DefaultLogStore.PushLogLine(stream.Context(), getLogPathFromResult(res), tfLogBytes)

			if err != nil {
				return err
			}
		}
	}
}
