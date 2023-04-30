package vcs

import (
	"fmt"
	"io"
	"net/url"

	"github.com/hatchet-dev/hatchet/internal/models"
)

type VCSRepositoryKind string

const (
	// enumeration of in-tree repository kinds
	VCSRepositoryKindGithub VCSRepositoryKind = "github"
	VCSRepositoryKindGitlab VCSRepositoryKind = "gitlab"
)

type VCSCheckRunStatus string

const (
	VCSCheckRunStatusQueued     VCSCheckRunStatus = "queued"
	VCSCheckRunStatusInProgress VCSCheckRunStatus = "in_progress"
	VCSCheckRunStatusCompleted  VCSCheckRunStatus = "completed"
)

type VCSCheckRunConclusion string

const (
	VCSCheckRunConclusionSuccess        VCSCheckRunConclusion = "success"
	VCSCheckRunConclusionFailure        VCSCheckRunConclusion = "failure"
	VCSCheckRunConclusionCancelled      VCSCheckRunConclusion = "cancelled"
	VCSCheckRunConclusionSkipped        VCSCheckRunConclusion = "skipped"
	VCSCheckRunConclusionTimedOut       VCSCheckRunConclusion = "timed_out"
	VCSCheckRunConclusionActionRequired VCSCheckRunConclusion = "action_required"
)

type VCSCheckRun struct {
	Name       string
	Status     VCSCheckRunStatus
	Conclusion VCSCheckRunConclusion
}

// VCSRepository provides an interface to implement for new version control providers
// (github, gitlab, etc)
type VCSRepository interface {
	// GetKind returns the kind of VCS provider -- used for downstream integrations
	GetKind() VCSRepositoryKind

	// GetRepoOwner returns the owner of the repository
	GetRepoOwner() string

	// GetRepoName returns the name of the repository
	GetRepoName() string

	// SetupRepository sets up a VCS repository on Hatchet.
	SetupRepository(teamID string) error

	// GetArchiveLink returns an archive link for a specific repo SHA
	GetArchiveLink(ref string) (*url.URL, error)

	// GetBranch gets a full branch (name and sha)
	GetBranch(name string) (VCSBranch, error)

	// GetPR retrieves a pull request from the VCS provider
	GetPR(mod *models.Module, run *models.ModuleRun) (VCSRepositoryPullRequest, error)

	// CreateOrUpdatePRInDatabase stores pull request information using this VCS provider
	CreateOrUpdatePRInDatabase(teamID string, pr VCSRepositoryPullRequest) error

	// CreateCheckRun creates a new "check" to run against a PR
	CreateCheckRun(
		pr VCSRepositoryPullRequest,
		mod *models.Module,
		checkRun VCSCheckRun,
	) (VCSObjectID, error)

	// CreateOrUpdateCheckRun creates a new "check" to run against a PR
	UpdateCheckRun(pr VCSRepositoryPullRequest, mod *models.Module, run *models.ModuleRun, checkRun VCSCheckRun) (VCSObjectID, error)

	// CreateOrUpdateComment creates or updates a comment on a pull or merge request
	// Note that run MAY be nil, implementations should check if the run is nil before reading from it
	CreateOrUpdateComment(pr VCSRepositoryPullRequest, mod *models.Module, run *models.ModuleRun, body string) (VCSObjectID, error)

	// PopulateModuleRun adds additional fields to the module run config. Should be called
	// before creating the module run.
	PopulateModuleRun(run *models.ModuleRun, prID, checkRunID, commentID VCSObjectID)

	// ReadFile returns a file by a SHA reference or path
	ReadFile(ref, path string) (io.ReadCloser, error)

	// CompareCommits compares a base commit with a head commit
	CompareCommits(base, head string) (VCSCommitsComparison, error)
}

// VCSObjectID is a generic method for retrieving IDs from the underlying VCS repository.
// Depending on the provider, object IDs may be int64 or strings.
//
// Object ids are meant to be passed between methods in the same VCS repository, so they should
// only be read by the same VCS repository that wrote them.
type VCSObjectID interface {
	GetIDString() *string
	GetIDInt64() *int64
}

type vcsObjectString struct {
	id string
}

func NewVCSObjectString(id string) VCSObjectID {
	return vcsObjectString{id}
}

func (v vcsObjectString) GetIDString() *string {
	return &v.id
}

func (v vcsObjectString) GetIDInt64() *int64 {
	return nil
}

type vcsObjectInt struct {
	id int64
}

func NewVCSObjectInt(id int64) VCSObjectID {
	return vcsObjectInt{id}
}

func (v vcsObjectInt) GetIDString() *string {
	return nil
}

func (v vcsObjectInt) GetIDInt64() *int64 {
	return &v.id
}

type VCSProvider interface {
	// GetVCSRepositoryFromModule returns the corresponding VCS repository for the module.
	// Callers should likely use the package method GetVCSProviderFromModule.
	GetVCSRepositoryFromModule(depl *models.ModuleDeploymentConfig) (VCSRepository, error)
}

// GetVCSRepositoryFromModule returns the corresponding VCS repository for the module
func GetVCSRepositoryFromModule(allProviders map[VCSRepositoryKind]VCSProvider, mod *models.Module) (VCSRepository, error) {
	var repoKind VCSRepositoryKind

	if mod.DeploymentMechanism == models.DeploymentMechanismGithub {
		repoKind = VCSRepositoryKindGithub

	} else if mod.DeploymentMechanism == models.DeploymentMechanismGitlab {
		repoKind = VCSRepositoryKindGitlab
	} else {
		return nil, fmt.Errorf("module %s does not use a VCS integration", mod.ID)
	}

	provider, exists := allProviders[repoKind]

	if !exists {
		return nil, fmt.Errorf("VCS provider kind '%s' is not registered on this Hatchet instance", repoKind)
	}

	return provider.GetVCSRepositoryFromModule(&mod.DeploymentConfig)
}

// VCSRepositoryPullRequest abstracts the underlying pull or merge request methods to only
// extract relevant information.
type VCSRepositoryPullRequest interface {
	GetRepoOwner() string
	GetVCSID() VCSObjectID
	GetPRNumber() int64
	GetRepoName() string
	GetBaseSHA() string
	GetHeadSHA() string
	GetBaseBranch() string
	GetHeadBranch() string
	GetTitle() string
	GetState() string
}

type VCSBranch interface {
	GetName() string
	GetLatestRef() string
}

type VCSCommitsComparison interface {
	GetFiles() []CommitFile
}

type CommitFile struct {
	Name string
}

func (f CommitFile) GetFilename() string {
	return f.Name
}
