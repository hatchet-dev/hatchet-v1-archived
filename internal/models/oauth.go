package models

import (
	"time"

	"github.com/hatchet-dev/hatchet/internal/encryption"
)

// SharedOAuthFields stores general fields needed for an oauth integration
type SharedOAuthFields struct {
	HasEncryptedFields

	// The ID issued to the client
	ClientID []byte `json:"client-id"`

	// The end-users's access token
	AccessToken []byte `json:"access-token"`

	// The end-user's refresh token
	RefreshToken []byte `json:"refresh-token"`

	// Time token expires and needs to be refreshed.
	// If 0, token will never refresh
	Expiry time.Time

	// The id of the user that linked oauth
	UserID string
	User   User `gorm:"foreignKey:UserID"`
}

func (s *SharedOAuthFields) Encrypt(key *[32]byte) error {
	if !s.HasEncryptedFields.FieldsAreEncrypted {
		ciphertext, err := encryption.Encrypt(s.AccessToken, key)

		if err != nil {
			return err
		}

		s.AccessToken = ciphertext

		ciphertext, err = encryption.Encrypt(s.RefreshToken, key)

		if err != nil {
			return err
		}

		s.RefreshToken = ciphertext

		s.HasEncryptedFields.FieldsAreEncrypted = true
	}

	return nil
}

func (s *SharedOAuthFields) Decrypt(key *[32]byte) error {
	if s.HasEncryptedFields.FieldsAreEncrypted {
		plaintext, err := encryption.Decrypt(s.AccessToken, key)

		if err != nil {
			return err
		}

		s.AccessToken = plaintext

		plaintext, err = encryption.Decrypt(s.RefreshToken, key)

		if err != nil {
			return err
		}

		s.RefreshToken = plaintext

		s.HasEncryptedFields.FieldsAreEncrypted = false
	}

	return nil
}
