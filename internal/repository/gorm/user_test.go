package gorm_test

import (
	"testing"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		user, err := conf.Repository.User().CreateUser(&models.User{
			DisplayName: "A Belanger",
			Email:       "abelanger5@hatchet.run",
			Password:    "Abcdefgh123",
		})

		if err != nil {
			t.Fatalf("could not create user: %v", err)
		}

		assert.Equal(t, "abelanger5@hatchet.run", user.Email)
		assert.Equal(t, "A Belanger", user.DisplayName)

		return nil
	})
}

func TestCreateDuplicateUser(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		duplUser1 := models.User{
			DisplayName: "A Belanger",
			Email:       "abelanger5@hatchet.run",
			Password:    "Abcdefgh123",
		}

		duplUser2 := duplUser1

		_, err := conf.Repository.User().CreateUser(&duplUser1)

		if err != nil {
			t.Fatalf("could not create first user: %v", err)
		}

		_, failingErr := conf.Repository.User().CreateUser(&duplUser2)

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryUniqueConstraintFailed, failingErr)

		return nil
	})
}

func TestReadUser(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		user, err := conf.Repository.User().ReadUserByEmail(testutils.DeclaredUserModels[0].Email)

		if err != nil {
			t.Fatalf("could not read user: %v", err)
		}

		assert.Equal(t, testutils.DeclaredUserModels[0].Email, user.Email, "email of queried user should match")

		return nil
	}, testutils.InitUsers)
}

func TestFailedReadUser(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		_, failingErr := conf.Repository.User().ReadUserByEmail("notanemail@gmail.com")

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)

		return nil
	}, testutils.InitUsers)
}
