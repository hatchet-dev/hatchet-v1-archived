package models

import (
	"encoding/json"
	"sort"

	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/encryption"
)

type ModuleEnvVarsVersion struct {
	Base
	HasEncryptedFields

	ModuleID string
	Version  uint

	// JSON-based representation of module values, encrypted before storage
	EnvVars []byte
}

func NewModuleEnvVarsVersion(moduleID string, prevVersion uint, vars map[string]string) (*ModuleEnvVarsVersion, error) {
	envVarsBytes, err := json.Marshal(&vars)

	if err != nil {
		return nil, err
	}

	return &ModuleEnvVarsVersion{
		ModuleID: moduleID,
		Version:  prevVersion + 1,
		EnvVars:  envVarsBytes,
	}, nil
}

func (m *ModuleEnvVarsVersion) GetEnvVars(key *[32]byte) (map[string]string, error) {
	err := m.Decrypt(key)

	if err != nil {
		return nil, err
	}

	res := make(map[string]string)

	err = json.Unmarshal(m.EnvVars, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *ModuleEnvVarsVersion) Encrypt(key *[32]byte) error {
	if !m.HasEncryptedFields.FieldsAreEncrypted {
		ciphertext, err := encryption.Encrypt(m.EnvVars, key)

		if err != nil {
			return err
		}

		m.EnvVars = ciphertext

		m.HasEncryptedFields.FieldsAreEncrypted = true
	}

	return nil
}

func (m *ModuleEnvVarsVersion) Decrypt(key *[32]byte) error {
	if m.HasEncryptedFields.FieldsAreEncrypted {
		plaintext, err := encryption.Decrypt(m.EnvVars, key)

		if err != nil {
			return err
		}

		m.EnvVars = plaintext

		m.HasEncryptedFields.FieldsAreEncrypted = false
	}

	return nil
}

func (m *ModuleEnvVarsVersion) ToAPIType(key *[32]byte) (*types.ModuleEnvVarsVersion, error) {
	vars, err := m.GetEnvVars(key)

	if err != nil {
		return nil, err
	}

	envVars := make([]types.ModuleEnvVar, 0)

	for key, val := range vars {
		envVars = append(envVars, types.ModuleEnvVar{
			Key: key,
			Val: val,
		})
	}

	// sort by stable alphanumeric ordering
	sort.SliceStable(envVars, func(i, j int) bool {
		return envVars[i].Key > envVars[j].Key
	})

	return &types.ModuleEnvVarsVersion{
		APIResourceMeta: m.ToAPITypeMetadata(),
		Version:         m.Version,
		EnvVars:         envVars,
	}, nil
}
