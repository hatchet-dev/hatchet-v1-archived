package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/encryption"
)

type DeploymentMechanism string

const (
	DeploymentMechanismGithub DeploymentMechanism = "github"
	DeploymentMechanismAPI    DeploymentMechanism = "api"
)

type Module struct {
	Base

	TeamID string
	Team   Team `gorm:"foreignKey:TeamID"`

	Name string

	DeploymentMechanism DeploymentMechanism

	DeploymentConfig ModuleDeploymentConfig

	Runs []ModuleRun
}

func (m *Module) ToAPIType() *types.Module {
	return &types.Module{
		APIResourceMeta:  m.ToAPITypeMetadata(),
		Name:             m.Name,
		DeploymentConfig: *m.DeploymentConfig.ToAPIType(),
	}
}

type ModuleDeploymentConfig struct {
	Base

	ModuleID string

	ModulePath string

	GithubRepoName   string
	GithubRepoOwner  string
	GithubRepoBranch string

	GithubAppInstallationID string
	GithubAppInstallation   GithubAppInstallation `gorm:"foreignKey:GithubAppInstallationID"`
}

func (m *ModuleDeploymentConfig) ToAPIType() *types.ModuleDeploymentConfig {
	return &types.ModuleDeploymentConfig{
		Path:                    m.ModulePath,
		GithubRepoName:          m.GithubRepoName,
		GithubRepoOwner:         m.GithubRepoOwner,
		GithubRepoBranch:        m.GithubRepoBranch,
		GithubAppInstallationID: m.GithubAppInstallationID,
	}
}

type ModuleRunStatus string

const (
	ModuleRunStatusQueued     ModuleRunStatus = "queued"
	ModuleRunStatusInProgress ModuleRunStatus = "in_progress"
	ModuleRunStatusCompleted  ModuleRunStatus = "completed"
	ModuleRunStatusFailed     ModuleRunStatus = "failed"
)

type ModuleRunKind string

const (
	ModuleRunKindPlan    ModuleRunKind = "plan"
	ModuleRunKindApply   ModuleRunKind = "apply"
	ModuleRunKindDestroy ModuleRunKind = "destroy"
)

type ModuleRun struct {
	Base

	ModuleID string

	Status            ModuleRunStatus
	StatusDescription string

	Kind ModuleRunKind

	LockID        string
	LockOperation string
	LockInfo      string
	LockWho       string
	LockVersion   string
	LockCreated   string
	LockPath      string

	Tokens []ModuleRunToken

	ModuleRunConfig ModuleRunConfig
}

func (m *ModuleRun) ToTerraformLockType() *types.TerraformLock {
	return &types.TerraformLock{
		ID:        m.LockID,
		Operation: m.LockOperation,
		Info:      m.LockInfo,
		Who:       m.LockWho,
		Version:   m.LockVersion,
		Created:   m.LockCreated,
		Path:      m.LockPath,
	}
}

func (m *ModuleRun) ToAPIType() *types.ModuleRun {
	return &types.ModuleRun{
		APIResourceMeta:   m.ToAPITypeMetadata(),
		Status:            types.ModuleRunStatus(m.Status),
		StatusDescription: m.StatusDescription,
		Kind:              types.ModuleRunKind(m.Kind),
	}
}

type ModuleRunToken struct {
	Base
	HasEncryptedFields

	// The subject of the token (service account user)
	UserID string

	// The run id that this token was created for
	ModuleRunID string

	// When this PAT expires. This should match what's in the JWT data
	Expires *time.Time

	// Whether the personal access token has been revoked
	Revoked bool

	// Encrypted data that contains the token signing secret for that specific token
	SigningSecret []byte
}

func NewModuleRunTokenFromRunID(userID, runID string) (*ModuleRunToken, error) {
	mrt := &ModuleRunToken{
		UserID:      userID,
		ModuleRunID: runID,
	}

	// in this case, we generate the UUID ahead of time (rather than BeforeCreate), as the token's UUID
	// is needed by the `token` package to generate the JWT.
	mrt.Base.ID = uuid.New().String()

	// we set the default expiry of module run tokens to be 6 hours
	expires := time.Now().Add(6 * time.Hour)

	mrt.Expires = &expires

	secretData, err := encryption.GenerateRandomBytes(32)

	if err != nil {
		return nil, err
	}

	mrt.SigningSecret = []byte(secretData)

	return mrt, err
}

func (m *ModuleRunToken) IsExpired() bool {
	timeLeft := m.Expires.Sub(time.Now())
	return timeLeft < 0
}

func (m *ModuleRunToken) Encrypt(key *[32]byte) error {
	if !m.HasEncryptedFields.FieldsAreEncrypted {
		ciphertext, err := encryption.Encrypt(m.SigningSecret, key)

		if err != nil {
			return err
		}

		m.SigningSecret = ciphertext

		m.HasEncryptedFields.FieldsAreEncrypted = true
	}

	return nil
}

func (m *ModuleRunToken) Decrypt(key *[32]byte) error {
	if m.HasEncryptedFields.FieldsAreEncrypted {
		plaintext, err := encryption.Decrypt(m.SigningSecret, key)

		if err != nil {
			return err
		}

		m.SigningSecret = plaintext

		m.HasEncryptedFields.FieldsAreEncrypted = false
	}

	return nil
}

type ModuleRunTriggerKind string

const (
	ModuleRunTriggerKindGithub ModuleRunTriggerKind = "github"
	ModuleRunTriggerKindManual ModuleRunTriggerKind = "manual"
)

type ModuleRunConfig struct {
	Base

	ModuleRunID string

	TriggerKind ModuleRunTriggerKind

	GithubCheckID   int64
	GithubCommentID int64
	GithubCommitSHA string
}
