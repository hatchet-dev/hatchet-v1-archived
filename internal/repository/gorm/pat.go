package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
)

// PersonalAccessTokenRepository uses gorm.DB for querying the database
type PersonalAccessTokenRepository struct {
	db  *gorm.DB
	key *[32]byte
}

// NewPersonalAccessTokenRepository returns a DefaultPersonalAccessTokenRepository which uses
// gorm.DB for querying the database
func NewPersonalAccessTokenRepository(db *gorm.DB, key *[32]byte) repository.PersonalAccessTokenRepository {
	return &PersonalAccessTokenRepository{db, key}
}

// CreateUser adds a new User row to the Users table in the database
func (repo *PersonalAccessTokenRepository) CreatePersonalAccessToken(pat *models.PersonalAccessToken) (*models.PersonalAccessToken, repository.RepositoryError) {
	err := pat.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Create(pat).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	// return the PAT decrypted
	err = pat.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return pat, nil
}

// ReadPersonalAccessToken reads a PAT by both the user ID and the token ID
func (repo *PersonalAccessTokenRepository) ReadPersonalAccessToken(userID, tokenID string) (*models.PersonalAccessToken, repository.RepositoryError) {
	pat := &models.PersonalAccessToken{}

	if err := repo.db.Where("user_id = ? AND id = ?", userID, tokenID).First(&pat).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	err := pat.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return pat, nil
}

func (repo *PersonalAccessTokenRepository) ReadPersonalAccessTokenByDisplayName(userID, displayName string) (
	*models.PersonalAccessToken,
	repository.RepositoryError,
) {
	pat := &models.PersonalAccessToken{}

	if err := repo.db.Where("user_id = ? AND display_name = ?", userID, displayName).First(&pat).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	err := pat.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return pat, nil
}

func (repo *PersonalAccessTokenRepository) ListPersonalAccessTokensByUserID(
	userID string,
	opts ...repository.QueryOption,
) ([]*models.PersonalAccessToken, *repository.PaginatedResult, repository.RepositoryError) {

	var pats []*models.PersonalAccessToken

	db := repo.db.Model(&models.PersonalAccessToken{})

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&pats).Error; err != nil {
		return nil, nil, err
	}

	return pats, paginatedResult, nil
}

func (repo *PersonalAccessTokenRepository) UpdatePersonalAccessToken(
	pat *models.PersonalAccessToken,
) (*models.PersonalAccessToken, repository.RepositoryError) {
	err := pat.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Save(pat).Error; err != nil {
		return nil, err
	}

	// return the PAT decrypted
	err = pat.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return pat, nil
}

func (repo *PersonalAccessTokenRepository) DeletePersonalAccessToken(
	pat *models.PersonalAccessToken,
) (*models.PersonalAccessToken, repository.RepositoryError) {
	// encrypt the PAT
	err := pat.Encrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	if err := repo.db.Delete(pat).Error; err != nil {
		return nil, err
	}

	// return the PAT decrypted
	err = pat.Decrypt(repo.key)

	if err != nil {
		return nil, repository.UnknownRepositoryError(err)
	}

	return pat, nil
}
