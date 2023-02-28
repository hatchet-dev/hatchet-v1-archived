package provisioner

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type ProvisionOpts struct {
	Env []string

	WaitForRunFinished bool
}

type Provisioner interface {
	RunPlan(opts *ProvisionOpts) error
	RunApply(opts *ProvisionOpts) error
	RunDestroy(opts *ProvisionOpts) error

	RunStateMonitor(opts *ProvisionOpts, monitorID string, policy []byte) error
	RunPlanMonitor(opts *ProvisionOpts, monitorID string, policy []byte) error
	RunBeforePlanMonitor(opts *ProvisionOpts, monitorID string, policy []byte) error
	RunAfterPlanMonitor(opts *ProvisionOpts, monitorID string, policy []byte) error
	RunBeforeApplyMonitor(opts *ProvisionOpts, monitorID string, policy []byte) error
	RunAfterApplyMonitor(opts *ProvisionOpts, monitorID string, policy []byte) error
}

type GetEnvOpts struct {
	Team      *models.Team
	Module    *models.Module
	ModuleRun *models.ModuleRun
	EnvVars   map[string]string

	TokenOpts  token.TokenOpts
	Repository repository.Repository
	ServerURL  string
}

func GetHatchetRunnerEnv(opts *GetEnvOpts, currEnv []string) ([]string, error) {
	tok, err := GetRunnerToken(opts)

	if err != nil {
		return nil, err
	}

	if currEnv == nil {
		currEnv = make([]string, 0)
	}

	// TODO: configure grpc server address separately
	// env = append(env, "GRPC_SERVER_ADDRESS=%s", config.ServerRuntimeConfig.ServerURL)
	currEnv = append(currEnv, fmt.Sprintf("GRPC_TOKEN=%s", tok))
	currEnv = append(currEnv, fmt.Sprintf("TEAM_ID=%s", opts.Team.ID))
	currEnv = append(currEnv, fmt.Sprintf("MODULE_ID=%s", opts.Module.ID))
	currEnv = append(currEnv, fmt.Sprintf("MODULE_RUN_ID=%s", opts.ModuleRun.ID))
	currEnv = append(currEnv, fmt.Sprintf("API_TOKEN=%s", tok))

	// TODO: configure api server address separately
	// env = append(env, fmt.Sprintf("API_SERVER_ADDRESS=%s", opts.ServerURL))
	currEnv = append(currEnv, fmt.Sprintf("GITHUB_SHA=%s", opts.ModuleRun.ModuleRunConfig.GithubCommitSHA))
	currEnv = append(currEnv, fmt.Sprintf("GITHUB_REPOSITORY_NAME=%s", opts.Module.DeploymentConfig.GithubRepoName))
	currEnv = append(currEnv, fmt.Sprintf("GITHUB_MODULE_PATH=%s", opts.Module.DeploymentConfig.ModulePath))

	for key, val := range opts.EnvVars {
		currEnv = append(currEnv, fmt.Sprintf("%s=%s", key, val))
	}

	return currEnv, nil
}

func GetRunnerToken(opts *GetEnvOpts) (string, error) {
	mrt, err := models.NewModuleRunTokenFromRunID(opts.Team.ServiceAccountRunnerID, opts.ModuleRun.ID)

	if err != nil {
		return "", err
	}

	rawTok, err := token.GenerateTokenFromMRT(mrt, &opts.TokenOpts)

	if err != nil {
		return "", err
	}

	mrt, err = opts.Repository.Module().CreateModuleRunToken(mrt)

	if err != nil {
		return "", err
	}

	return rawTok, nil
}
