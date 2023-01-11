package models

import (
	"time"

	"github.com/hatchet-dev/hatchet/internal/encryption"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// VerifyEmailToken represents an email verification request for a user. The token is hashed and stored
// in the database, with the raw (unhashed) token sent to the user. The token is discovered via token
// ID and compared to its hashed value.
type VerifyEmailToken struct {
	Base

	// Email is the target email of the user that initiated the password reset request
	Email string

	// Revoked represents whether the token has been revoked (used) or not
	Revoked bool

	// Expiry time
	Expires *time.Time

	// Token is hashed before storage
	Token []byte
}

func NewVerifyEmailTokenFromEmail(targetEmail string) (*VerifyEmailToken, error) {
	p := &VerifyEmailToken{
		Email: targetEmail,
	}

	// we set the default expiry of PAT's to be 1 hour
	expires := time.Now().Add(24 * time.Hour)

	p.Expires = &expires

	rawToken, err := encryption.GenerateRandomBytes(32)

	if err != nil {
		return nil, err
	}

	p.Token = []byte(rawToken)

	return p, err
}

func (p *VerifyEmailToken) BeforeCreate(tx *gorm.DB) error {
	err := p.Base.BeforeCreate(tx)

	if err != nil {
		return err
	}

	// hash the password before create using bcrypt
	hashedTok, err := bcrypt.GenerateFromPassword([]byte(p.Token), 8)

	if err != nil {
		return err
	}

	p.Token = hashedTok

	return nil
}

func (p *VerifyEmailToken) IsExpired() bool {
	timeLeft := p.Expires.Sub(time.Now())
	return timeLeft < 0
}

func (p *VerifyEmailToken) VerifyToken(tok string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Token, []byte(tok))

	return err == nil, err
}
