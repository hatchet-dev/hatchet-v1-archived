package github

import (
	githubsdk "github.com/google/go-github/v49/github"
)

type GithubBranch struct {
	*githubsdk.Branch
}

func (g *GithubBranch) GetName() string {
	return g.GetName()
}

func (g *GithubBranch) GetLatestRef() string {
	return g.GetCommit().GetSHA()
}
