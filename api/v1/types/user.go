package types

// swagger:model
type User struct {
	*APIResourceMeta

	DisplayName   string `json:"display_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Icon          string `json:"icon"`
}

// swagger:model
type GetUserResponse User

// swagger:model
type CreateUserRequest struct {
	DisplayName string `json:"display_name" form:"required,max=255"`
	Email       string `json:"email" form:"required,max=255,email"`
	Password    string `json:"password" form:"required,max=255,password"`
}

// swagger:model
type CreateUserResponse User

const InvalidEmailOrPasswordCode APIErrorCode = 2403

var InvalidEmailOrPassword APIError = APIError{
	Code:        InvalidEmailOrPasswordCode,
	Description: "Invalid email or password combination",
}
