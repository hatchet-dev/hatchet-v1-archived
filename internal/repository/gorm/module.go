package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
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
func (repo *ModuleRepository) ReadModuleByID(id string) (*models.Module, repository.RepositoryError) {
	mod := &models.Module{}

	if err := repo.db.Preload("DeploymentConfig").Preload("Runs").Where("id = ?", id).First(&mod).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mod, nil
}

// UpdateModule updates an module in the database
func (repo *ModuleRepository) UpdateModule(mod *models.Module) (*models.Module, repository.RepositoryError) {
	if err := repo.db.Save(mod).Error; err != nil {
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
