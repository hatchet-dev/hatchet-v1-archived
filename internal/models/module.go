package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/encryption"
	"gorm.io/gorm"
)

type DeploymentMechanism string

const (
	DeploymentMechanismGithub DeploymentMechanism = "github"
	DeploymentMechanismAPI    DeploymentMechanism = "api"
)

type ModuleLockKind string

const (
	ModuleLockKindGithubBranch ModuleLockKind = "github_branch"
	ModuleLockKindManual       ModuleLockKind = "manual"
)

type Module struct {
	Base

	TeamID string
	Team   Team `gorm:"foreignKey:TeamID"`

	Name string

	DeploymentMechanism DeploymentMechanism

	DeploymentConfig ModuleDeploymentConfig

	CurrentModuleValuesVersionID string
	CurrentModuleValuesVersion   ModuleValuesVersion `gorm:"foreignKey:CurrentModuleValuesVersionID"`

	CurrentModuleEnvVarsVersionID string
	CurrentModuleEnvVarsVersion   ModuleEnvVarsVersion `gorm:"foreignKey:CurrentModuleEnvVarsVersionID"`

	ModuleRunQueueID string
	ModuleRunQueue   ModuleRunQueue `gorm:"foreignKey:ModuleRunQueueID"`

	// LockID represents a unique lock ID for the module. This operates at a higher level than the
	// Terraform state lock. For a LockKind of type "github," this corresponds to a commit SHA.
	LockID string

	// LockKind describes the type of lock.
	LockKind ModuleLockKind

	Runs []ModuleRun
}

func (m *Module) ToAPIType() *types.Module {
	return &types.Module{
		APIResourceMeta:         m.ToAPITypeMetadata(),
		Name:                    m.Name,
		DeploymentConfig:        *m.DeploymentConfig.ToAPIType(),
		LockID:                  m.LockID,
		LockKind:                types.ModuleLockKind(m.LockKind),
		CurrentValuesVersionID:  m.CurrentModuleValuesVersionID,
		CurrentEnvVarsVersionID: m.CurrentModuleEnvVarsVersionID,
	}
}

func (m *Module) AfterFind(tx *gorm.DB) (err error) {
	// this ensures that AfterFind is called on the invite link even if called with a
	// Joins method, instead of just Preload
	return m.CurrentModuleEnvVarsVersion.AfterFind(tx)
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

const LogLocationFileStorage string = "file"

type ModuleRun struct {
	Base

	// TeamID is only used by some queries where the team id is not implicit. This is not
	// written to the module run table.
	TeamID string `gorm:"-"`

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

	LogLocation string
}

func (m *ModuleRun) AfterFind(tx *gorm.DB) (err error) {
	// this ensures that AfterFind is called on the invite link even if called with a
	// Joins method, instead of just Preload
	return m.ModuleRunConfig.AfterFind(tx)
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

func (m *ModuleRun) ToAPITypeOverview() *types.ModuleRunOverview {
	return &types.ModuleRunOverview{
		APIResourceMeta:   m.ToAPITypeMetadata(),
		Status:            types.ModuleRunStatus(m.Status),
		StatusDescription: m.StatusDescription,
		Kind:              types.ModuleRunKind(m.Kind),
	}
}

func (m *ModuleRun) ToAPIType(pr *GithubPullRequest) *types.ModuleRun {
	res := &types.ModuleRun{
		ModuleRunOverview: m.ToAPITypeOverview(),
	}

	if pr != nil {
		res.ModuleRunPullRequest = pr.ToAPIType()
	}

	if mc := m.ModuleRunConfig; mc.ID != "" {
		res.ModuleRunConfig = &types.ModuleRunConfig{
			TriggerKind:     types.ModuleRunTriggerKind(mc.TriggerKind),
			GithubCommitSHA: mc.GithubCommitSHA,
			EnvVarVersionID: mc.ModuleEnvVarsVersionID,
			ValuesVersionID: mc.ModuleValuesVersionID,
		}
	}

	return res
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

	ModuleValuesVersionID string
	ModuleValuesVersion   ModuleValuesVersion `gorm:"foreignKey:ModuleValuesVersionID"`

	ModuleEnvVarsVersionID string
	ModuleEnvVarsVersion   ModuleEnvVarsVersion `gorm:"foreignKey:ModuleEnvVarsVersionID"`
}

func (m *ModuleRunConfig) AfterFind(tx *gorm.DB) (err error) {
	// this ensures that AfterFind is called on the invite link even if called with a
	// Joins method, instead of just Preload
	return m.ModuleEnvVarsVersion.AfterFind(tx)
}
