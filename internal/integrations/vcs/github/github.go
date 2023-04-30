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
	repo      repository.Repository
	appConf   *GithubAppConf
	serverURL string
}

func NewGithubVCSProvider(appConf *GithubAppConf, repo repository.Repository, serverURL string) GithubVCSProvider {
	return GithubVCSProvider{
		appConf:   appConf,
		repo:      repo,
		serverURL: serverURL,
	}
}

func ToGithubVCSProvider(provider vcs.VCSProvider) (res GithubVCSProvider, err error) {
	res, ok := provider.(GithubVCSProvider)

	if !ok {
		return res, fmt.Errorf("could not convert VCS provider to Github VCS provider: %w", err)
	}

	return res, nil
}

func (g GithubVCSProvider) GetGithubAppConfig() *GithubAppConf {
	return g.appConf
}

func (g GithubVCSProvider) GetVCSRepositoryFromModule(depl *models.ModuleDeploymentConfig) (vcs.VCSRepository, error) {
	if depl.GithubAppInstallationID == "" {
		return nil, fmt.Errorf("module does not have github app installation id param set")
	}

	gai, err := g.repo.GithubAppInstallation().ReadGithubAppInstallationByID(depl.GithubAppInstallationID)

	if err != nil {
		return nil, err
	}

	client, err := g.appConf.GetGithubClient(gai.InstallationID)

	if err != nil {
		return nil, err
	}

	return &GithubVCSRepository{
		repoOwner: depl.GitRepoOwner,
		repoName:  depl.GitRepoName,
		serverURL: g.serverURL,
		client:    client,
		repo:      g.repo,
	}, nil
}

// func (g GithubVCSProvider) GetVCSRepositoryFromGAI(gai *models.GithubAppInstallation) (vcs.VCSRepository, error) {
// }

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

// GetArchiveLink returns an archive link for a specific repo SHA
func (g *GithubVCSRepository) GetArchiveLink(ref string) (*url.URL, error) {
	gURL, _, err := g.client.Repositories.GetArchiveLink(
		context.TODO(),
		g.GetRepoOwner(),
		g.GetRepoName(),
		githubsdk.Zipball,
		&githubsdk.RepositoryContentGetOptions{
			Ref: ref,
		},
		false,
	)

	return gURL, err
}

// GetBranch gets a full branch (name and sha)
func (g *GithubVCSRepository) GetBranch(name string) (vcs.VCSBranch, error) {
	branchResp, _, err := g.client.Repositories.GetBranch(
		context.TODO(),
		g.GetRepoOwner(),
		g.GetRepoName(),
		name,
		false,
	)

	if err != nil {
		return nil, err
	}

	return &GithubBranch{branchResp}, nil
}

func (g *GithubVCSRepository) GetPR(mod *models.Module, run *models.ModuleRun) (vcs.VCSRepositoryPullRequest, error) {
	if run.ModuleRunConfig.GithubPullRequestID == 0 {
		return nil, fmt.Errorf("module run does not have github pull request id param set")
	}

	repoPR, err := g.repo.GithubPullRequest().ReadGithubPullRequestByGithubID(mod.TeamID, run.ModuleRunConfig.GithubPullRequestID)

	if err != nil {
		return nil, err
	}

	ghPR, _, err := g.client.PullRequests.Get(
		context.Background(),
		g.GetRepoOwner(),
		g.GetRepoName(),
		int(repoPR.GithubPullRequestNumber),
	)

	if err != nil {
		return nil, err
	}

	return ToVCSRepositoryPullRequest(g.GetRepoOwner(), g.GetRepoName(), ghPR), nil
}

// CreateOrUpdatePRInDatabase stores pull request information using this VCS provider
func (g *GithubVCSRepository) CreateOrUpdatePRInDatabase(teamID string, pr vcs.VCSRepositoryPullRequest) error {
	ghID := pr.GetVCSID().GetIDInt64()
	repoPR, err := g.repo.GithubPullRequest().ReadGithubPullRequestByGithubID(teamID, *ghID)

	if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
		return err
	}

	if repoPR == nil {
		// create a PR object in the db
		repoPR, err = g.repo.GithubPullRequest().CreateGithubPullRequest(&models.GithubPullRequest{
			TeamID:                      teamID,
			GithubRepositoryOwner:       g.GetRepoOwner(),
			GithubRepositoryName:        g.GetRepoName(),
			GithubPullRequestID:         *ghID,
			GithubPullRequestTitle:      pr.GetTitle(),
			GithubPullRequestNumber:     pr.GetPRNumber(),
			GithubPullRequestHeadBranch: pr.GetHeadBranch(),
			GithubPullRequestBaseBranch: pr.GetBaseBranch(),
			GithubPullRequestState:      pr.GetState(),
		})

		if err != nil {
			return err
		}
	} else {
		repoPR.GithubPullRequestTitle = pr.GetTitle()
		repoPR.GithubPullRequestHeadBranch = pr.GetHeadBranch()
		repoPR.GithubPullRequestBaseBranch = pr.GetBaseBranch()
		repoPR.GithubPullRequestState = pr.GetState()

		repoPR, err = g.repo.GithubPullRequest().UpdateGithubPullRequest(repoPR)

		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GithubVCSRepository) CreateCheckRun(
	pr vcs.VCSRepositoryPullRequest,
	mod *models.Module,
	checkRun vcs.VCSCheckRun,
) (vcs.VCSObjectID, error) {
	checkResp, _, err := g.client.Checks.CreateCheckRun(
		context.Background(),
		g.GetRepoOwner(),
		g.GetRepoName(),
		githubsdk.CreateCheckRunOptions{
			Name:    checkRun.Name,
			HeadSHA: pr.GetHeadSHA(),
			Status:  githubsdk.String(string(checkRun.Status)),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error creating new github check run for owner: %s repo %s prNumber: %d. Error: %w",
			g.GetRepoOwner(), g.GetRepoName(), pr.GetPRNumber(), err)
	}

	return vcs.NewVCSObjectInt(checkResp.GetID()), nil
}

func (g *GithubVCSRepository) UpdateCheckRun(pr vcs.VCSRepositoryPullRequest, mod *models.Module, run *models.ModuleRun, checkRun vcs.VCSCheckRun) (vcs.VCSObjectID, error) {
	if run == nil {
		return nil, fmt.Errorf("run cannot be nil")
	}

	checkResp, _, err := g.client.Checks.UpdateCheckRun(
		context.Background(),
		g.GetRepoOwner(),
		g.GetRepoName(),
		run.ModuleRunConfig.GithubCheckID,
		githubsdk.UpdateCheckRunOptions{
			Name:       checkRun.Name,
			Status:     githubsdk.String(string(checkRun.Status)),
			Conclusion: githubsdk.String(string(checkRun.Conclusion)),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error updating new github check run for owner: %s repo %s prNumber: %d. Error: %w",
			g.GetRepoOwner(), g.GetRepoName(), pr.GetPRNumber(), err)
	}

	return vcs.NewVCSObjectInt(checkResp.GetID()), nil
}

// CreateOrUpdateComment creates a comment on a pull/merge request
func (g *GithubVCSRepository) CreateOrUpdateComment(pr vcs.VCSRepositoryPullRequest, mod *models.Module, run *models.ModuleRun, body string) (vcs.VCSObjectID, error) {
	if run != nil && run.ModuleRunConfig.GithubCommentID != 0 {
		commentResp, _, err := g.client.Issues.EditComment(
			context.Background(),
			g.GetRepoOwner(),
			g.GetRepoName(),
			run.ModuleRunConfig.GithubCommentID,
			&githubsdk.IssueComment{
				Body: githubsdk.String(body),
			},
		)

		if err != nil {
			return nil, fmt.Errorf("error updating github comment for owner: %s repo %s prNumber: %d. Error: %w",
				g.GetRepoOwner(), g.GetRepoName(), pr.GetPRNumber(), err)
		}

		return vcs.NewVCSObjectInt(commentResp.GetID()), nil
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
	file, _, err := g.client.Repositories.DownloadContents(
		context.Background(),
		g.GetRepoOwner(),
		g.GetRepoName(),
		path,
		&githubsdk.RepositoryContentGetOptions{
			Ref: ref,
		},
	)

	return file, err
}

// CompareCommits compares a base commit with a head commit
func (g *GithubVCSRepository) CompareCommits(base, head string) (vcs.VCSCommitsComparison, error) {
	commitsRes, _, err := g.client.Repositories.CompareCommits(
		context.Background(),
		g.GetRepoOwner(),
		g.GetRepoName(),
		base,
		head,
		&githubsdk.ListOptions{},
	)

	if err != nil {
		return nil, err
	}

	return &GithubCommitsComparison{commitsRes}, nil
}

// PopulateModuleRun adds additional fields to the module run config. Should be called
// before creating the module run.
func (g *GithubVCSRepository) PopulateModuleRun(run *models.ModuleRun, prID, checkRunID, commentID vcs.VCSObjectID) {
	run.ModuleRunConfig.GithubCheckID = *checkRunID.GetIDInt64()
	run.ModuleRunConfig.GithubCommentID = *commentID.GetIDInt64()
	run.ModuleRunConfig.GithubPullRequestID = *prID.GetIDInt64()
}
