package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// ModuleRepository represents the set of queries on the Module model
type ModuleRepository interface {
	// --- Module queries ---
	//
	// CreateModule creates a new module in the database
	CreateModule(mod *models.Module) (*models.Module, RepositoryError)

	// ReadModuleByID reads the module by its unique UUID
	ReadModuleByID(teamID, moduleID string) (*models.Module, RepositoryError)

	// ListModulesByTeamID lists all modules for a team
	ListModulesByTeamID(teamID string, opts ...QueryOption) ([]*models.Module, *PaginatedResult, RepositoryError)

	// ListModulesByIDs lists all modules matching a list of IDs
	ListModulesByIDs(teamID string, ids []string, opts ...QueryOption) ([]*models.Module, *PaginatedResult, RepositoryError)

	// ListGithubRepositoryModules lists modules that use the github deployment mechanism belonging
	// to a specific repo owner and name
	ListGithubRepositoryModules(teamID, repoOwner, repoName string) ([]*models.Module, RepositoryError)

	// UpdateModule updates any modified values for a module
	UpdateModule(module *models.Module) (*models.Module, RepositoryError)

	// DeleteModule soft-deletes a module
	DeleteModule(module *models.Module) (*models.Module, RepositoryError)

	// --- Run queries ---
	//
	// CreateModuleRun creates a new run in the database
	CreateModuleRun(run *models.ModuleRun) (*models.ModuleRun, RepositoryError)

	// ReadModuleRunByID reads the run by its unique UUID
	ReadModuleRunByID(moduleID, moduleRunID string) (*models.ModuleRun, RepositoryError)

	// ReadModuleRunByGithubSHA finds a run by its Github SHA
	ListModuleRunsByGithubSHA(moduleID, githubSHA string, kind *models.ModuleRunKind) ([]*models.ModuleRun, RepositoryError)

	// ListCompletedModuleRunsByLogLocation lists all module runs with a given log location
	ListCompletedModuleRunsByLogLocation(location string, opts ...QueryOption) ([]*models.ModuleRun, *PaginatedResult, RepositoryError)

	// ReadModuleRunWithStateLock returns a module run that has a lock on the module state
	ReadModuleRunWithStateLock(moduleID string) (*models.ModuleRun, RepositoryError)

	// ListRunsByModuleID lists all runs for a module
	ListRunsByModuleID(moduleID string, status *models.ModuleRunStatus, opts ...QueryOption) ([]*models.ModuleRun, *PaginatedResult, RepositoryError)

	// UpdateModuleRun updates any modified values for a module
	UpdateModuleRun(run *models.ModuleRun) (*models.ModuleRun, RepositoryError)

	// AppendModuleRunMonitors adds a list of monitors to a module run
	AppendModuleRunMonitors(run *models.ModuleRun, monitors []*models.ModuleMonitor) (*models.ModuleRun, RepositoryError)

	// AppendModuleRunMonitorResult adds a single monitor result to the module run monitor results
	AppendModuleRunMonitorResult(run *models.ModuleRun, result *models.ModuleMonitorResult) (*models.ModuleRun, RepositoryError)

	// DeleteModuleRun soft-deletes a run
	DeleteModuleRun(run *models.ModuleRun) (*models.ModuleRun, RepositoryError)

	// --- Run token queries ---
	//
	// CreateModuleRunToken creates a new module run token in the database
	CreateModuleRunToken(mrt *models.ModuleRunToken) (*models.ModuleRunToken, RepositoryError)

	// ReadModuleRunToken reads the module run token by its token ID
	ReadModuleRunToken(userID, runID, tokenID string) (*models.ModuleRunToken, RepositoryError)

	// UpdateModuleRunToken updates a module run token
	UpdateModuleRunToken(mrt *models.ModuleRunToken) (*models.ModuleRunToken, RepositoryError)

	// DeleteModuleRunToken soft-deletes a module run token in the DB
	DeleteModuleRunToken(mrt *models.ModuleRunToken) (*models.ModuleRunToken, RepositoryError)
}
