package models

import (
	"github.com/google/uuid"
	"github.com/hatchet-dev/hatchet/internal/encryption"
)

// GithubWebhook contains data for a Github webhook
type GithubWebhook struct {
	Base
	HasEncryptedFields

	TeamID string

	GithubRepositoryOwner string
	GithubRepositoryName  string

	// Encrypted data that contains the webhook signing secret
	SigningSecret []byte

	GithubAppInstallations []GithubAppInstallation `gorm:"many2many:github_app_installations;"`
}

func NewGithubWebhook(teamID, repoOwner, repoName string) (*GithubWebhook, error) {
	gw := &GithubWebhook{
		TeamID:                teamID,
		GithubRepositoryOwner: repoOwner,
		GithubRepositoryName:  repoName,
	}

	// in this case, we generate the UUID ahead of time (rather than BeforeCreate), as the token's UUID
	// is needed by the `token` package to generate the JWT.
	gw.Base.ID = uuid.New().String()

	secretData, err := encryption.GenerateRandomBytes(32)

	if err != nil {
		return nil, err
	}

	gw.SigningSecret = []byte(secretData)

	return gw, nil
}

func (gw *GithubWebhook) Encrypt(key *[32]byte) error {
	if !gw.HasEncryptedFields.FieldsAreEncrypted {
		ciphertext, err := encryption.Encrypt(gw.SigningSecret, key)

		if err != nil {
			return err
		}

		gw.SigningSecret = ciphertext

		gw.HasEncryptedFields.FieldsAreEncrypted = true
	}

	return nil
}

func (gw *GithubWebhook) Decrypt(key *[32]byte) error {
	if gw.HasEncryptedFields.FieldsAreEncrypted {
		plaintext, err := encryption.Decrypt(gw.SigningSecret, key)

		if err != nil {
			return err
		}

		gw.SigningSecret = plaintext

		gw.HasEncryptedFields.FieldsAreEncrypted = false
	}

	return nil
}
