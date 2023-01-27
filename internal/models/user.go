package models

import (
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserAccountKind string

const (
	UserAccountEmail   UserAccountKind = "email"
	UserAccountService UserAccountKind = "serviceaccount"
)

type User struct {
	Base

	DisplayName   string
	Email         string `gorm:"unique"`
	EmailVerified bool
	Password      string
	Icon          string

	UserAccountKind UserAccountKind
}

func (u *User) ToAPIType() *types.User {
	return &types.User{
		APIResourceMeta: u.Base.ToAPITypeMetadata(),
		DisplayName:     u.DisplayName,
		Email:           u.Email,
		EmailVerified:   u.EmailVerified,
		Icon:            u.Icon,
	}
}

func (u *User) ToOrgUserPublishedData() *types.UserOrgPublishedData {
	return &types.UserOrgPublishedData{
		DisplayName: u.DisplayName,
		Email:       u.Email,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	err := u.Base.BeforeCreate(tx)

	if err != nil {
		return err
	}

	// hash the password before create using bcrypt
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)

	if err != nil {
		return err
	}

	u.Password = string(hashedPw)

	return nil
}

func (u *User) HashPassword() error {
	// hash the new password using bcrypt
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)

	if err != nil {
		return err
	}

	u.Password = string(hashedPw)

	return nil
}

func (u *User) VerifyPassword(pw string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))

	return err == nil, err
}
