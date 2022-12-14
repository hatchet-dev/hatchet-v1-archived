package types

import "time"

const PersonalAccessTokenURLParam URLParam = "pat_id"

// swagger:model
type PersonalAccessToken struct {
	*APIResourceMeta

	// the display name for the personal access token
	// example: cli-token-1234
	DisplayName string `json:"display_name"`

	// when the token expires
	// example: 2023-01-12T17:09:28.35059-05:00
	Expires *time.Time `json:"expires"`

	// whether the token has been revoked
	// example: false
	Revoked bool `json:"revoked"`
}

// swagger:model
type GetPATResponse PersonalAccessToken

// swagger:model
type CreatePATRequest struct {
	// the display name for the personal access token
	// required: true
	// example: cli-token-1234
	DisplayName string `json:"display_name" form:"required,max=255"`
}

// swagger:model
type CreatePATResponse struct {
	// the personal access token object
	PersonalAccessToken PersonalAccessToken `json:"pat"`

	// the raw JWT token. see API documentation for details
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9....
	Token string `json:"token"`
}

// swagger:parameters listPersonalAccessTokens
type ListPATsRequest struct {
	*PaginationRequest
}

// swagger:model
type ListPATsResponse struct {
	Pagination *PaginationResponse    `json:"pagination"`
	Rows       []*PersonalAccessToken `json:"rows"`
}

// swagger:model
type RevokePATResponse PersonalAccessToken

// swagger:model
type RevokePATResponseExample struct {
	*PersonalAccessToken

	// whether the token is revoked
	// example: true
	Revoked bool `json:"revoked"`
}

// swagger:model
type DeletePATResponse PersonalAccessToken
