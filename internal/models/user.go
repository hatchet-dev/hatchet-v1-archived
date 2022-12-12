package models

import (
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"gorm.io/gorm"
)

// User type that extends gorm.Model
type User struct {
	gorm.Model

	DisplayName   string
	Email         string `gorm:"unique"`
	EmailVerified bool
	Password      string
	Icon          string
}

func (u *User) ToAPIType() *types.User {
	return &types.User{
		DisplayName:   u.DisplayName,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Icon:          u.Icon,
	}
}
