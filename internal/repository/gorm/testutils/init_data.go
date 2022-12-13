package testutils

import (
	"testing"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
)

var DeclaredUserModels []models.User = []models.User{
	{
		Email:         "user1@gmail.com",
		Password:      "Abcdefgh123",
		DisplayName:   "User 1",
		EmailVerified: true,
	},
	{
		Email:         "user2@gmail.com",
		Password:      "Abcdefgh123",
		DisplayName:   "User 2",
		EmailVerified: false,
	},
}

type InitData struct {
	Users []*models.User
}

type InitDataFunc func(t *testing.T, conf *database.Config, i *InitData) error

func InitUsers(t *testing.T, conf *database.Config, i *InitData) error {
	users := make([]*models.User, 0)

	for _, declaredUser := range DeclaredUserModels {
		userCp := declaredUser

		user, err := conf.Repository.User().CreateUser(&userCp)

		if err != nil {
			return err
		}

		users = append(users, user)
	}

	i.Users = users

	return nil
}
