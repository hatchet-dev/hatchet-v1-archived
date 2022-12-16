package types

const URLParamOrgID URLParam = "org_id"

// swagger:model
type Organization struct {
	*APIResourceMeta

	// the display name for the personal access token
	// example: cli-token-1234
	DisplayName string `json:"display_name"`

	// information about the organization owner
	Owner UserOrgPublishedData `json:"owner"`
}

// swagger:model
type CreateOrganizationRequest struct {
	// the display name for this user
	//
	// required: true
	// example: User 1
	DisplayName string `json:"display_name" form:"required,max=255"`
}

// swagger:model
type CreateOrganizationResponse Organization

// swagger:model
type GetOrganizationResponse Organization
