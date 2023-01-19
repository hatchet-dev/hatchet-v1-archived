package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// TeamRepository represents the set of queries on the Team model
type TeamRepository interface {
	// --- Team queries ---
	//
	// CreateTeam creates a new team in the database
	CreateTeam(team *models.Team) (*models.Team, RepositoryError)

	// ReadTeamByID reads the team by it's unique UUID
	ReadTeamByID(id string) (*models.Team, RepositoryError)

	// ListTeamsByOrgID lists all teams for an organization
	ListTeamsByOrgID(orgID string, opts ...QueryOption) ([]*models.Team, *PaginatedResult, RepositoryError)

	// UpdateTeam updates any modified values for a team
	UpdateTeam(team *models.Team) (*models.Team, RepositoryError)

	// DeleteTeam soft-deletes a team
	DeleteTeam(team *models.Team) (*models.Team, RepositoryError)

	// --- Team member queries ---
	//
	// CreateTeamMember creates a new team member for that team
	CreateTeamMember(team *models.Team, teamMember *models.TeamMember) (*models.TeamMember, RepositoryError)

	// ReadTeamMemberByUserID finds a team member by their user UUID. Not to be confused with ReadTeamMemberByID,
	// which finds an team member by the **team member ID**.
	ReadTeamMemberByOrgMemberID(teamID, orgMemberID string) (*models.TeamMember, RepositoryError)

	// ReadTeamMemberByID finds a team member by their unique team member UUID.
	ReadTeamMemberByID(teamID, memberID string) (*models.TeamMember, RepositoryError)

	// ListTeamMembersByTeamID lists team members that are part of that team.
	ListTeamMembersByTeamID(teamID string, opts ...QueryOption) ([]*models.TeamMember, *PaginatedResult, RepositoryError)

	// UpdateTeamMember updates team members. This MAY have the side effect of updating dependent models,
	// depending on the implementation. Gorm is inconsistent about this so make sure any update methods
	// are tested.
	UpdateTeamMember(teamMember *models.TeamMember) (*models.TeamMember, RepositoryError)

	// DeleteTeamMember deletes a team member.
	DeleteTeamMember(teamMember *models.TeamMember) (*models.TeamMember, RepositoryError)

	// AppendTeamPolicyToTeamMember adds a team policy to that member.
	AppendTeamPolicyToTeamMember(teamMember *models.TeamMember, teamPolicy *models.TeamPolicy) (*models.TeamMember, RepositoryError)

	// ReplaceTeamPoliciesForTeamMember replaces all team policies for that team member.
	ReplaceTeamPoliciesForTeamMember(teamMember *models.TeamMember, policies []*models.TeamPolicy) (*models.TeamMember, RepositoryError)

	// RemoveTeamPolicyFromTeamMember removes a team policy for that member.
	RemoveTeamPolicyFromTeamMember(teamMember *models.TeamMember, teamPolicy *models.TeamPolicy) (*models.TeamMember, RepositoryError)

	// --- Team policy queries ---
	//
	// CreateTeamPolicy creates a new team policy
	CreateTeamPolicy(team *models.Team, teamPolicy *models.TeamPolicy) (*models.TeamPolicy, RepositoryError)

	// ReadPresetTeamPolicyByName finds a preset policy for a team by its PolicyName
	ReadPresetTeamPolicyByName(teamID string, presetName models.PresetTeamPolicyName) (*models.TeamPolicy, RepositoryError)

	// ReadPolicyByID finds a policy by its ID
	ReadPolicyByID(teamID, policyID string) (*models.TeamPolicy, RepositoryError)

	// ListTeamPoliciesByTeamID lists policies for the team
	ListTeamPoliciesByTeamID(teamID string, opts ...QueryOption) ([]*models.TeamPolicy, *PaginatedResult, RepositoryError)
}
