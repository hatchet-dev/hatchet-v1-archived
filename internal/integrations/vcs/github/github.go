package github

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	githubsdk "github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

type GithubVCSProvider struct {
	repo    repository.Repository
	appConf *GithubAppConf
}

func NewGithubVCSProvider(appConf *GithubAppConf, repo repository.Repository) GithubVCSProvider {
	return GithubVCSProvider{
		appConf: appConf,
		repo:    repo,
	}
}

func ToGithubVCSProviderFactory(provider vcs.VCSProvider) (res GithubVCSProvider, err error) {
	res, ok := provider.(GithubVCSProvider)

	if !ok {
		return res, fmt.Errorf("could not convert VCS provider factory to Github VCS provider factory: %w", err)
	}

	return res, nil
}

func (g GithubVCSProvider) GetGithubAppConfig() *GithubAppConf {
	return g.appConf
}

func (g GithubVCSProvider) GetVCSRepositoryFromModule(mod *models.Module) (vcs.VCSRepository, error) {
	if mod.DeploymentConfig.GithubAppInstallationID == "" {
		return nil, fmt.Errorf("module does not have github app installation id param set")
	}

	gai, err := g.repo.GithubAppInstallation().ReadGithubAppInstallationByID(mod.DeploymentConfig.GithubAppInstallationID)

	if err != nil {
		return nil, err
	}

	return g.GetVCSRepositoryFromGAI(gai)
}

func (g GithubVCSProvider) GetVCSRepositoryFromGAI(gai *models.GithubAppInstallation) (vcs.VCSRepository, error) {
	client, err := g.appConf.GetGithubClient(gai.InstallationID)

	if err != nil {
		return nil, err
	}

	return &GithubVCSRepository{
		client: client,
	}, nil
}

type GithubVCSRepository struct {
	repoOwner, repoName string
	client              *githubsdk.Client
	repo                repository.Repository
	serverURL           string
}

// GetKind returns the kind of VCS provider -- used for downstream integrations
func (g *GithubVCSRepository) GetKind() vcs.VCSRepositoryKind {
	return vcs.VCSRepositoryKindGithub
}

func (g *GithubVCSRepository) GetRepoOwner() string {
	return g.repoOwner
}

func (g *GithubVCSRepository) GetRepoName() string {
	return g.repoName
}

// SetupRepository sets up a VCS repository on Hatchet.
func (g *GithubVCSRepository) SetupRepository(teamID string) error {
	repoOwner := g.GetRepoOwner()
	repoName := g.GetRepoName()

	_, err := g.repo.GithubWebhook().ReadGithubWebhookByTeamID(teamID, repoOwner, repoName)

	if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
		return err
	} else if err != nil {
		gw, err := models.NewGithubWebhook(teamID, repoOwner, repoName)

		if err != nil {
			return err
		}

		gw, err = g.repo.GithubWebhook().CreateGithubWebhook(gw)

		if err != nil {
			return err
		}

		webhookURL := fmt.Sprintf("%s/api/v1/teams/%s/github_incoming/%s", g.serverURL, teamID, gw.ID)

		_, _, err = g.client.Repositories.CreateHook(
			context.Background(), repoOwner, repoName, &githubsdk.Hook{
				Config: map[string]interface{}{
					"url":          webhookURL,
					"content_type": "json",
					"secret":       string(gw.SigningSecret),
				},
				Events: []string{"pull_request", "push"},
				Active: githubsdk.Bool(true),
			},
		)

		return err
	}

	return nil
}

// OnPullRequestEvent is triggered when a new pull request is opened, closed, edited, etc.
func (g *GithubVCSRepository) OnPullRequestEvent(pr vcs.VCSRepositoryPullRequest) error {
	panic("unimplemented")
}

// GetArchiveLink returns an archive link for a specific repo SHA
func (g *GithubVCSRepository) GetArchiveLink(ref string) (*url.URL, error) {
	panic("unimplemented")
}

// GetBranch gets a full branch (name and sha)
func (g *GithubVCSRepository) GetBranch(name string) (vcs.VCSBranch, error) {
	panic("unimplemented")
}

// CreateOrUpdatePR stores pull request information using this VCS provider
func (g *GithubVCSRepository) CreateOrUpdatePR(pr vcs.VCSRepositoryPullRequest) (vcs.VCSObjectID, error) {
	panic("unimplemented")
}

// CreateOrUpdateCheckRun creates a new "check" to run against a PR
func (g *GithubVCSRepository) CreateOrUpdateCheckRun(pr vcs.VCSRepositoryPullRequest, mod *models.Module, run *models.ModuleRun) (vcs.VCSObjectID, error) {
	panic("unimplemented")
}

// CreateOrUpdateComment creates or updates a comment on a pull or merge request
func (g *GithubVCSRepository) CreateOrUpdateComment(pr vcs.VCSRepositoryPullRequest, mod *models.Module, run *models.ModuleRun, body string) (vcs.VCSObjectID, error) {
	if run != nil && run.ModuleRunConfig.GithubCommentID != 0 {
		g.repo.GithubPullRequest().ReadGithubPullRequestCommentByGithubID(mod.ID, run.ModuleRunConfig.GithubCommentID)

		// g.client.
	}

	// if no existing id, create the id
	commentResp, _, err := g.client.Issues.CreateComment(
		context.Background(),
		g.GetRepoOwner(),
		g.GetRepoName(),
		int(pr.GetPRNumber()),
		&githubsdk.IssueComment{
			Body: &body,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error creating new github comment for owner: %s repo %s prNumber: %d. Error: %w",
			g.GetRepoOwner(), g.GetRepoName(), pr.GetPRNumber(), err)
	}

	return vcs.NewVCSObjectInt(commentResp.GetID()), nil
}

// ReadFile returns a file by a SHA reference or path
func (g *GithubVCSRepository) ReadFile(ref, path string) (io.ReadCloser, error) {
	panic("unimplemented")
}

// CompareCommits compares a base commit with a head commit
func (g *GithubVCSRepository) CompareCommits(base, head string) (vcs.CommitsComparison, error) {
	panic("unimplemented")
}

// PopulateModuleRun adds additional fields to the module run config. Should be called
// before creating the module run.
func (g *GithubVCSRepository) PopulateModuleRun(run *models.ModuleRun, checkRunID, commentID vcs.VCSObjectID) {
	run.ModuleRunConfig.GithubCheckID = *checkRunID.GetIDInt64()
	run.ModuleRunConfig.GithubCommentID = *checkRunID.GetIDInt64()
}
