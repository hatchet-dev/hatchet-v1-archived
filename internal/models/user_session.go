package models

import "time"

type UserSession struct {
	Base

	// Key contains the session id
	Key string `gorm:"unique"`

	// Contains the encrypted cookie data
	Data []byte

	// Time the session will expire
	ExpiresAt time.Time
}
