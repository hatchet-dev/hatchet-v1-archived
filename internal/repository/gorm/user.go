package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

// UserRepository uses gorm.DB for querying the database
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a DefaultUserRepository which uses
// gorm.DB for querying the database
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db}
}

// CreateUser adds a new User row to the Users table in the database
func (repo *UserRepository) CreateUser(user *models.User) (*models.User, repository.RepositoryError) {
	if err := repo.db.Create(user).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}
	return user, nil
}

// ReadUserByID finds a single user based on their unique id
func (repo *UserRepository) ReadUserByID(id string) (*models.User, repository.RepositoryError) {
	user := &models.User{}

	if err := repo.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return user, nil
}

// ReadUserByEmail finds a single user based on their unique email
func (repo *UserRepository) ReadUserByEmail(email string) (*models.User, repository.RepositoryError) {
	user := &models.User{}

	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return user, nil
}

// UpdateUser updates a user in the database
func (repo *UserRepository) UpdateUser(user *models.User) (*models.User, repository.RepositoryError) {
	if err := repo.db.Save(user).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return user, nil
}

// DeleteUser deletes a single user using their unique id
func (repo *UserRepository) DeleteUser(user *models.User) (*models.User, repository.RepositoryError) {
	del := repo.db.Delete(&user)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return user, nil
}
