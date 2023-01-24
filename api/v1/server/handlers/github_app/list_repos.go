package github_app

import (
	"context"
	"net/http"
	"sync"

	"github.com/google/go-github/v49/github"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

type ListGithubReposHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewListGithubReposHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ListGithubReposHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (g *ListGithubReposHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client, err := GetGithubAppClientFromRequest(g.Config(), r)

	if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// figure out number of repositories
	opt := &github.ListOptions{
		PerPage: 100,
	}

	repoList, resp, err := client.Apps.ListRepos(context.Background(), opt)

	if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	allRepos := repoList.Repositories

	// make workers to get pages concurrently
	const WCOUNT = 5
	numPages := resp.LastPage + 1
	var workerErr error
	var mu sync.Mutex
	var wg sync.WaitGroup

	worker := func(cp int) {
		defer wg.Done()

		for cp < numPages {
			cur_opt := &github.ListOptions{
				Page:    cp,
				PerPage: 100,
			}

			repos, _, err := client.Apps.ListRepos(context.Background(), cur_opt)

			if err != nil {
				mu.Lock()
				workerErr = err
				mu.Unlock()
				return
			}

			mu.Lock()
			allRepos = append(allRepos, repos.Repositories...)
			mu.Unlock()

			cp += WCOUNT
		}
	}

	var numJobs int
	if numPages > WCOUNT {
		numJobs = WCOUNT
	} else {
		numJobs = numPages
	}

	wg.Add(numJobs)

	// page 1 is already loaded so we start with 2
	for i := 1; i <= numJobs; i++ {
		go worker(i + 1)
	}

	wg.Wait()

	if workerErr != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	res := make(types.ListGithubReposResponse, 0)

	for _, repo := range allRepos {
		res = append(res, types.GithubRepo{
			RepoName:  *repo.Name,
			RepoOwner: *repo.Owner.Login,
		})
	}

	g.WriteResult(w, r, res)
}
