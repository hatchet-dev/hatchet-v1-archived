package models

// GithubPullRequest contains data about a Github PR
type GithubPullRequest struct {
	Base

	TeamID string

	GithubRepositoryOwner       string
	GithubRepositoryName        string
	GithubPullRequestID         int64
	GithubPullRequestTitle      string
	GithubPullRequestNumber     int64
	GithubPullRequestHeadBranch string
	GithubPullRequestBaseBranch string
	GithubPullRequestState      string

	GithubPullRequestComments []GithubPullRequestComment
}

// GithubPullRequestComment are identified by their parent pull request along with
// a parent module ID. That is, all modules that are triggered by this PR will have
// their own comment
type GithubPullRequestComment struct {
	Base

	GithubPullRequestID string
	ModuleID            string

	GithubCommentID int64
}
