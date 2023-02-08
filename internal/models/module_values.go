package models

import "github.com/hatchet-dev/hatchet/internal/encryption"

type ModuleValuesVersionKind string

const (
	ModuleValuesVersionKindDatabase ModuleValuesVersionKind = "db"
	ModuleValuesVersionKindGithub   ModuleValuesVersionKind = "github"
)

type ModuleValuesVersion struct {
	Base

	ModuleID string
	Version  uint

	Kind ModuleValuesVersionKind

	// Github configuration params
	GithubValuesPath        string
	GithubRepoName          string
	GithubRepoOwner         string
	GithubRepoBranch        string
	GithubAppInstallationID string
	GithubAppInstallation   GithubAppInstallation `gorm:"foreignKey:GithubAppInstallationID"`
}

type ModuleValues struct {
	Base
	HasEncryptedFields

	ModuleValuesVersionID string

	// JSON-based representation of module values, encrypted before storage
	Values []byte
}

func (m *ModuleValues) Encrypt(key *[32]byte) error {
	if !m.HasEncryptedFields.FieldsAreEncrypted {
		ciphertext, err := encryption.Encrypt(m.Values, key)

		if err != nil {
			return err
		}

		m.Values = ciphertext

		m.HasEncryptedFields.FieldsAreEncrypted = true
	}

	return nil
}

func (m *ModuleValues) Decrypt(key *[32]byte) error {
	if m.HasEncryptedFields.FieldsAreEncrypted {
		plaintext, err := encryption.Decrypt(m.Values, key)

		if err != nil {
			return err
		}

		m.Values = plaintext

		m.HasEncryptedFields.FieldsAreEncrypted = false
	}

	return nil
}
