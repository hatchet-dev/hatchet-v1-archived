package github

import (
	githubsdk "github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
)

type ghPullRequest struct {
	repoOwner, repoName string
	pr                  *githubsdk.PullRequest
}

func ToVCSRepositoryPullRequest(repoOwner, repoName string, pr *githubsdk.PullRequest) vcs.VCSRepositoryPullRequest {
	return &ghPullRequest{repoOwner, repoName, pr}
}

func (g *ghPullRequest) GetRepoOwner() string {
	return g.repoOwner
}

func (g *ghPullRequest) GetRepoName() string {
	return g.repoName
}

func (g *ghPullRequest) GetVCSID() vcs.VCSObjectID {
	return vcs.NewVCSObjectInt(g.pr.GetID())
}

func (g *ghPullRequest) GetPRNumber() int64 {
	return g.pr.GetID()
}

func (g *ghPullRequest) GetBaseSHA() string {
	return g.pr.GetBase().GetSHA()
}

func (g *ghPullRequest) GetHeadSHA() string {
	return g.pr.GetHead().GetSHA()
}

func (g *ghPullRequest) GetBaseBranch() string {
	return g.pr.GetBase().GetRef()
}

func (g *ghPullRequest) GetHeadBranch() string {
	return g.pr.GetHead().GetRef()
}

func (g *ghPullRequest) GetTitle() string {
	return g.pr.GetTitle()
}

func (g *ghPullRequest) GetState() string {
	return g.pr.GetState()
}
