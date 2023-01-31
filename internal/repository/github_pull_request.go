package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// GithubPullRequestRepository represents the set of queries on the GithubPullRequest model
type GithubPullRequestRepository interface {
	CreateGithubPullRequest(gpr *models.GithubPullRequest) (*models.GithubPullRequest, RepositoryError)
	ReadGithubPullRequestByID(teamID, id string) (*models.GithubPullRequest, RepositoryError)
	ReadGithubPullRequestByGithubID(teamID string, ghID int64) (*models.GithubPullRequest, RepositoryError)
	ListGithubPullRequestsByHeadBranch(teamID, repoOwner, repoName, branchName string) ([]*models.GithubPullRequest, RepositoryError)
	UpdateGithubPullRequest(gpr *models.GithubPullRequest) (*models.GithubPullRequest, RepositoryError)
	DeleteGithubPullRequest(gpr *models.GithubPullRequest) (*models.GithubPullRequest, RepositoryError)

	CreateGithubPullRequestComment(gc *models.GithubPullRequestComment) (*models.GithubPullRequestComment, RepositoryError)
	ReadGithubPullRequestCommentByID(moduleID, id string) (*models.GithubPullRequestComment, RepositoryError)
	ReadGithubPullRequestCommentByGithubID(moduleID string, ghID int64) (*models.GithubPullRequestComment, RepositoryError)
	UpdateGithubPullRequestComment(gc *models.GithubPullRequestComment) (*models.GithubPullRequestComment, RepositoryError)
	DeleteGithubPullRequestComment(gc *models.GithubPullRequestComment) (*models.GithubPullRequestComment, RepositoryError)
}
