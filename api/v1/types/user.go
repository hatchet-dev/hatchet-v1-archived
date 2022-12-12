package types

// swagger:model
type User struct {
	*APIResourceMeta

	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Icon        string `json:"icon"`
}

// swagger:model
type GetUserResponse User

// swagger:model
type CreateUserRequest struct {
	Email    string `json:"email" form:"required,max=255,email"`
	Password string `json:"password" form:"required,max=255"`
}

// swagger:model
type CreateUserResponse User
