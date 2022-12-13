package models

import (
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User type that extends gorm.Model
type User struct {
	Base

	DisplayName   string
	Email         string `gorm:"unique"`
	EmailVerified bool
	Password      string
	Icon          string
}

func (u *User) ToAPIType() *types.User {
	return &types.User{
		APIResourceMeta: u.ToAPITypeMetadata(),
		DisplayName:     u.DisplayName,
		Email:           u.Email,
		EmailVerified:   u.EmailVerified,
		Icon:            u.Icon,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	err := u.Base.BeforeCreate(tx)

	if err != nil {
		return err
	}

	// hash the password before create using bcrypt
	// hash the password using bcrypt
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)

	if err != nil {
		return err
	}

	u.Password = string(hashedPw)

	return nil
}
