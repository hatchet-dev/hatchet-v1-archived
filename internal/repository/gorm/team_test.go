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

func TestCreateTeam(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		orgMemberOwner, err := conf.Repository.Org().ReadOrgMemberByUserID(testutils.OrgModels[0].ID, testutils.OrgModels[0].OwnerID, false)

		if err != nil {
			t.Fatalf("%v", err)
		}

		team, err := conf.Repository.Team().CreateTeam(&models.Team{
			DisplayName:    "Team 1",
			OrganizationID: testutils.OrgModels[0].ID,
			TeamMembers: []models.TeamMember{
				{
					OrgMemberID: orgMemberOwner.ID,
				},
			},
		})

		if err != nil {
			t.Fatalf("could not create org: %v", err)
		}

		assert.True(t, uuidutils.IsValidUUID(team.ID))
		assert.Equal(t, "Team 1", team.DisplayName)

		// verify that the preset policies were assigned IDs
		for _, policy := range team.TeamPolicies {
			assert.True(t, uuidutils.IsValidUUID(policy.ID))
		}

		assert.Equal(t, 2, len(team.TeamPolicies), "length of team policies should be 2")

		// verify that the org members were assigned IDs
		for _, member := range team.TeamMembers {
			assert.True(t, uuidutils.IsValidUUID(member.ID))
		}

		assert.Equal(t, 1, len(team.TeamMembers), "length of team members should be 1")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestReadTeamByID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		team, err := conf.Repository.Team().ReadTeamByID(testutils.TeamModels[0].ID)

		if err != nil {
			t.Fatalf("could not read org: %v", err)
		}

		assert.Equal(t, testutils.TeamModels[0].DisplayName, team.DisplayName, "display names should be equal")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams)
}

func TestFailedReadTeamByID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		_, failingErr := conf.Repository.Team().ReadTeamByID("not-an-id")

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams)
}

func TestUpdateTeam(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expTeam := testutils.TeamModels[0]

		expTeam.DisplayName = "Team-Rename 1"

		team, err := conf.Repository.Team().UpdateTeam(expTeam)

		if err != nil {
			t.Fatalf("could not update org: %v", err)
		}

		// ensure that display name was written
		team, err = conf.Repository.Team().ReadTeamByID(testutils.TeamModels[0].ID)

		if err != nil {
			t.Fatalf("could not read org: %v", err)
		}

		assert.Equal(t, "Team-Rename 1", team.DisplayName, "display name should be new")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams)
}

func TestDeleteTeam(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		expTeam := testutils.TeamModels[0]

		team, err := conf.Repository.Team().DeleteTeam(expTeam)

		if err != nil {
			t.Fatalf("could not delete team: %v", err)
		}

		// ensure that display name was written
		team, failingErr := conf.Repository.Team().ReadTeamByID(testutils.TeamModels[0].ID)

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)
		assert.Nil(t, team, "team should be nil")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams)
}

func TestFailedDeleteTeam(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		notATeam := &models.Team{
			Base: models.Base{
				ID: "not-an-id",
			},
		}

		_, failingErr := conf.Repository.Team().DeleteTeam(notATeam)

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryNoRowsAffected, failingErr)

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams)
}

func TestListTeamsByOrgID(t *testing.T) {
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		teams, pagination, err := conf.Repository.Team().ListTeamsByOrgID(testutils.OrgModels[0].ID)

		if err != nil {
			t.Fatalf("could not read org: %v", err)
		}

		assert.Equal(t, &repository.PaginatedResult{
			NumPages:    1,
			CurrentPage: 0,
			NextPage:    0,
		}, pagination, "pagination should be equal")

		assert.Equal(t, testutils.TeamModels[0].DisplayName, teams[0].DisplayName, "display names should be equal")
		assert.Equal(t, 1, len(teams), "length of teams should be 1")

		return nil
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitTeams)
}
