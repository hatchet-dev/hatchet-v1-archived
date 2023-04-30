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

type ListGithubRepoBranchesHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewListGithubRepoBranchesHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &ListGithubRepoBranchesHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (g *ListGithubRepoBranchesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client, reqErr := GetGithubAppClientFromRequest(g.Config(), r)

	if reqErr != nil {
		g.HandleAPIError(w, r, reqErr)
		return
	}

	owner, reqErr := handlerutils.GetURLParamString(r, types.URLParamGithubRepoOwner)

	if reqErr != nil {
		g.HandleAPIError(w, r, reqErr)
		return
	}

	name, reqErr := handlerutils.GetURLParamString(r, types.URLParamGithubRepoName)

	if reqErr != nil {
		g.HandleAPIError(w, r, reqErr)
		return
	}

	repo, _, err := client.Repositories.Get(
		context.TODO(),
		owner,
		name,
	)

	if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	defaultBranch := repo.GetDefaultBranch()

	// List all branches for a specified repo
	allBranches, resp, err := client.Repositories.ListBranches(context.Background(), owner, name, &github.BranchListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})

	if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// make workers to get branches concurrently
	const WCOUNT = 5
	numPages := resp.LastPage + 1
	var workerErr error
	var mu sync.Mutex
	var wg sync.WaitGroup

	worker := func(cp int) {
		defer wg.Done()

		for cp < numPages {
			opts := &github.BranchListOptions{
				ListOptions: github.ListOptions{
					Page:    cp,
					PerPage: 100,
				},
			}

			branches, _, err := client.Repositories.ListBranches(context.Background(), owner, name, opts)

			if err != nil {
				mu.Lock()
				workerErr = err
				mu.Unlock()
				return
			}

			mu.Lock()
			allBranches = append(allBranches, branches...)
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

	res := make(types.ListGithubRepoBranchesResponse, 0)

	for _, branch := range allBranches {
		res = append(res, types.GithubBranch{
			BranchName: *branch.Name,
			IsDefault:  defaultBranch == *branch.Name,
		})
	}

	g.WriteResult(w, r, res)
}
