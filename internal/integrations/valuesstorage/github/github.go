package github

import (
	"context"
	"encoding/json"
	"fmt"

	githubsdk "github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/internal/models"
)

// GithubValuesStore represents a values storage manager which uses a path in the Github
// repository to reference values
type GithubValuesStore struct {
	client    *githubsdk.Client
	githubRef string
}

func NewGithubValuesStore(client *githubsdk.Client, githubRef string) *GithubValuesStore {
	return &GithubValuesStore{client, githubRef}
}

func (d *GithubValuesStore) WriteValues(mvv *models.ModuleValuesVersion, values map[string]interface{}) error {
	return fmt.Errorf("github-based values storage does not support writing values")
}

func (d *GithubValuesStore) ReadValues(mvv *models.ModuleValuesVersion) (map[string]interface{}, error) {
	file, _, err := d.client.Repositories.DownloadContents(
		context.Background(),
		mvv.GithubRepoOwner,
		mvv.GithubRepoName,
		mvv.GithubValuesPath,
		&githubsdk.RepositoryContentGetOptions{
			Ref: d.githubRef,
		},
	)

	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})

	err = json.NewDecoder(file).Decode(&res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
