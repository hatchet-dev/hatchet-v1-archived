package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/encryption"
)

// PersonalAccessToken contains additional data about the JWT token, and provides a mechanism
// for revoking an otherwise-valid JWT.
type PersonalAccessToken struct {
	Base

	// How this token gets displayed to the user
	DisplayName string

	// When this PAT expires. This should match what's in the JWT data
	Expires *time.Time

	// Whether the personal access token has been revoked
	Revoked bool

	// Encrypted data that contains the token signing secret for that specific token
	SigningSecret []byte

	// The user id that this PAT has been written for.
	UserID string
}

func (p *PersonalAccessToken) ToAPIType() *types.PersonalAccessToken {
	return &types.PersonalAccessToken{
		APIResourceMeta: p.Base.ToAPITypeMetadata(),
		DisplayName:     p.DisplayName,
		Expires:         p.Expires,
		Revoked:         p.Revoked,
	}
}

func NewPATFromUserID(displayName, userID string) (*PersonalAccessToken, error) {
	pat := &PersonalAccessToken{
		DisplayName: displayName,
		UserID:      userID,
	}

	// in this case, we generate the UUID ahead of time (rather than BeforeCreate), as the token's UUID
	// is needed by the `token` package to generate the JWT.
	pat.Base.ID = uuid.New().String()

	// we set the default expiry of PAT's to be 30 days
	expires := time.Now().Add(30 * 24 * time.Hour)

	pat.Expires = &expires

	secretData, err := encryption.GenerateRandomBytes(32)

	if err != nil {
		return nil, err
	}

	pat.SigningSecret = []byte(secretData)

	return pat, err
}

func (p *PersonalAccessToken) IsExpired() bool {
	timeLeft := p.Expires.Sub(time.Now())
	return timeLeft < 0
}

func (p *PersonalAccessToken) Encrypt(key *[32]byte) error {
	ciphertext, err := encryption.Encrypt(p.SigningSecret, key)

	if err != nil {
		return err
	}

	p.SigningSecret = ciphertext

	return nil
}

func (p *PersonalAccessToken) Decrypt(key *[32]byte) error {
	plaintext, err := encryption.Decrypt(p.SigningSecret, key)

	if err != nil {
		return err
	}

	p.SigningSecret = plaintext

	return nil
}
