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

func TestCreateOrg(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		org, err := conf.Repository.Org().CreateOrg(&models.Organization{
			DisplayName: "Org 1",
			Icon:        "https://icon.example.com",
			OwnerID:     testutils.UserModels[0].ID,
			OrgMembers: []models.OrganizationMember{
				{
					InviteAccepted: true,
					UserID:         testutils.UserModels[0].ID,
				},
			},
		})

		if err != nil {
			t.Fatalf("could not create org: %v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(org.ID))
		assert.Equal(t, "Org 1", org.DisplayName)

		// verify that the preset policies were assigned IDs
		for _, policy := range org.OrgPolicies {
			assert.True(t, uuidutils.IsValidUUID(policy.ID))
		}

		assert.Equal(t, 3, len(org.OrgPolicies), "length of org members should be 3")

		// verify that the org members were assigned IDs
		for _, member := range org.OrgMembers {
			assert.True(t, uuidutils.IsValidUUID(member.ID))
		}

		assert.Equal(t, 1, len(org.OrgMembers), "length of org members should be 1")

		// verify that the org member is the owner
		assert.Equal(t, testutils.UserModels[0].ID, org.OrgMembers[0].UserID, "org member is user id")

		return nil
	}, testutils.InitUsers)
}

func TestReadOrgByID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		org, err := conf.Repository.Org().ReadOrgByID(testutils.OrgModels[0].ID)

		if err != nil {
			t.Fatalf("could not read org: %v", err)
		}

		assert.Equal(t, org.Owner.ID, org.OwnerID, "owner is populated with correct id")
		assert.Equal(t, org.Owner.ID, testutils.UserModels[0].ID, "owner is set to the first user")

		testutils.AssertOrgsEqual(t, testutils.OrgModels[0], org, false, false)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestFailedReadOrgByID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		_, failingErr := conf.Repository.Org().ReadOrgByID("not-an-id")

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestUpdateOrg(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expOrg := testutils.OrgModels[0]

		expOrg.DisplayName = "Org-Rename 1"

		org, err := conf.Repository.Org().UpdateOrg(expOrg)

		if err != nil {
			t.Fatalf("could not update org: %v", err)
		}

		// ensure that display name was written
		org, err = conf.Repository.Org().ReadOrgByID(testutils.OrgModels[0].ID)

		if err != nil {
			t.Fatalf("could not read org: %v", err)
		}

		testutils.AssertOrgsEqual(t, expOrg, org, false, false)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestDeleteOrg(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expOrg := testutils.OrgModels[0]

		org, err := conf.Repository.Org().DeleteOrg(expOrg)

		if err != nil {
			t.Fatalf("could not delete org: %v", err)
		}

		// ensure that display name was written
		org, failingErr := conf.Repository.Org().ReadOrgByID(testutils.OrgModels[0].ID)

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)
		assert.Nil(t, org, "org should be nil")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestFailedDeleteOrg(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		notAOrg := &models.Organization{
			Base: models.Base{
				ID: "not-an-id",
			},
		}

		_, failingErr := conf.Repository.Org().DeleteOrg(notAOrg)

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryNoRowsAffected, failingErr)

		return nil
	}, testutils.InitOrgs)
}

func TestCreateOrgMember(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgMember := &models.OrganizationMember{
			UserID: testutils.UserModels[1].ID,
		}

		orgMember, err := conf.Repository.Org().CreateOrgMember(testutils.OrgModels[0], orgMember)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(orgMember.ID), "org member should have a valid id")
		assert.Equal(t, testutils.UserModels[1].ID, orgMember.UserID, "user ids should be equal")
		assert.Equal(t, testutils.OrgModels[0].ID, orgMember.OrganizationID, "org ids should be equal")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestReadOrgMemberByInviteeEmail(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgMember, err := conf.Repository.Org().ReadOrgMemberByUserOrInviteeEmail(testutils.OrgModels[0].ID, testutils.OrgInviteLinks[0].InviteeEmail)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(orgMember.ID), "org member should have a valid id")

		// invite link should be populated with correct invitee email
		assert.Equal(t, testutils.OrgInviteLinks[0].InviteeEmail, orgMember.InviteLink.InviteeEmail, "emails should be equal")
		assert.False(t, orgMember.InviteLink.Used, "invite should not be used")

		// assert that this is admin user
		adminPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(testutils.OrgModels[0].ID, models.PresetPolicyNameAdmin)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, adminPolicy.ID, orgMember.OrgPolicies[0].ID, "org policy should be admin")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgInviteLinks)
}

func TestReadOrgMemberByUserEmail(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgMember, err := conf.Repository.Org().ReadOrgMemberByUserOrInviteeEmail(testutils.OrgModels[0].ID, testutils.OrgAdditionalMembers[0].User.Email)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(orgMember.ID), "org member should have a valid id")

		// invite link should be populated with correct invitee email
		assert.Equal(t, testutils.OrgAdditionalMembers[0].User.Email, orgMember.User.Email, "emails should be equal")

		// assert that this is admin user
		adminPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(testutils.OrgModels[0].ID, models.PresetPolicyNameAdmin)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, adminPolicy.ID, orgMember.OrgPolicies[0].ID, "org policy should be admin")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}

func TestReadOrgMemberByID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgMember, err := conf.Repository.Org().ReadOrgMemberByID(testutils.OrgModels[0].ID, testutils.OrgAdditionalMembers[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(orgMember.ID), "org member should have a valid id")

		// invite link should be populated with correct invitee email
		assert.Equal(t, testutils.OrgAdditionalMembers[0].User.Email, orgMember.User.Email, "emails should be equal")

		// assert that this is admin user
		adminPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(testutils.OrgModels[0].ID, models.PresetPolicyNameAdmin)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, adminPolicy.ID, orgMember.OrgPolicies[0].ID, "org policy should be admin")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}

func TestReadOrgMemberByUserID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgMember, err := conf.Repository.Org().ReadOrgMemberByUserID(testutils.OrgModels[0].ID, testutils.UserModels[2].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(orgMember.ID), "org member should have a valid id")

		// invite link should be populated with correct invitee email
		assert.Equal(t, testutils.OrgAdditionalMembers[0].User.Email, orgMember.User.Email, "emails should be equal")

		// assert that this is admin user
		adminPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(testutils.OrgModels[0].ID, models.PresetPolicyNameAdmin)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, adminPolicy.ID, orgMember.OrgPolicies[0].ID, "org policy should be admin")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}

func TestListOrgMembersByOrgID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgMembers, pagination, err := conf.Repository.Org().ListOrgMembersByOrgID(testutils.OrgModels[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, &repository.PaginatedResult{
			NumPages:    1,
			CurrentPage: 0,
			NextPage:    0,
		}, pagination, "pagination should be equal")

		assert.Equal(t, 2, len(orgMembers), "length of org members should be 1")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}

func TestUpdateOrgMember(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgMember, err := conf.Repository.Org().ReadOrgMemberByID(testutils.OrgModels[0].ID, testutils.OrgAdditionalMembers[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		orgMember.InviteAccepted = false

		orgMember, err = conf.Repository.Org().UpdateOrgMember(orgMember)

		if err != nil {
			t.Fatalf("%v", err)
		}

		// read the member again to ensure successful write
		orgMember, err = conf.Repository.Org().ReadOrgMemberByID(testutils.OrgModels[0].ID, testutils.OrgAdditionalMembers[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.False(t, orgMember.InviteAccepted)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}

func TestReadOrgInviteByID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgInvite, err := conf.Repository.Org().ReadOrgInviteByID(testutils.OrgInviteLinks[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, testutils.OrgInviteLinks[0].ID, orgInvite.ID)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgInviteLinks)
}

func TestUpdateOrgInvite(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgInvite, err := conf.Repository.Org().ReadOrgInviteByID(testutils.OrgInviteLinks[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		orgInvite.Used = true

		orgInvite, err = conf.Repository.Org().UpdateOrgInvite(orgInvite)

		if err != nil {
			t.Fatalf("%v", err)
		}

		// read invite again to ensure successful write
		orgInvite, err = conf.Repository.Org().ReadOrgInviteByID(testutils.OrgInviteLinks[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, testutils.OrgInviteLinks[0].ID, orgInvite.ID)
		assert.True(t, orgInvite.Used)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgInviteLinks)
}

func TestDeleteOrgMember(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgMember, err := conf.Repository.Org().ReadOrgMemberByID(testutils.OrgModels[0].ID, testutils.OrgAdditionalMembers[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		orgMember, err = conf.Repository.Org().DeleteOrgMember(orgMember)

		if err != nil {
			t.Fatalf("%v", err)
		}

		// read the member again to ensure successful write
		orgMember, err = conf.Repository.Org().ReadOrgMemberByID(testutils.OrgModels[0].ID, testutils.OrgAdditionalMembers[0].ID)

		assert.ErrorIs(t, err, repository.RepositoryErrorNotFound)
		assert.Nil(t, orgMember)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}

func TestAppendOrgPolicyToOrgMember(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		// get the admin policy
		memberPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(testutils.OrgModels[0].ID, models.PresetPolicyNameMember)

		if err != nil {
			t.Fatalf("%v", err)
		}

		// get the org member
		orgMember, err := conf.Repository.Org().ReadOrgMemberByUserID(testutils.OrgModels[0].ID, testutils.UserModels[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		orgMember, err = conf.Repository.Org().AppendOrgPolicyToOrgMember(orgMember, memberPolicy)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, 2, len(orgMember.OrgPolicies), "org policy length on members should be 1")

		var gotMemberPolicy *models.OrganizationPolicy
		var gotOwnerPolicy *models.OrganizationPolicy

		for _, policy := range orgMember.OrgPolicies {
			policyCp := policy

			if policy.PolicyName == string(models.PresetPolicyNameOwner) {
				gotOwnerPolicy = &policyCp
			} else if policy.PolicyName == string(models.PresetPolicyNameMember) {
				gotMemberPolicy = &policyCp
			}
		}

		assert.NotNil(t, gotMemberPolicy)
		assert.Equal(t, memberPolicy.ID, gotMemberPolicy.ID, "org policy IDs for member are equal")

		assert.NotNil(t, gotOwnerPolicy)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestReplaceOrgPoliciesForOrgMember(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		// get the admin policy
		memberPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(testutils.OrgModels[0].ID, models.PresetPolicyNameMember)

		if err != nil {
			t.Fatalf("%v", err)
		}

		// get the org member
		orgMember, err := conf.Repository.Org().ReadOrgMemberByUserID(testutils.OrgModels[0].ID, testutils.UserModels[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		orgMember, err = conf.Repository.Org().ReplaceOrgPoliciesForOrgMember(orgMember, []*models.OrganizationPolicy{memberPolicy})

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, 1, len(orgMember.OrgPolicies), "org policy length on members should be 1")
		assert.Equal(t, memberPolicy.ID, orgMember.OrgPolicies[0].ID, "org policy IDs are equal")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestRemoveOrgPolicyFromOrgMember(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		ownerPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(testutils.OrgModels[0].ID, models.PresetPolicyNameOwner)

		if err != nil {
			t.Fatalf("%v", err)
		}

		// get the org member
		orgMember, err := conf.Repository.Org().ReadOrgMemberByUserID(testutils.OrgModels[0].ID, testutils.UserModels[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		orgMember, err = conf.Repository.Org().RemoveOrgPolicyFromOrgMember(orgMember, ownerPolicy)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, 0, len(orgMember.OrgPolicies), "org policy length on members should be 1")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestCreateOrgPolicy(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgPolicy := &models.OrganizationPolicy{
			IsCustom:   true,
			PolicyName: "admin-policy-copy",
			Policy: []byte(`
			package hatchet_org_presets.admin

			import future.keywords.contains
			import future.keywords.every
			import future.keywords.if
			import future.keywords.in
			
			allow if {
				resource := input.endpoint.resources[_]
			
				# admins cannot perform any actions on the org owner resource
				resource.scope != "org_owner_scope"
			}
			`),
		}

		orgPolicy, err := conf.Repository.Org().CreateOrgPolicy(testutils.OrgModels[0], orgPolicy)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(orgPolicy.ID), "has a valid uuid assigned")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestReadPresetPolicyByName(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgPolicy, err := conf.Repository.Org().ReadPresetPolicyByName(testutils.OrgModels[0].ID, models.PresetPolicyNameAdmin)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(orgPolicy.ID), "has a valid uuid assigned")
		assert.Equal(t, string(models.PresetPolicyNameAdmin), orgPolicy.PolicyName, "correct name")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestReadPresetPolicyByID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgPolicy, err := conf.Repository.Org().ReadPolicyByID(testutils.OrgModels[0].ID, testutils.OrgModels[0].OrgPolicies[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(orgPolicy.ID), "has a valid uuid assigned")
		assert.Equal(t, testutils.OrgModels[0].OrgPolicies[0].PolicyName, orgPolicy.PolicyName, "correct name")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestListOrgPoliciesByOrgID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgPolicies, pagination, err := conf.Repository.Org().ListOrgPoliciesByOrgID(testutils.OrgModels[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, &repository.PaginatedResult{
			NumPages:    1,
			CurrentPage: 0,
			NextPage:    0,
		}, pagination, "pagination should be equal")

		assert.Equal(t, 3, len(orgPolicies), "length of org members should be 3")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}
