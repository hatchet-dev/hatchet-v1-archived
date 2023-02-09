package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
)

// ModuleEnvVarsRepository uses gorm.DB for querying the database
type ModuleEnvVarsRepository struct {
	db  *gorm.DB
	key *[32]byte
}

func NewModuleEnvVarsRepository(db *gorm.DB, key *[32]byte) repository.ModuleEnvVarsRepository {
	return &ModuleEnvVarsRepository{db, key}
}

func (repo *ModuleEnvVarsRepository) CreateModuleEnvVarsVersion(mev *models.ModuleEnvVarsVersion) (*models.ModuleEnvVarsVersion, repository.RepositoryError) {
	err := mev.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Create(mev).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	err = mev.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return mev, nil
}

func (repo *ModuleEnvVarsRepository) ReadModuleEnvVarsVersionByID(moduleID, moduleEnvVarsVersionID string) (*models.ModuleEnvVarsVersion, repository.RepositoryError) {
	mev := &models.ModuleEnvVarsVersion{}

	if err := repo.db.Where("module_id = ? AND id = ?", moduleID, moduleEnvVarsVersionID).First(&mev).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	err := mev.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return mev, nil
}

func (repo *ModuleEnvVarsRepository) ListModuleValueVersionsByModuleID(
	moduleID string,
	opts ...repository.QueryOption,
) ([]*models.ModuleEnvVarsVersion, *repository.PaginatedResult, repository.RepositoryError) {
	var mevs []*models.ModuleEnvVarsVersion

	db := repo.db.Model(&models.ModuleEnvVarsVersion{}).Where("module_id = ?", moduleID)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&mevs).Error; err != nil {
		return nil, nil, err
	}

	return mevs, paginatedResult, nil
}

func (repo *ModuleEnvVarsRepository) DeleteModuleEnvVarsVersion(mev *models.ModuleEnvVarsVersion) (*models.ModuleEnvVarsVersion, repository.RepositoryError) {
	err := mev.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Delete(mev).Error; err != nil {
		return nil, err
	}

	err = mev.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return mev, nil
}
