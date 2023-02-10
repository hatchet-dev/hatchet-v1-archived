package models

import (
	"encoding/json"

	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/encryption"
)

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

func (m *ModuleValuesVersion) ToAPIType(mv *ModuleValues) (*types.ModuleValues, error) {
	res := &types.ModuleValues{
		APIResourceMeta: m.ToAPITypeMetadata(),
		Version:         m.Version,
	}

	if m.Kind == ModuleValuesVersionKindDatabase {
		if mv != nil {
			jsonVal := make(map[string]interface{})

			err := json.Unmarshal(mv.Values, &jsonVal)

			if err != nil {
				return nil, err
			}

			res.Values = jsonVal
		}
	} else {
		gh := &types.ModuleValuesGithubConfig{}

		gh.Path = m.GithubValuesPath
		gh.GithubRepoOwner = m.GithubRepoOwner
		gh.GithubRepoName = m.GithubRepoName
		gh.GithubAppInstallationID = m.GithubAppInstallationID
		gh.GithubRepoBranch = m.GithubRepoBranch

		res.Github = gh
	}

	return res, nil

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
