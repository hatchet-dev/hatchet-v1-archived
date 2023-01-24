package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
)

// GithubAppInstallationRepository uses gorm.DB for querying the database
type GithubAppInstallationRepository struct {
	db *gorm.DB
}

// NewGithubAppInstallationRepository returns a DefaultGithubAppInstallationRepository which uses
// gorm.DB for querying the database
func NewGithubAppInstallationRepository(db *gorm.DB) repository.GithubAppInstallationRepository {
	return &GithubAppInstallationRepository{db}
}

func (repo *GithubAppInstallationRepository) CreateGithubAppInstallation(gai *models.GithubAppInstallation) (*models.GithubAppInstallation, repository.RepositoryError) {
	if err := repo.db.Create(gai).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return gai, nil
}

func (repo *GithubAppInstallationRepository) ReadGithubAppInstallationByInstallationAndAccountID(installationID, accountID int64) (*models.GithubAppInstallation, repository.RepositoryError) {
	gai := &models.GithubAppInstallation{}

	if err := repo.db.Where("installation_id = ? AND account_id = ?", installationID, accountID).First(&gai).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return gai, nil
}

func (repo *GithubAppInstallationRepository) ListGithubAppInstallationsByUserID(userID string, opts ...repository.QueryOption) ([]*models.GithubAppInstallation, *repository.PaginatedResult, repository.RepositoryError) {
	// get the corresponding oauth integration
	gao := &models.GithubAppOAuth{}

	if err := repo.db.Where("user_id = ?", userID).First(&gao).Error; err != nil {
		return nil, nil, toRepoError(repo.db, err)
	}

	var gais []*models.GithubAppInstallation

	db := repo.db.Model(&models.GithubAppInstallation{}).Where("github_app_o_auth_id = ?", gao.ID)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&gais).Error; err != nil {
		return nil, nil, err
	}

	return gais, paginatedResult, nil
}

func (repo *GithubAppInstallationRepository) UpdateGithubAppInstallation(
	gai *models.GithubAppInstallation,
) (*models.GithubAppInstallation, repository.RepositoryError) {
	if err := repo.db.Save(gai).Error; err != nil {
		return nil, err
	}

	return gai, nil
}

func (repo *GithubAppInstallationRepository) DeleteGithubAppInstallation(
	gai *models.GithubAppInstallation,
) (*models.GithubAppInstallation, repository.RepositoryError) {
	if err := repo.db.Delete(gai).Error; err != nil {
		return nil, err
	}

	return gai, nil
}
