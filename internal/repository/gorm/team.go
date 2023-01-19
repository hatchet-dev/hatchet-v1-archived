package gorm

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
)

// TeamRepository uses gorm.DB for querying the database
type TeamRepository struct {
	db *gorm.DB
}

// NewTeamRepository returns a TeamRepository which uses
// gorm.DB for querying the database
func NewTeamRepository(db *gorm.DB) repository.TeamRepository {
	return &TeamRepository{db}
}

// CreateOrg adds a new Org row to the Orgs table in the database
func (repo *TeamRepository) CreateTeam(team *models.Team) (*models.Team, repository.RepositoryError) {
	if err := repo.db.Create(team).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return team, nil
}

// ReadTeamByID finds a single team by its unique id
func (repo *TeamRepository) ReadTeamByID(id string) (*models.Team, repository.RepositoryError) {
	team := &models.Team{}

	if err := repo.db.Where("id = ?", id).First(&team).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return team, nil
}

// UpdateTeam updates a team in the database
func (repo *TeamRepository) UpdateTeam(team *models.Team) (*models.Team, repository.RepositoryError) {
	if err := repo.db.Save(team).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return team, nil
}

// DeleteTeam soft-deletes a single team by its unique id
func (repo *TeamRepository) DeleteTeam(team *models.Team) (*models.Team, repository.RepositoryError) {
	del := repo.db.Delete(&team)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return team, nil
}

func (repo *TeamRepository) ListTeamsByOrgID(orgID string, opts ...repository.QueryOption) ([]*models.Team, *repository.PaginatedResult, repository.RepositoryError) {
	var teams []*models.Team

	db := repo.db.Model(&models.Team{}).Where("organization_id = ?", orgID)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&teams).Error; err != nil {
		return nil, nil, err
	}

	return teams, paginatedResult, nil
}

func (repo *TeamRepository) CreateTeamMember(team *models.Team, teamMember *models.TeamMember) (*models.TeamMember, repository.RepositoryError) {
	if err := repo.db.Model(team).Association("TeamMembers").Append(teamMember); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return teamMember, nil
}

func (repo *TeamRepository) ReadTeamMemberByID(teamID, memberID string) (*models.TeamMember, repository.RepositoryError) {
	member := &models.TeamMember{}

	if err := repo.db.Preload("TeamPolicies").Joins("OrgMember").Where("team_members.team_id = ? AND team_members.id = ?", teamID, memberID).First(&member).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return member, nil
}

func (repo *TeamRepository) ReadTeamMemberByOrgMemberID(teamID, orgMemberID string) (*models.TeamMember, repository.RepositoryError) {
	member := &models.TeamMember{}

	if err := repo.db.Preload("TeamPolicies").Joins("OrgMember").Where("team_members.team_id = ? AND team_members.org_member_id = ?", teamID, orgMemberID).First(&member).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return member, nil
}

func (repo *TeamRepository) ListTeamMembersByTeamID(teamID string, opts ...repository.QueryOption) ([]*models.TeamMember, *repository.PaginatedResult, repository.RepositoryError) {
	var members []*models.TeamMember

	db := repo.db.Model(&models.TeamMember{}).Where("team_id = ?", teamID)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Preload("OrgMember").Preload("TeamPolicies").Find(&members).Error; err != nil {
		return nil, nil, err
	}

	return members, paginatedResult, nil
}

func (repo *TeamRepository) UpdateTeamMember(teamMember *models.TeamMember) (*models.TeamMember, repository.RepositoryError) {
	if err := repo.db.Save(teamMember).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return teamMember, nil
}

func (repo *TeamRepository) DeleteTeamMember(teamMember *models.TeamMember) (*models.TeamMember, repository.RepositoryError) {
	del := repo.db.Delete(&teamMember)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return teamMember, nil
}

func (repo *TeamRepository) AppendTeamPolicyToTeamMember(teamMember *models.TeamMember, teamPolicy *models.TeamPolicy) (*models.TeamMember, repository.RepositoryError) {
	// we add an additional check to verify that the team ids are the same
	if teamMember.TeamID != teamPolicy.TeamID {
		return nil, repository.UnknownRepositoryError(
			fmt.Errorf("team ids are not equal: member has team %s, policy has team %s", teamMember.TeamID, teamPolicy.TeamID),
		)
	}

	if err := repo.db.Model(teamMember).Association("TeamPolicies").Append(teamPolicy); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return teamMember, nil
}

func (repo *TeamRepository) ReplaceTeamPoliciesForTeamMember(teamMember *models.TeamMember, policies []*models.TeamPolicy) (*models.TeamMember, repository.RepositoryError) {
	if err := repo.db.Model(teamMember).Association("TeamPolicies").Clear(); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	if err := repo.db.Model(teamMember).Association("TeamPolicies").Append(policies); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return teamMember, nil
}

func (repo *TeamRepository) RemoveTeamPolicyFromTeamMember(teamMember *models.TeamMember, teamPolicy *models.TeamPolicy) (*models.TeamMember, repository.RepositoryError) {
	// we add an additional check to verify that the team ids are the same
	if teamMember.TeamID != teamPolicy.TeamID {
		return nil, repository.UnknownRepositoryError(
			fmt.Errorf("team ids are not equal: member has team %s, policy has team %s", teamMember.TeamID, teamPolicy.TeamID),
		)
	}

	if err := repo.db.Model(teamMember).Association("TeamPolicies").Delete(teamPolicy); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return teamMember, nil
}

func (repo *TeamRepository) CreateTeamPolicy(team *models.Team, teamPolicy *models.TeamPolicy) (*models.TeamPolicy, repository.RepositoryError) {
	if err := repo.db.Model(team).Association("TeamPolicies").Append(teamPolicy); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return teamPolicy, nil
}

func (repo *TeamRepository) ReadPresetTeamPolicyByName(orgID string, presetName models.PresetTeamPolicyName) (*models.TeamPolicy, repository.RepositoryError) {
	policy := &models.TeamPolicy{}

	if err := repo.db.Where("team_id = ? AND is_custom = ? AND policy_name = ?", orgID, false, presetName).First(&policy).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return policy, nil
}

func (repo *TeamRepository) ReadPolicyByID(orgID, policyID string) (*models.TeamPolicy, repository.RepositoryError) {
	policy := &models.TeamPolicy{}

	if err := repo.db.Where("team_id = ? AND id = ?", orgID, policyID).First(&policy).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return policy, nil
}

func (repo *TeamRepository) ListTeamPoliciesByTeamID(orgID string, opts ...repository.QueryOption) ([]*models.TeamPolicy, *repository.PaginatedResult, repository.RepositoryError) {
	var policies []*models.TeamPolicy

	db := repo.db.Model(&models.TeamPolicy{}).Where("team_id = ?", orgID)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&policies).Error; err != nil {
		return nil, nil, err
	}

	return policies, paginatedResult, nil
}
