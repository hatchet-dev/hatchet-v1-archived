package repository

import (
	"github.com/hatchet-dev/hatchet/internal/models"
)

// UserRepository represents the set of queries on the User model
type UserRepository interface {
	CreateUser(user *models.User) (*models.User, RepositoryError)
	ReadUserByEmail(email string) (*models.User, RepositoryError)
	ReadUserByID(id string) (*models.User, RepositoryError)
	DeleteUser(user *models.User) (*models.User, RepositoryError)
}
