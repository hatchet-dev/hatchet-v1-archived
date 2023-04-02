package github

import (
	githubsdk "github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
)

type ghPullRequest struct {
	e *githubsdk.PullRequestEvent
}

func ToVCSRepositoryPullRequest(event *githubsdk.PullRequestEvent) vcs.VCSRepositoryPullRequest {
	return &ghPullRequest{event}
}

func (g *ghPullRequest) GetRepoOwner() string {
	return g.e.GetRepo().GetOwner().GetLogin()
}

func (g *ghPullRequest) GetRepoName() string {
	return g.e.GetRepo().GetName()
}

func (g *ghPullRequest) GetPRNumber() int64 {
	return g.e.GetPullRequest().GetID()
}

func (g *ghPullRequest) GetBaseSHA() string {
	return g.e.GetPullRequest().GetBase().GetSHA()
}

func (g *ghPullRequest) GetHeadSHA() string {
	return g.e.GetPullRequest().GetHead().GetSHA()
}

func (g *ghPullRequest) GetBaseBranch() string {
	return g.e.GetPullRequest().GetBase().GetRef()
}

func (g *ghPullRequest) GetHeadBranch() string {
	return g.e.GetPullRequest().GetHead().GetRef()
}

func (g *ghPullRequest) GetTitle() string {
	return g.e.GetPullRequest().GetTitle()
}

func (g *ghPullRequest) GetState() string {
	return g.e.GetPullRequest().GetState()
}
