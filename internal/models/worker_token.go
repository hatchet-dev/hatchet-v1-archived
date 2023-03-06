package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/hatchet-dev/hatchet/internal/encryption"
)

type WorkerToken struct {
	Base
	HasEncryptedFields

	// The subject of the token (team ID)
	TeamID string

	// When this worker token expires. This should match what's in the JWT data
	Expires *time.Time

	// Whether the worker token has been revoked
	Revoked bool

	// Encrypted data that contains the token signing secret for that specific token
	SigningSecret []byte
}

func NewWorkerTokenFromTeamID(teamID string) (*WorkerToken, error) {
	wt := &WorkerToken{
		TeamID: teamID,
	}

	// in this case, we generate the UUID ahead of time (rather than BeforeCreate), as the token's UUID
	// is needed by the `token` package to generate the worker token.
	wt.Base.ID = uuid.New().String()

	// we set the default expiry of module run tokens to be 6 hours
	expires := time.Now().Add(6 * time.Hour)

	wt.Expires = &expires

	secretData, err := encryption.GenerateRandomBytes(32)

	if err != nil {
		return nil, err
	}

	wt.SigningSecret = []byte(secretData)

	return wt, err
}

func (w *WorkerToken) IsExpired() bool {
	timeLeft := w.Expires.Sub(time.Now())
	return timeLeft < 0
}

func (w *WorkerToken) Encrypt(key *[32]byte) error {
	if !w.HasEncryptedFields.FieldsAreEncrypted {
		ciphertext, err := encryption.Encrypt(w.SigningSecret, key)

		if err != nil {
			return err
		}

		w.SigningSecret = ciphertext

		w.HasEncryptedFields.FieldsAreEncrypted = true
	}

	return nil
}

func (w *WorkerToken) Decrypt(key *[32]byte) error {
	if w.HasEncryptedFields.FieldsAreEncrypted {
		plaintext, err := encryption.Decrypt(w.SigningSecret, key)

		if err != nil {
			return err
		}

		w.SigningSecret = plaintext

		w.HasEncryptedFields.FieldsAreEncrypted = false
	}

	return nil
}
