package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
)

type ModuleMonitorRepository struct {
	db *gorm.DB
}

func NewModuleMonitorRepository(db *gorm.DB) repository.ModuleMonitorRepository {
	return &ModuleMonitorRepository{db}
}

func (repo *ModuleMonitorRepository) CreateModuleMonitor(monitor *models.ModuleMonitor) (*models.ModuleMonitor, repository.RepositoryError) {
	if err := repo.db.Create(monitor).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return monitor, nil
}

func (repo *ModuleMonitorRepository) ReadModuleMonitorByID(teamID, moduleMonitorID string) (*models.ModuleMonitor, repository.RepositoryError) {
	monitor := &models.ModuleMonitor{}

	if err := repo.db.Where("team_id = ? AND id = ?", teamID, moduleMonitorID).First(&monitor).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return monitor, nil
}

func (repo *ModuleMonitorRepository) ListModuleMonitorsByTeamID(teamID string, opts ...repository.QueryOption) ([]*models.ModuleMonitor, *repository.PaginatedResult, repository.RepositoryError) {
	var results []*models.ModuleMonitor

	query := repo.db.Model(&models.ModuleMonitor{}).Where("team_id = ?", teamID)

	paginatedResult := &repository.PaginatedResult{}

	query = query.Scopes(queryutils.Paginate(opts, query, paginatedResult))

	if err := query.Find(&results).Error; err != nil {
		return nil, nil, err
	}

	return results, paginatedResult, nil
}

func (repo *ModuleMonitorRepository) UpdateModuleMonitor(monitor *models.ModuleMonitor) (*models.ModuleMonitor, repository.RepositoryError) {
	if err := repo.db.Save(monitor).Error; err != nil {
		return nil, err
	}

	return monitor, nil
}

func (repo *ModuleMonitorRepository) DeleteModuleMonitor(monitor *models.ModuleMonitor) (*models.ModuleMonitor, repository.RepositoryError) {
	if err := repo.db.Delete(monitor).Error; err != nil {
		return nil, err
	}

	return monitor, nil
}

func (repo *ModuleMonitorRepository) CreateModuleMonitorResult(monitor *models.ModuleMonitor, result *models.ModuleMonitorResult) (*models.ModuleMonitorResult, repository.RepositoryError) {
	if err := repo.db.Create(result).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return result, nil
}

func (repo *ModuleMonitorRepository) ListModuleMonitorResults(
	teamID string,
	filterOpts *repository.ListModuleMonitorResultsOpts,
	opts ...repository.QueryOption,
) ([]*models.ModuleMonitorResult, *repository.PaginatedResult, repository.RepositoryError) {
	var results []*models.ModuleMonitorResult

	query := repo.db.Model(&models.ModuleMonitorResult{}).Where("team_id = ?", teamID)

	if filterOpts.ModuleID != "" {
		query = query.Where("module_id = ?", filterOpts.ModuleID)
	}

	if filterOpts.ModuleMonitorID != "" {
		query = query.Where("module_monitor_id = ?", filterOpts.ModuleMonitorID)
	}

	if filterOpts.Result != "" {
		query = query.Where("result = ?", filterOpts.Result)
	}

	if filterOpts.Severity != "" {
		query = query.Where("severity = ?", filterOpts.Severity)
	}

	paginatedResult := &repository.PaginatedResult{}

	query = query.Scopes(queryutils.Paginate(opts, query, paginatedResult))

	if err := query.Find(&results).Error; err != nil {
		return nil, nil, err
	}

	return results, paginatedResult, nil
}

func (repo *ModuleMonitorRepository) ReadModuleMonitorResultByID(moduleID, monitorID, resultID string) (*models.ModuleMonitorResult, repository.RepositoryError) {
	result := &models.ModuleMonitorResult{}

	if err := repo.db.Where("module_id = ? AND module_monitor_id = ? AND id = ?", moduleID, monitorID, resultID).First(&result).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return result, nil
}

func (repo *ModuleMonitorRepository) DeleteModuleMonitorResult(result *models.ModuleMonitorResult) (*models.ModuleMonitorResult, repository.RepositoryError) {
	if err := repo.db.Delete(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
