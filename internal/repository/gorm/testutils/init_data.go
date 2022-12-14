package testutils

import (
	"testing"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type InitData struct {
	Users []*models.User
	PATs  []*models.PersonalAccessToken
}

// All models will be populated with data and IDs after init methods are called
var InitDataAll = &InitData{
	Users: UserModels,
	PATs:  PATModels,
}

var UserModels []*models.User = []*models.User{
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

var PATModels []*models.PersonalAccessToken = []*models.PersonalAccessToken{
	{
		DisplayName: "test-pat-1",
	},
}

type InitDataFunc func(t *testing.T, conf *database.Config) error

func InitUsers(t *testing.T, conf *database.Config) error {
	for i, declaredUser := range UserModels {
		userCp := declaredUser

		user, err := conf.Repository.User().CreateUser(userCp)

		if err != nil {
			return err
		}

		UserModels[i] = user
	}

	return nil
}

// Note that the declared PATs are assigned to the users in a round-robin fashion
func InitPATs(t *testing.T, conf *database.Config) error {
	for i, declaredPAT := range PATModels {
		patCp := declaredPAT

		parentUser := UserModels[i%len(UserModels)]

		// call population method
		pat, err := models.NewPATFromUserID(patCp.DisplayName, parentUser.ID)

		if err != nil {
			return err
		}

		pat, err = conf.Repository.PersonalAccessToken().CreatePersonalAccessToken(pat)

		PATModels[i] = pat
	}

	return nil
}
