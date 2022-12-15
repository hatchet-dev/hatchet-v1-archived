package gorm_test

import (
	"testing"
	"time"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserSession(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expires := time.Now().Add(24 * 30 * time.Hour)

		userSession, err := conf.Repository.UserSession().CreateUserSession(&models.UserSession{
			Key:       "1",
			Data:      []byte("1"),
			ExpiresAt: expires,
		})

		if err != nil {
			t.Fatalf("could not create user: %v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(userSession.ID))

		assert.Equal(t, "1", userSession.Key)
		assert.Equal(t, []byte("1"), userSession.Data)
		assert.Equal(t, userSession.ExpiresAt.Unix(), expires.Unix())

		return nil
	}, testutils.InitUsers)
}

func TestCreateUserSessionDuplicateKeys(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expires := time.Now().Add(24 * 30 * time.Hour)

		_, failingErr := conf.Repository.UserSession().CreateUserSession(&models.UserSession{
			Key:       testutils.InitDataAll.UserSessions[0].Key,
			Data:      []byte("1"),
			ExpiresAt: expires,
		})

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryUniqueConstraintFailed, failingErr)

		return nil
	}, testutils.InitUsers, testutils.InitUserSessions)
}

func TestReadUserSessionByKey(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expUserSess := testutils.InitDataAll.UserSessions[0]
		userSession, err := conf.Repository.UserSession().ReadUserSessionByKey(expUserSess.Key)

		assert.Nil(t, err, "error should be nil")

		testutils.AssertUserSessionsEqual(t, userSession, expUserSess)

		return nil
	}, testutils.InitUsers, testutils.InitUserSessions)
}

func TestReadUserSessionByKeyNotFound(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		_, failingErr := conf.Repository.UserSession().ReadUserSessionByKey("not-set")

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)

		return nil
	}, testutils.InitUsers, testutils.InitUserSessions)
}

func TestUpdateUserSession(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expUserSess := testutils.InitDataAll.UserSessions[0]
		expUserSess.Data = []byte("newdata")

		userSession, err := conf.Repository.UserSession().UpdateUserSession(expUserSess)

		assert.Nil(t, err, "error is nil")

		testutils.AssertUserSessionsEqual(t, userSession, expUserSess)

		// ensure session was written properly
		readUserSession, err := conf.Repository.UserSession().ReadUserSessionByKey(expUserSess.Key)

		assert.Nil(t, err, "error is nil")

		testutils.AssertUserSessionsEqual(t, readUserSession, expUserSess)

		return nil
	}, testutils.InitUsers, testutils.InitUserSessions)
}

func TestDeleteUserSession(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expUserSess := testutils.InitDataAll.UserSessions[0]

		userSession, err := conf.Repository.UserSession().DeleteUserSession(expUserSess)

		assert.Nil(t, err, "err is nil")

		testutils.AssertUserSessionsEqual(t, userSession, expUserSess)

		// ensure session cannot be read
		_, failingErr := conf.Repository.UserSession().ReadUserSessionByKey(expUserSess.Key)

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)

		return nil
	}, testutils.InitUsers, testutils.InitUserSessions)
}
