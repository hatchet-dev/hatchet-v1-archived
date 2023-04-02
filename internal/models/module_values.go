package models

import (
	"encoding/json"

	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/encryption"
)

type ModuleValuesVersionKind string

const (
	ModuleValuesVersionKindDatabase ModuleValuesVersionKind = "db"
	ModuleValuesVersionKindVCS      ModuleValuesVersionKind = "vcs"
)

type ModuleValuesVersion struct {
	Base

	ModuleID string
	Version  uint

	Kind ModuleValuesVersionKind

	// Git-specific params
	GitValuesPath string
	GitRepoName   string
	GitRepoOwner  string
	GitRepoBranch string

	// Github-specific params
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

		gh.Path = m.GitValuesPath
		gh.GithubRepoOwner = m.GitRepoOwner
		gh.GithubRepoName = m.GitRepoName
		gh.GithubAppInstallationID = m.GithubAppInstallationID
		gh.GithubRepoBranch = m.GitRepoBranch

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
