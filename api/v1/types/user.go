package types

// swagger:model
type User struct {
	*APIResourceMeta

	// the display name for this user
	// example: User 1
	DisplayName string `json:"display_name"`

	// the email address for this user
	// example: user1@gmail.com
	Email string `json:"email"`

	// whether this user's email address has been verified
	// example: false
	EmailVerified bool `json:"email_verified"`

	// a URI for the user icon
	// example: https://avatars.githubusercontent.com/u/25448214?v=4
	Icon string `json:"icon"`
}

// swagger:model
type GetUserResponse User

// swagger:model
type CreateUserRequest struct {
	// the display name for this user
	//
	// required: true
	// example: User 1
	DisplayName string `json:"display_name" form:"required,max=255"`

	// the email address for this user
	//
	// required: true
	// example: user1@gmail.com
	Email string `json:"email" form:"required,max=255,email"`

	// the password for this user
	//
	// required: true
	// example: Securepassword123
	Password string `json:"password" form:"required,max=255,password"`
}

// swagger:model
type CreateUserResponse User

const InvalidEmailOrPasswordCode uint = 2403

var InvalidEmailOrPassword APIError = APIError{
	Code:        InvalidEmailOrPasswordCode,
	Description: "Invalid email or password combination",
}
