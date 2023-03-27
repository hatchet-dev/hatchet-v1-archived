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

	// CreateOrUpdatePR stores pull request information using this VCS provider
	CreateOrUpdatePR(pr VCSRepositoryPullRequest)

	// CreateOrUpdateCheckRun creates a new "check" to run against a PR
	CreateOrUpdateCheckRun(pr VCSRepositoryPullRequest, mod *models.Module) error

	// CreateOrUpdateComment creates or updates a comment on a pull or merge request
	CreateOrUpdateComment(pr VCSRepositoryPullRequest, mod *models.Module, body string) error

	// ReadFile returns a file by a SHA reference or path
	ReadFile(ref, path string) (io.ReadCloser, error)

	// CompareCommits compares a base commit with a head commit
	CompareCommits(base, head string) (CommitsComparison, error)
}

type VCSProvider interface {
	// GetVCSRepositoryFromModule returns the corresponding VCS repository for the module.
	// Callers should likely use the package method GetVCSProviderFromModule.
	GetVCSRepositoryFromModule(mod *models.Module) (VCSRepository, error)
}

// GetVCSProviderFromModule returns the corresponding VCS provider for the module
func GetVCSProviderFromModule(allProviders map[VCSRepositoryKind]VCSProvider, mod *models.Module) (VCSRepository, error) {
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

	return provider.GetVCSRepositoryFromModule(mod)
}

// VCSRepositoryPullRequest abstracts the underlying pull or merge request methods to only
// extract relevant information.
type VCSRepositoryPullRequest interface {
	GetRepoOwner() string
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

type CommitsComparison struct {
	Files []CommitFile
}

type CommitFile struct {
	Name string
}

func (f CommitFile) GetFilename() string {
	return f.Name
}
