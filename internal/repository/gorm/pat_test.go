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

func TestCreatePAT(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expExpires := time.Now().Add(24 * 30 * time.Hour)

		displayName := "my-pat-1"

		// call population method
		pat, err := models.NewPATFromUserID(displayName, testutils.UserModels[0].ID)

		if err != nil {
			return err
		}

		pat, err = conf.Repository.PersonalAccessToken().CreatePersonalAccessToken(pat)

		if err != nil {
			return err
		}

		// ensure that id is set and fields are valid
		assert.True(t, uuidutils.IsValidUUID(pat.ID))
		assert.Equal(t, "my-pat-1", pat.DisplayName)
		assert.Equal(t, testutils.UserModels[0].ID, pat.UserID)
		assert.InDelta(t, expExpires.Unix(), pat.Expires.Unix(), 10)

		return nil
	}, testutils.InitUsers)
}

func TestReadPATSuccessful(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expPAT := testutils.PATModels[0]

		pat, err := conf.Repository.PersonalAccessToken().ReadPersonalAccessToken(expPAT.UserID, expPAT.ID)

		if err != nil {
			return err
		}

		testutils.AssertPATsEqual(t, pat, expPAT)

		return nil
	}, testutils.InitUsers, testutils.InitPATs)
}

func TestReadPATNotFound(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		_, failingErr := conf.Repository.PersonalAccessToken().ReadPersonalAccessToken(testutils.UserModels[0].ID, "not-an-id")

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)

		return nil
	}, testutils.InitUsers, testutils.InitPATs)
}

func TestListPATs(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expPAT := testutils.PATModels[0]
		expPAT.FieldsAreEncrypted = true

		pat, pagination, err := conf.Repository.PersonalAccessToken().ListPersonalAccessTokensByUserID(expPAT.UserID)

		if err != nil {
			return err
		}

		assert.Equal(t, 1, len(pat), "length of pats should be 1")

		assert.Equal(t, &repository.PaginatedResult{
			NumPages:    1,
			CurrentPage: 0,
			NextPage:    0,
		}, pagination, "pagination should be equal")

		// we skip looking at the signing secret as its still encrypted when returned from list
		pat[0].SigningSecret = expPAT.SigningSecret

		testutils.AssertPATsEqual(t, pat[0], expPAT)

		return nil
	}, testutils.InitUsers, testutils.InitPATs)
}

func TestUpdatePAT(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expPAT := testutils.PATModels[0]
		expPAT.Revoked = true

		pat, err := conf.Repository.PersonalAccessToken().UpdatePersonalAccessToken(expPAT)

		if err != nil {
			return err
		}

		assert.Nil(t, err, "err is nil")

		testutils.AssertPATsEqual(t, pat, expPAT)

		// ensure resulting PAT is revoked
		pat, err = conf.Repository.PersonalAccessToken().ReadPersonalAccessToken(expPAT.UserID, expPAT.ID)

		assert.Nil(t, err, "err is nil")
		assert.True(t, pat.Revoked, "pat is revoked")

		return nil
	}, testutils.InitUsers, testutils.InitPATs)
}

func TestDeletePAT(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expPAT := testutils.PATModels[0]

		pat, err := conf.Repository.PersonalAccessToken().DeletePersonalAccessToken(expPAT)

		if err != nil {
			return err
		}

		assert.Nil(t, err, "err is nil")

		testutils.AssertPATsEqual(t, pat, expPAT)

		// ensure session cannot be read
		_, failingErr := conf.Repository.PersonalAccessToken().ReadPersonalAccessToken(expPAT.UserID, expPAT.ID)

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)

		return nil
	}, testutils.InitUsers, testutils.InitPATs)
}
