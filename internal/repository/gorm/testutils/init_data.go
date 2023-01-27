package testutils

import (
	"testing"
	"time"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/models"
)

var declaredUserModels = []*models.User{
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
	{
		Email:         "invite-0@example.com",
		Password:      "Abcdefgh123",
		DisplayName:   "Invite Test User",
		EmailVerified: false,
	},
}

var UserModels []*models.User = declaredUserModels

var declaredPATModels = []*models.PersonalAccessToken{
	{
		DisplayName: "test-pat-1",
	},
}

var PATModels []*models.PersonalAccessToken = declaredPATModels

var declaredUserSessionModels = []*models.UserSession{
	{
		Key:       "1",
		Data:      []byte("1"),
		ExpiresAt: time.Now().Add(24 * 30 * time.Hour),
	},
}

// UserSessionModels do not represent actually valid user sessions
var UserSessionModels []*models.UserSession = declaredUserSessionModels

var declaredOrgModels = []*models.Organization{
	{
		DisplayName: "My Org 1",
	},
}

var OrgModels []*models.Organization = declaredOrgModels

var declaredTeamModels = []*models.Team{
	{
		DisplayName: "Team 1",
	},
}

var TeamModels []*models.Team = declaredTeamModels

var declaredOrgAdditionalMembers = []*models.OrganizationMember{
	{
		User: models.User{
			Email: "invite-0@example.com",
		},
	},
}

var OrgAdditionalMembers []*models.OrganizationMember = declaredOrgAdditionalMembers

var declaredOrgInviteLinks = []*models.OrganizationInviteLink{
	{
		InviteeEmail: "invite-0@example.com",
	},
}

var OrgInviteLinks []*models.OrganizationInviteLink = declaredOrgInviteLinks

var OrgInviteLinksUnencryptedTok map[string]string = make(map[string]string)

func InitAll() {
	UserModels = copyVals(declaredUserModels)
	PATModels = copyVals(declaredPATModels)
	UserSessionModels = copyVals(declaredUserSessionModels)
	OrgModels = copyVals(declaredOrgModels)
	OrgAdditionalMembers = copyVals(declaredOrgAdditionalMembers)
	OrgInviteLinks = copyVals(declaredOrgInviteLinks)
	TeamModels = copyVals(declaredTeamModels)
}

func copyVals[
	V models.User |
		models.PersonalAccessToken |
		models.UserSession |
		models.Organization |
		models.OrganizationMember |
		models.OrganizationInviteLink |
		models.Team](oldArr []*V) []*V {
	c := make([]*V, len(oldArr))

	for i, p := range oldArr {
		if p == nil {
			continue
		}

		v := *p
		c[i] = &v
	}

	return c
}

type InitDataFunc func(t *testing.T, conf *database.Config) error

func InitUsers(t *testing.T, conf *database.Config) error {
	for i, declaredUser := range declaredUserModels {
		userCp := *declaredUser

		user, err := conf.Repository.User().CreateUser(&userCp)

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

		// make sure org policy is created for owner
		orgPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(org.ID, models.PresetPolicyNameOwner)

		if err != nil {
			return err
		}

		orgMember := &org.OrgMembers[0]

		orgMember, err = conf.Repository.Org().AppendOrgPolicyToOrgMember(orgMember, orgPolicy)

		if err != nil {
			return err
		}

		OrgModels[i] = org
	}

	return nil
}

// InitOrgInviteLinks generates invite links with the declared invitee email addresses
func InitOrgInviteLinks(t *testing.T, conf *database.Config) error {
	for _, declaredOrg := range OrgModels {
		orgCp := declaredOrg

		for i, declaredOrgInvite := range OrgInviteLinks {
			invite, err := models.NewOrganizationInviteLink("http://test.example.com", orgCp.ID)

			if err != nil {
				return err
			}

			err = invite.Encrypt(conf.GetEncryptionKey())

			if err != nil {
				return err
			}

			OrgInviteLinksUnencryptedTok[invite.ID] = string(invite.Token)

			invite.InviteeEmail = declaredOrgInvite.InviteeEmail

			adminPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(orgCp.ID, models.PresetPolicyNameAdmin)

			if err != nil {
				return err
			}

			orgMember := &models.OrganizationMember{
				InviteLink: *invite,
				OrgPolicies: []models.OrganizationPolicy{
					*adminPolicy,
				},
			}

			orgMember, err = conf.Repository.Org().CreateOrgMember(orgCp, orgMember)

			if err != nil {
				return err
			}

			OrgInviteLinks[i] = &orgMember.InviteLink
		}
	}

	return nil
}

// InitOrgAdditionalMemberAdmin adds additional members to the org matching the email address,
// adding them as new admin members. Assigned in round-robin fashion.
func InitOrgAdditionalMemberAdmin(t *testing.T, conf *database.Config) error {
	for i, declaredOrgAdditionalMember := range OrgAdditionalMembers {
		orgCp := OrgModels[i%len(OrgModels)]
		orgMemberCp := declaredOrgAdditionalMember

		orgMemberCp.InviteAccepted = true

		var targetUser models.User

		for _, user := range UserModels {
			if user.Email == orgMemberCp.User.Email {
				targetUser = *user
				break
			}
		}

		orgMemberCp.User = targetUser

		orgMember, err := conf.Repository.Org().CreateOrgMember(orgCp, orgMemberCp)

		if err != nil {
			return err
		}

		// get admin policy
		orgPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(orgCp.ID, models.PresetPolicyNameAdmin)

		if err != nil {
			return err
		}

		orgMember, err = conf.Repository.Org().AppendOrgPolicyToOrgMember(orgMember, orgPolicy)

		if err != nil {
			return err
		}

		OrgAdditionalMembers[i] = orgMember
	}

	return nil
}

// Note that the declared teams are assigned to the orgs in a round-robin fashion
func InitTeams(t *testing.T, conf *database.Config) error {
	for i, declaredTeam := range TeamModels {
		teamCp := declaredTeam

		parentOrg := OrgModels[i%len(OrgModels)]

		teamCp.OrganizationID = parentOrg.ID

		orgMemberOwner, err := conf.Repository.Org().ReadOrgMemberByUserID(parentOrg.ID, parentOrg.OwnerID, false)

		if err != nil {
			return err
		}

		teamCp.TeamMembers = []models.TeamMember{
			{
				OrgMemberID: orgMemberOwner.ID,
			},
		}

		// call population method
		team, err := conf.Repository.Team().CreateTeam(teamCp)

		if err != nil {
			return err
		}

		TeamModels[i] = team
	}

	return nil
}
