package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ModuleRepository uses gorm.DB for querying the database
type ModuleRepository struct {
	db *gorm.DB
}

// NewModuleRepository returns a DefaultModuleRepository which uses
// gorm.DB for querying the database
func NewModuleRepository(db *gorm.DB) repository.ModuleRepository {
	return &ModuleRepository{db}
}

// CreateModule adds a new Module row to the Modules table in the database
func (repo *ModuleRepository) CreateModule(mod *models.Module) (*models.Module, repository.RepositoryError) {
	if err := repo.db.Create(mod).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mod, nil
}

// ReadModuleByID finds a single mod by its unique id
func (repo *ModuleRepository) ReadModuleByID(teamID, moduleID string) (*models.Module, repository.RepositoryError) {
	mod := &models.Module{}

	if err := repo.db.Preload("DeploymentConfig").Preload("Runs").Joins("CurrentModuleValuesVersion").Joins("CurrentModuleEnvVarsVersion").Where("team_id = ? AND modules.id = ?", teamID, moduleID).First(&mod).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mod, nil
}

// UpdateModule updates an module in the database
func (repo *ModuleRepository) UpdateModule(mod *models.Module) (*models.Module, repository.RepositoryError) {
	if err := repo.db.Omit(clause.Associations).Save(mod).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mod, nil
}

// DeleteModule deletes a single mod by its unique id
func (repo *ModuleRepository) DeleteModule(mod *models.Module) (*models.Module, repository.RepositoryError) {
	del := repo.db.Delete(&mod)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return mod, nil
}

func (repo *ModuleRepository) ListModulesByTeamID(teamID string, opts ...repository.QueryOption) ([]*models.Module, *repository.PaginatedResult, repository.RepositoryError) {
	var mods []*models.Module

	db := repo.db.Model(&models.Module{}).Where("team_id = ?", teamID)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Preload("DeploymentConfig").Find(&mods).Error; err != nil {
		return nil, nil, err
	}

	return mods, paginatedResult, nil
}

func (repo *ModuleRepository) ListModulesByIDs(teamID string, ids []string, opts ...repository.QueryOption) ([]*models.Module, *repository.PaginatedResult, repository.RepositoryError) {
	var mods []*models.Module

	db := repo.db.Model(&models.Module{}).Where("team_id = ? AND id IN (?)", teamID, ids)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Preload("DeploymentConfig").Find(&mods).Error; err != nil {
		return nil, nil, err
	}

	return mods, paginatedResult, nil
}

func (repo *ModuleRepository) ListGithubRepositoryModules(teamID, repoOwner, repoName string) ([]*models.Module, repository.RepositoryError) {
	var mods []*models.Module

	db := repo.db.Joins("DeploymentConfig").Where("team_id = ? AND DeploymentConfig.github_repo_owner = ? AND DeploymentConfig.github_repo_name = ?", teamID, repoOwner, repoName)

	if err := db.Find(&mods).Error; err != nil {
		return nil, err
	}

	return mods, nil
}

func (repo *ModuleRepository) CreateModuleRun(run *models.ModuleRun) (*models.ModuleRun, repository.RepositoryError) {
	if err := repo.db.Create(run).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return run, nil
}

func (repo *ModuleRepository) ReadModuleRunByID(moduleID, moduleRunID string) (*models.ModuleRun, repository.RepositoryError) {
	mod := &models.ModuleRun{}

	if err := repo.db.Joins("ModuleRunConfig").Where("module_id = ? AND module_runs.id = ?", moduleID, moduleRunID).First(&mod).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mod, nil
}

func (repo *ModuleRepository) ReadModuleRunWithStateLock(moduleID string) (*models.ModuleRun, repository.RepositoryError) {
	mod := &models.ModuleRun{}

	if err := repo.db.Joins("ModuleRunConfig").Where("module_id = ? AND lock_id != ''", moduleID).First(&mod).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mod, nil
}

func (repo *ModuleRepository) ListModuleRunsByGithubSHA(moduleID, githubSHA string, kind *models.ModuleRunKind) ([]*models.ModuleRun, repository.RepositoryError) {
	mods := make([]*models.ModuleRun, 0)

	db := repo.db.Joins("ModuleRunConfig").Where("module_id = ? AND ModuleRunConfig.github_commit_sha = ?", moduleID, githubSHA)

	if kind != nil {
		db = db.Where("kind = ?", kind)
	}

	if err := db.Find(&mods).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mods, nil
}

func (repo *ModuleRepository) UpdateModuleRun(run *models.ModuleRun) (*models.ModuleRun, repository.RepositoryError) {
	if err := repo.db.Omit("ModuleRunConfig").Save(run).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return run, nil
}

func (repo *ModuleRepository) DeleteModuleRun(run *models.ModuleRun) (*models.ModuleRun, repository.RepositoryError) {
	del := repo.db.Delete(&run)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return run, nil
}

func (repo *ModuleRepository) ListCompletedModuleRunsByLogLocation(location string, opts ...repository.QueryOption) ([]*models.ModuleRun, *repository.PaginatedResult, repository.RepositoryError) {
	var runs []*models.ModuleRun

	db := repo.db.Model(&models.ModuleRun{}).Where("status = ? AND log_location = ?", models.ModuleRunStatusCompleted, location)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&runs).Error; err != nil {
		return nil, nil, err
	}

	// populate the team ids in the runs as the team IDs are not implicit
	moduleToTeamID := make(map[string]string)

	for _, run := range runs {
		moduleToTeamID[run.ModuleID] = ""
	}

	moduleIDs := make([]string, 0)

	for moduleID := range moduleToTeamID {
		moduleIDs = append(moduleIDs, moduleID)
	}

	var mods []*models.Module

	if err := repo.db.Where("id IN (?)", moduleIDs).Find(&mods).Error; err != nil {
		return nil, nil, err
	}

	for _, mod := range mods {
		moduleToTeamID[mod.ID] = mod.TeamID
	}

	for _, run := range runs {
		run.TeamID = moduleToTeamID[run.ModuleID]
	}

	return runs, paginatedResult, nil
}

func (repo *ModuleRepository) ListRunsByModuleID(moduleID string, status *models.ModuleRunStatus, opts ...repository.QueryOption) ([]*models.ModuleRun, *repository.PaginatedResult, repository.RepositoryError) {
	var runs []*models.ModuleRun

	db := repo.db.Model(&models.ModuleRun{}).Where("module_id = ?", moduleID)

	if status != nil && *status != "" {
		db = db.Where("status = ?", *status)
	}

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&runs).Error; err != nil {
		return nil, nil, err
	}

	return runs, paginatedResult, nil
}

func (repo *ModuleRepository) CreateModuleRunToken(mrt *models.ModuleRunToken) (*models.ModuleRunToken, repository.RepositoryError) {
	if err := repo.db.Create(mrt).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mrt, nil
}

func (repo *ModuleRepository) ReadModuleRunToken(userID, runID, tokenID string) (*models.ModuleRunToken, repository.RepositoryError) {
	mrt := &models.ModuleRunToken{}

	if err := repo.db.Where("user_id = ? AND module_run_id = ? AND id = ?", userID, runID, tokenID).First(&mrt).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mrt, nil
}

func (repo *ModuleRepository) UpdateModuleRunToken(mrt *models.ModuleRunToken) (*models.ModuleRunToken, repository.RepositoryError) {
	if err := repo.db.Save(mrt).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mrt, nil
}

func (repo *ModuleRepository) DeleteModuleRunToken(mrt *models.ModuleRunToken) (*models.ModuleRunToken, repository.RepositoryError) {
	del := repo.db.Delete(&mrt)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return mrt, nil
}
