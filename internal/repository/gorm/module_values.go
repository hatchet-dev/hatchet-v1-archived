package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
)

// ModuleValuesRepository uses gorm.DB for querying the database
type ModuleValuesRepository struct {
	db  *gorm.DB
	key *[32]byte
}

func NewModuleValuesRepository(db *gorm.DB, key *[32]byte) repository.ModuleValuesRepository {
	return &ModuleValuesRepository{db, key}
}

func (repo *ModuleValuesRepository) CreateModuleValuesVersion(mvv *models.ModuleValuesVersion) (*models.ModuleValuesVersion, repository.RepositoryError) {
	if err := repo.db.Create(mvv).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mvv, nil
}

// ReadPersonalAccessToken reads a PAT by both the user ID and the token ID
func (repo *ModuleValuesRepository) ReadModuleValuesVersionByID(moduleID, moduleValuesVersionID string) (*models.ModuleValuesVersion, repository.RepositoryError) {
	mvv := &models.ModuleValuesVersion{}

	if err := repo.db.Where("module_id = ? AND id = ?", moduleID, moduleValuesVersionID).First(&mvv).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return mvv, nil
}

func (repo *ModuleValuesRepository) ListModuleValueVersionsByModuleID(
	moduleID string,
	opts ...repository.QueryOption,
) ([]*models.ModuleValuesVersion, *repository.PaginatedResult, repository.RepositoryError) {
	var mvvs []*models.ModuleValuesVersion

	db := repo.db.Model(&models.ModuleValuesVersion{}).Where("module_id = ?", moduleID)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&mvvs).Error; err != nil {
		return nil, nil, err
	}

	return mvvs, paginatedResult, nil
}

func (repo *ModuleValuesRepository) DeleteModuleValuesVersion(mvv *models.ModuleValuesVersion) (*models.ModuleValuesVersion, repository.RepositoryError) {
	if err := repo.db.Delete(mvv).Error; err != nil {
		return nil, err
	}

	return mvv, nil
}

func (repo *ModuleValuesRepository) CreateModuleValues(mv *models.ModuleValues) (*models.ModuleValues, repository.RepositoryError) {
	err := mv.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Create(mv).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	// return the PAT decrypted
	err = mv.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return mv, nil
}

// ReadPersonalAccessToken reads a PAT by both the user ID and the token ID
func (repo *ModuleValuesRepository) ReadModuleValuesByID(moduleID, moduleValuesID string) (*models.ModuleValues, repository.RepositoryError) {
	mv := &models.ModuleValues{}

	if err := repo.db.Where("module_id = ? AND id = ?", moduleID, moduleValuesID).First(&mv).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	err := mv.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return mv, nil
}

func (repo *ModuleValuesRepository) ReadModuleValuesByVersionID(moduleID, moduleValuesVersionID string) (*models.ModuleValues, repository.RepositoryError) {
	mv := &models.ModuleValues{}

	if err := repo.db.Where("module_id = ? AND module_values_version_id = ?", moduleID, moduleValuesVersionID).First(&mv).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	err := mv.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return mv, nil
}

func (repo *ModuleValuesRepository) DeleteModuleValues(mv *models.ModuleValues) (*models.ModuleValues, repository.RepositoryError) {
	err := mv.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Delete(mv).Error; err != nil {
		return nil, err
	}

	// return the PAT decrypted
	err = mv.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return mv, nil
}
