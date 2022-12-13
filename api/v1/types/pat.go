package types

import "time"

// swagger:model
type PersonalAccessToken struct {
	*APIResourceMeta

	DisplayName string     `json:"display_name"`
	Expires     *time.Time `json:"expires"`
	Revoked     bool       `json:"revoked"`
}

// swagger:model
type GetPATResponse PersonalAccessToken

// swagger:model
type CreatePATRequest struct {
	DisplayName string `json:"display_name" form:"required,max=255"`
}

// swagger:model
type CreatePATResponse struct {
	PersonalAccessToken PersonalAccessToken `json:"pat"`
	Token               string              `json:"token"`
}
