package gorm_test

import (
	"testing"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
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

		assert.True(t, uuidutils.IsValidUUID(user.ID))
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

func TestReadUserByEmail(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		user, err := conf.Repository.User().ReadUserByEmail(testutils.UserModels[0].Email)

		if err != nil {
			t.Fatalf("could not read user: %v", err)
		}

		assert.Equal(t, testutils.UserModels[0].Email, user.Email, "email of queried user should match")

		return nil
	}, testutils.InitUsers)
}

func TestFailedReadUserByEmail(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		_, failingErr := conf.Repository.User().ReadUserByEmail("notanemail@gmail.com")

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)

		return nil
	}, testutils.InitUsers)
}

func TestReadUserByID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		user, err := conf.Repository.User().ReadUserByID(testutils.UserModels[0].ID)

		if err != nil {
			t.Fatalf("could not read user: %v", err)
		}

		testutils.AssertUsersEqual(t, testutils.UserModels[0], user)

		return nil
	}, testutils.InitUsers)
}

func TestFailedReadUserByID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		_, failingErr := conf.Repository.User().ReadUserByID("not-an-id")

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)

		return nil
	}, testutils.InitUsers)
}

func TestDeleteUser(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		user, err := conf.Repository.User().DeleteUser(testutils.UserModels[0])

		if err != nil {
			t.Fatalf("could not read user: %v", err)
		}

		testutils.AssertUsersEqual(t, testutils.UserModels[0], user)

		// ensure that user no longer exists in the DB
		user, failingErr := conf.Repository.User().ReadUserByEmail(testutils.UserModels[0].Email)

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)
		assert.Nil(t, user, "user should be nil")

		return nil
	}, testutils.InitUsers)
}

func TestFailedDeleteUser(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		notAUser := &models.User{
			Base: models.Base{
				ID: "not-an-id",
			},
		}

		_, failingErr := conf.Repository.User().DeleteUser(notAUser)

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryNoRowsAffected, failingErr)

		return nil
	}, testutils.InitUsers)
}
