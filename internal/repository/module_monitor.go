package repository

import "github.com/hatchet-dev/hatchet/internal/models"

type ListModuleMonitorResultsOpts struct {
	ModuleID        string
	ModuleMonitorID string
	Severity        models.MonitorResultSeverity
	Result          models.MonitorResultStatus
}

// ModuleMonitorRepository represents the set of queries on the ModuleMonitor model
type ModuleMonitorRepository interface {
	// --- Module monitor queries ---
	//
	// CreateModuleMonitor creates a new module monitor in the database, associating it
	// with the parent module
	CreateModuleMonitor(monitor *models.ModuleMonitor) (*models.ModuleMonitor, RepositoryError)

	// ReadModuleMonitorByID reads the module by its unique UUID
	ReadModuleMonitorByID(teamID, moduleMonitorID string) (*models.ModuleMonitor, RepositoryError)

	// UpdateModuleMonitor updates a module monitor in the database
	UpdateModuleMonitor(monitor *models.ModuleMonitor) (*models.ModuleMonitor, RepositoryError)

	// ReplaceModuleMonitorModules replaces module monitors in the database
	ReplaceModuleMonitorModules(monitor *models.ModuleMonitor, modules []*models.Module) (*models.ModuleMonitor, RepositoryError)

	// ListModuleMonitorsByTeamID lists the module monitors by the team id
	ListModuleMonitorsByTeamID(teamID string, opts ...QueryOption) ([]*models.ModuleMonitor, *PaginatedResult, RepositoryError)

	// DeleteModuleMonitor soft-deletes and module monitor in the database
	DeleteModuleMonitor(monitor *models.ModuleMonitor) (*models.ModuleMonitor, RepositoryError)

	// --- Module monitor result queries ---
	//
	// CreateModuleMonitorResult creates a new module monitor result in the database, associating it
	// with the parent module and monitor
	CreateModuleMonitorResult(monitor *models.ModuleMonitor, result *models.ModuleMonitorResult) (*models.ModuleMonitorResult, RepositoryError)

	// ReadModuleMonitorResultByID reads the first module result corresponding to that module monitor id
	ReadModuleMonitorResultByID(moduleID, monitorID, resultID string) (*models.ModuleMonitorResult, RepositoryError)

	// ListModuleMonitorResults lists the module results based on a set of filters
	ListModuleMonitorResults(teamID string, filterOpts *ListModuleMonitorResultsOpts, opts ...QueryOption) ([]*models.ModuleMonitorResult, *PaginatedResult, RepositoryError)

	// DeleteModuleMonitorResult soft-deletes a module run queue item in the database
	DeleteModuleMonitorResult(result *models.ModuleMonitorResult) (*models.ModuleMonitorResult, RepositoryError)
}
