package gorm

import (
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"gorm.io/gorm"
)

// GithubPullRequestRepository uses gorm.DB for querying the database
type GithubPullRequestRepository struct {
	db *gorm.DB
}

// NewGithubPullRequestRepository returns a DefaultGithubPullRequestRepository which uses
// gorm.DB for querying the database
func NewGithubPullRequestRepository(db *gorm.DB) repository.GithubPullRequestRepository {
	return &GithubPullRequestRepository{db}
}

func (repo *GithubPullRequestRepository) CreateGithubPullRequest(gpr *models.GithubPullRequest) (*models.GithubPullRequest, repository.RepositoryError) {
	if err := repo.db.Create(gpr).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return gpr, nil
}

func (repo *GithubPullRequestRepository) ReadGithubPullRequestByID(teamID, id string) (*models.GithubPullRequest, repository.RepositoryError) {
	gpr := &models.GithubPullRequest{}

	if err := repo.db.Where("team_id = ? AND id = ?", teamID, id).First(&gpr).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return gpr, nil
}

func (repo *GithubPullRequestRepository) ReadGithubPullRequestByGithubID(teamID string, ghID int64) (*models.GithubPullRequest, repository.RepositoryError) {
	gpr := &models.GithubPullRequest{}

	if err := repo.db.Where("team_id = ? AND github_pull_request_id = ?", teamID, ghID).First(&gpr).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return gpr, nil
}

func (repo *GithubPullRequestRepository) ListGithubPullRequestsByHeadBranch(teamID, repoOwner, repoName, branchName string) ([]*models.GithubPullRequest, repository.RepositoryError) {
	res := make([]*models.GithubPullRequest, 0)

	if err := repo.db.Where("team_id = ? AND github_repository_owner = ? AND github_repository_name = ? AND github_pull_request_head_branch = ? AND github_pull_request_state = ?", teamID, repoOwner, repoName, branchName, "open").Find(&res).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return res, nil
}

func (repo *GithubPullRequestRepository) UpdateGithubPullRequest(
	gpr *models.GithubPullRequest,
) (*models.GithubPullRequest, repository.RepositoryError) {
	if err := repo.db.Save(gpr).Error; err != nil {
		return nil, err
	}

	return gpr, nil
}

func (repo *GithubPullRequestRepository) DeleteGithubPullRequest(gpr *models.GithubPullRequest) (*models.GithubPullRequest, repository.RepositoryError) {
	if err := repo.db.Delete(gpr).Error; err != nil {
		return nil, err
	}

	return gpr, nil
}

func (repo *GithubPullRequestRepository) CreateGithubPullRequestComment(gc *models.GithubPullRequestComment) (*models.GithubPullRequestComment, repository.RepositoryError) {
	if err := repo.db.Create(gc).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return gc, nil
}

func (repo *GithubPullRequestRepository) ReadGithubPullRequestCommentByID(moduleID, id string) (*models.GithubPullRequestComment, repository.RepositoryError) {
	gc := &models.GithubPullRequestComment{}

	if err := repo.db.Where("module_id = ? AND id = ?", moduleID, id).First(&gc).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return gc, nil
}

func (repo *GithubPullRequestRepository) ReadGithubPullRequestCommentByGithubID(moduleID string, ghID int64) (*models.GithubPullRequestComment, repository.RepositoryError) {
	gc := &models.GithubPullRequestComment{}

	if err := repo.db.Where("module_id = ? AND github_comment_id = ?", moduleID, ghID).First(&gc).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return gc, nil
}

func (repo *GithubPullRequestRepository) UpdateGithubPullRequestComment(
	gc *models.GithubPullRequestComment,
) (*models.GithubPullRequestComment, repository.RepositoryError) {
	if err := repo.db.Save(gc).Error; err != nil {
		return nil, err
	}

	return gc, nil
}

func (repo *GithubPullRequestRepository) DeleteGithubPullRequestComment(gc *models.GithubPullRequestComment) (*models.GithubPullRequestComment, repository.RepositoryError) {
	if err := repo.db.Delete(gc).Error; err != nil {
		return nil, err
	}

	return gc, nil
}
