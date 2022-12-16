package testutils

import (
	"testing"
	"time"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
)

type InitData struct {
	Users        []*models.User
	UserSessions []*models.UserSession
	PATs         []*models.PersonalAccessToken
	Orgs         []*models.Organization
	OrgMembers   []*models.OrganizationMember
}

// All models will be populated with data and IDs after init methods are called
var InitDataAll = &InitData{
	Users:        UserModels,
	UserSessions: UserSessionModels,
	PATs:         PATModels,
	Orgs:         OrgModels,
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

// UserSessionModels do not represent actually valid user sessions
var UserSessionModels []*models.UserSession = []*models.UserSession{
	{
		Key:       "1",
		Data:      []byte("1"),
		ExpiresAt: time.Now().Add(24 * 30 * time.Hour),
	},
}

var OrgModels []*models.Organization = []*models.Organization{
	{
		DisplayName: "My Org 1",
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

		if err != nil {
			return err
		}

		PATModels[i] = pat
	}

	return nil
}

func InitUserSessions(t *testing.T, conf *database.Config) error {
	for i, declaredUserSession := range UserSessionModels {
		sessCp := declaredUserSession

		userSession, err := conf.Repository.UserSession().CreateUserSession(sessCp)

		if err != nil {
			return err
		}

		UserSessionModels[i] = userSession
	}

	return nil
}

// Note that the declared orgs are assigned to the users in a round-robin fashion
func InitOrgs(t *testing.T, conf *database.Config) error {
	for i, declaredOrg := range OrgModels {
		orgCp := declaredOrg

		parentUser := UserModels[i%len(UserModels)]

		orgCp.OwnerID = parentUser.ID

		orgCp.OrgMembers = []models.OrganizationMember{
			{
				InviteAccepted: true,
				UserID:         parentUser.ID,
			},
		}

		// call population method
		org, err := conf.Repository.Org().CreateOrg(orgCp)

		if err != nil {
			return err
		}

		OrgModels[i] = org
	}

	return nil
}
