package grpc

import (
	"context"
	"fmt"

	"github.com/hatchet-dev/hatchet/api/v1/server/pb"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage"
	"github.com/hatchet-dev/hatchet/internal/models"
	"google.golang.org/grpc/metadata"
)

type ProvisionerServer struct {
	pb.UnimplementedProvisionerServer

	config *server.Config
}

func NewProvisionerServer(config *server.Config) *ProvisionerServer {
	return &ProvisionerServer{
		config: config,
	}
}

type verifyResult struct {
	mrt                           *models.ModuleRunToken
	teamID, moduleID, moduleRunID string
}

// verifyModuleRunToken ensures that the MRT token passed through auth is valid
func verifyModuleRunToken(config *server.Config, ctx context.Context) (res *verifyResult, valid bool, err error) {
	streamContext, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, false, fmt.Errorf("stream context could not be read")
	}

	tokenArr, exists := streamContext["token"]

	if !exists || len(tokenArr) != 1 {
		return nil, false, fmt.Errorf("stream context does not contain correct token parameter")
	}

	mrt, err := token.GetMRTFromEncoded(tokenArr[0], config.DB.Repository.Module(), config.TokenOpts)

	if err != nil {
		return nil, false, err
	}

	if mrt.Revoked || mrt.IsExpired() {
		return nil, false, fmt.Errorf("token with id %s not valid", mrt.ID)
	}

	teamID, exists := streamContext["team_id"]

	if !exists || len(teamID) != 1 {
		return nil, false, fmt.Errorf("team_id not found in stream context")
	}

	moduleID, exists := streamContext["module_id"]

	if !exists || len(moduleID) != 1 {
		return nil, false, fmt.Errorf("module_id not found in stream context")
	}

	moduleRunID, exists := streamContext["module_run_id"]

	if !exists || len(moduleRunID) != 1 {
		return nil, false, fmt.Errorf("module_run_id not found in stream context")
	}

	// look at the team id, module id, and module run id and verify that it matches the MRT
	module, err := config.DB.Repository.Module().ReadModuleByID(teamID[0], moduleID[0])

	if err != nil {
		return nil, false, err
	}

	moduleRun, err := config.DB.Repository.Module().ReadModuleRunByID(module.ID, moduleRunID[0])

	if err != nil {
		return nil, false, err
	}

	if moduleRun.ID != mrt.ModuleRunID {
		return nil, false, fmt.Errorf("module run id [%s] does not match token's module run id [%s]", moduleRun.ID, mrt.ModuleRunID)
	}

	return &verifyResult{
		mrt:         mrt,
		teamID:      teamID[0],
		moduleID:    module.ID,
		moduleRunID: moduleRun.ID,
	}, true, nil
}

func getLogPathFromResult(res *verifyResult) string {
	return logstorage.GetLogStoragePath(res.teamID, res.moduleID, res.moduleRunID)
}
