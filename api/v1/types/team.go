package types

const (
	URLParamTeamID       URLParam = "team_id"
	URLParamTeamMemberID URLParam = "team_member_id"
)

// swagger:model
type Team struct {
	*APIResourceMeta

	// the display name for the team
	// example: Team 1
	DisplayName string `json:"display_name"`
}

// swagger:model
type TeamMember struct {
	*APIResourceMeta

	OrgMember OrganizationMemberSanitized `json:"org_member"`

	TeamPolicies []TeamPolicyMeta `json:"team_policies"`
}

// swagger:model
type TeamPolicyMeta struct {
	*APIResourceMeta

	Name string `json:"name"`
}

// swagger:model
type CreateTeamRequest struct {
	// the display name for the team
	//
	// required: true
	// example: Team 1
	DisplayName string `json:"display_name" form:"required,max=255"`
}

// swagger:model
type CreateTeamResponse Team

// swagger:parameters listTeams
type ListTeamsRequest struct {
	*PaginationRequest
}

// swagger:model
type ListTeamsResponse struct {
	Pagination *PaginationResponse `json:"pagination"`
	Rows       []*Team             `json:"rows"`
}

// swagger:parameters listTeamMembers
type ListTeamMembersRequest struct {
	*PaginationRequest
}

// swagger:model
type ListTeamMembersResponse struct {
	Pagination *PaginationResponse `json:"pagination"`
	Rows       []*TeamMember       `json:"rows"`
}

// swagger:model
type TeamAddMemberRequest struct {
	// the organization member id of the new team member
	OrgMemberID string `json:"org_member_id" form:"required"`

	// the set of policies for this user
	// required: true
	Policies []TeamPolicyReference `json:"policies" form:"required,min=1,dive"`
}

// swagger:model
type TeamAddMemberResponse TeamMember

// swagger:model
type TeamPolicyReference struct {
	Name string `json:"name" form:"omitempty,oneof=admin member"`
	ID   string `json:"id" form:"omitempty,uuid"`
}

// swagger:model
type TeamUpdateRequest struct {
	// the display name for the team
	//
	// required: true
	// example: Team 1
	DisplayName string `json:"display_name" form:"required,max=255"`
}

// swagger:model
type TeamUpdateResponse Team

// swagger:parameters listUserTeams
type ListUserTeamsRequest struct {
	*PaginationRequest

	// the id of the organization to filter by (optional)
	// in: query
	// example: bb214807-246e-43a5-a25d-41761d1cff9e
	OrganizationID string `schema:"organization_id"`
}

// swagger:model
type ListUserTeamsResponse struct {
	Pagination *PaginationResponse `json:"pagination"`
	Rows       []*Team             `json:"rows"`
}

// swagger:model
type DeleteTeamResponse Team
