package vcs

import (
	"fmt"
	"io"
	"net/url"

	"github.com/hatchet-dev/hatchet/internal/models"
)

type VCSProviderKind string

const (
	// enumeration of in-tree providers
	VCSProviderKindGithub VCSProviderKind = "github"
	VCSProviderKindGitlab VCSProviderKind = "gitlab"
)

// VCSProvider provides an interface to implement for new version control providers
// (github, gitlab, etc)
type VCSProvider interface {
	// GetKind returns the kind of VCS provider -- used for downstream integrations
	GetKind() VCSProviderKind

	// SetupRepository sets up a VCS repository on Hatchet.
	SetupRepository(teamID string, r VCSRepo) error

	// GetArchiveLink returns an archive link for a specific repo SHA
	GetArchiveLink(ref string) (*url.URL, error)

	// GetBranch gets a full branch (name and sha)
	GetBranch(name string) (VCSBranch, error)

	// CreateOrUpdatePR stores pull request information using this VCS provider
	CreateOrUpdatePR(pr VCSProviderPullRequest)

	// CreateOrUpdateCheckRun creates a new "check" to run against a PR
	CreateOrUpdateCheckRun(pr VCSProviderPullRequest, mod *models.Module) error

	// CreateOrUpdateComment creates or updates a comment on a pull or merge request
	CreateOrUpdateComment(pr VCSProviderPullRequest, mod *models.Module, body string) error

	// ReadFile returns a file by a SHA reference or path
	ReadFile(ref, path string) (io.ReadCloser, error)

	// CompareCommits compares a base commit with a head commit
	CompareCommits(base, head string) (CommitsComparison, error)
}

type VCSProviderFactory interface {
	// GetVCSProviderFromModule returns the corresponding VCS provider for the module.
	// Callers should likely use the package method GetVCSProviderFromModule.
	GetVCSProviderFromModule(mod *models.Module) (VCSProvider, error)
}

// GetVCSProviderFromModule returns the corresponding VCS provider for the module
func GetVCSProviderFromModule(allProviders map[VCSProviderKind]VCSProviderFactory, mod *models.Module) (VCSProvider, error) {
	var providerFact VCSProviderFactory
	var providerKind VCSProviderKind

	if mod.DeploymentMechanism == models.DeploymentMechanismGithub {
		providerKind = VCSProviderKindGithub

	} else if mod.DeploymentMechanism == models.DeploymentMechanismGitlab {
		providerKind = VCSProviderKindGitlab
	} else {
		return nil, fmt.Errorf("module %s does not use a VCS integration", mod.ID)
	}

	providerFact, exists := allProviders[providerKind]

	if !exists {
		return nil, fmt.Errorf("VCS provider kind '%s' is not registered on this Hatchet instance", providerKind)
	}

	return providerFact.GetVCSProviderFromModule(mod)
}

// VCSProviderPullRequest abstracts the underlying pull or merge request methods to only
// extract relevant information.
type VCSProviderPullRequest interface {
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

type VCSRepoKind string

const (
	VCSRepoKindGithub = "github"
	VCSRepoKindGitlab = "gitlab"
)

type VCSRepo interface {
	GetKind() VCSRepoKind
	GetRepoName() string
	GetRepoOwner() string
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
