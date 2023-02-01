package teams

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/git/github"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/provisioner"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/runmanager"

	githubsdk "github.com/google/go-github/v49/github"
)

type GithubIncomingWebhookHandler struct {
	handlers.HatchetHandlerReadWriter
}

func NewGithubIncomingWebhookHandler(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler {
	return &GithubIncomingWebhookHandler{
		HatchetHandlerReadWriter: handlers.NewDefaultHatchetHandler(config, decoderValidator, writer),
	}
}

func (g *GithubIncomingWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	teamID, reqErr := handlerutils.GetURLParamString(r, types.URLParamTeamID)

	if reqErr != nil {
		g.HandleAPIError(w, r, reqErr)
		return
	}

	team, err := g.Repo().Team().ReadTeamByID(teamID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			g.HandleAPIError(w, r, apierrors.NewErrForbidden(
				fmt.Errorf("team with id %s not found", teamID),
			))

			return
		}

		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	gwID, reqErr := handlerutils.GetURLParamString(r, types.URLParamGithubWebhookID)

	if reqErr != nil {
		g.HandleAPIError(w, r, reqErr)
		return
	}

	gw, err := g.Repo().GithubWebhook().ReadGithubWebhookByID(team.ID, gwID)

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			g.HandleAPIError(w, r, apierrors.NewErrForbidden(
				fmt.Errorf("github webhook with id %s not found in team %s", gwID, team.ID),
			))

			return
		}

		g.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// validate the payload using the github webhook signing secret
	payload, err := githubsdk.ValidatePayload(r, gw.SigningSecret)

	if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrForbidden(fmt.Errorf("error validating webhook payload: %w", err)))
		return
	}

	event, err := githubsdk.ParseWebHook(githubsdk.WebHookType(r), payload)

	if err != nil {
		g.HandleAPIError(w, r, apierrors.NewErrInternal(fmt.Errorf("error parsing webhook: %w", err)))
		return
	}

	switch event := event.(type) {
	case *githubsdk.PullRequestEvent:
		err = g.processPullRequestEvent(team, event, r)

		if err != nil {
			g.HandleAPIError(w, r, apierrors.NewErrInternal(fmt.Errorf("error processing pull request webhook event: %w", err)))
			return
		}
	case *githubsdk.PushEvent:
		err = g.processPushEvent(team, event, r)

		if err != nil {
			g.HandleAPIError(w, r, apierrors.NewErrInternal(fmt.Errorf("error processing push webhook event: %w", err)))
			return
		}
	}
}

func (g *GithubIncomingWebhookHandler) processPullRequestEvent(team *models.Team, event *githubsdk.PullRequestEvent, r *http.Request) error {
	// case on the event action
	switch *event.Action {
	case "opened", "reopened":
		return g.processPullRequestOpened(team, event)
	case "edited":
		return g.processPullRequestEdited(team, event)
	case "closed":
		if event.GetPullRequest().GetMerged() {
			return g.processPullRequestMerged(team, event)
		} else {
			return g.processPullRequestEdited(team, event)
		}
	}

	return nil
}

func (g *GithubIncomingWebhookHandler) processPushEvent(team *models.Team, event *githubsdk.PushEvent, r *http.Request) error {
	owner := event.GetRepo().GetOwner().GetLogin()
	repoName := event.GetRepo().GetName()
	baseSHA := event.GetBefore()
	headSHA := event.GetHeadCommit().GetID()

	headBranch := strings.TrimPrefix(event.GetRef(), "refs/heads/")

	// determine all modules that should trigger based on this PR
	mods, err := g.Repo().Module().ListGithubRepositoryModules(team.ID, owner, repoName)

	if err != nil {
		return err
	}

	// if there are no modules, continue
	if len(mods) == 0 {
		return nil
	}

	// find pull request that corresponds to this head branch
	prs, err := g.Repo().GithubPullRequest().ListGithubPullRequestsByHeadBranch(team.ID, owner, repoName, headBranch)

	if err != nil {
		return err
	}

	ghClients := make(map[string]*githubsdk.Client)

	errs := make([]string, 0)

	for _, mod := range mods {
		for _, pr := range prs {
			err := g.newPlanFromPR(team, mod, pr, &eventData{owner, repoName, baseSHA, headSHA, ""}, ghClients)

			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ", "))
	}

	return nil
}

func (g *GithubIncomingWebhookHandler) processPullRequestOpened(team *models.Team, event *githubsdk.PullRequestEvent) error {
	owner := event.GetRepo().GetOwner().GetLogin()
	repoName := event.GetRepo().GetName()
	baseSHA := event.GetPullRequest().GetBase().GetSHA()
	headSHA := event.GetPullRequest().GetHead().GetSHA()
	branch := event.GetPullRequest().GetBase().GetRef()

	// determine all modules that should trigger based on this PR
	mods, err := g.Repo().Module().ListGithubRepositoryModules(team.ID, owner, repoName)

	if err != nil {
		return err
	}

	// if there are no modules, return nil
	if len(mods) == 0 {
		return nil
	}

	// determine if pull request already exists
	ghPR, err := g.Repo().GithubPullRequest().ReadGithubPullRequestByGithubID(team.ID, event.GetPullRequest().GetID())

	if err != nil && !errors.Is(err, repository.RepositoryErrorNotFound) {
		return err
	}

	if ghPR == nil {
		// create a PR object in the db
		ghPR, err = g.Repo().GithubPullRequest().CreateGithubPullRequest(&models.GithubPullRequest{
			TeamID:                      team.ID,
			GithubRepositoryOwner:       owner,
			GithubRepositoryName:        repoName,
			GithubPullRequestID:         event.GetPullRequest().GetID(),
			GithubPullRequestTitle:      event.GetPullRequest().GetTitle(),
			GithubPullRequestNumber:     int64(event.GetPullRequest().GetNumber()),
			GithubPullRequestHeadBranch: event.GetPullRequest().GetHead().GetRef(),
			GithubPullRequestBaseBranch: event.GetPullRequest().GetBase().GetRef(),
			GithubPullRequestState:      event.GetPullRequest().GetState(),
		})

		if err != nil {
			return err
		}
	} else {
		ghPR.GithubPullRequestTitle = event.GetPullRequest().GetTitle()
		ghPR.GithubPullRequestHeadBranch = event.GetPullRequest().GetHead().GetRef()
		ghPR.GithubPullRequestBaseBranch = event.GetPullRequest().GetBase().GetRef()
		ghPR.GithubPullRequestState = event.GetPullRequest().GetState()

		ghPR, err = g.Repo().GithubPullRequest().UpdateGithubPullRequest(ghPR)

		if err != nil {
			return err
		}
	}

	// list of github app installation ids to clients that can be reused
	ghClients := make(map[string]*githubsdk.Client)

	errs := make([]string, 0)

	// create comments corresponding to each module
	for _, mod := range mods {
		err := g.newPlanFromPR(team, mod, ghPR, &eventData{owner, repoName, baseSHA, headSHA, branch}, ghClients)

		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ", "))
	}

	return nil
}

func (g *GithubIncomingWebhookHandler) processPullRequestEdited(team *models.Team, event *githubsdk.PullRequestEvent) error {
	owner := event.GetRepo().GetOwner().GetLogin()
	repoName := event.GetRepo().GetName()
	title := event.GetPullRequest().GetTitle()
	headBranch := event.GetPullRequest().GetHead().GetRef()
	baseBranch := event.GetPullRequest().GetBase().GetRef()
	state := event.GetPullRequest().GetState()

	// read the PRs from the database
	prs, err := g.Repo().GithubPullRequest().ListGithubPullRequestsByHeadBranch(team.ID, owner, repoName, headBranch)

	if err != nil {
		return err
	}

	var errs = make([]string, 0)

	for _, pr := range prs {
		pr.GithubPullRequestTitle = title
		pr.GithubPullRequestHeadBranch = headBranch
		pr.GithubPullRequestBaseBranch = baseBranch
		pr.GithubPullRequestState = state

		pr, err = g.Repo().GithubPullRequest().UpdateGithubPullRequest(pr)

		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ", "))
	}

	return nil
}

func (g *GithubIncomingWebhookHandler) processPullRequestMerged(team *models.Team, event *githubsdk.PullRequestEvent) error {
	owner := event.GetRepo().GetOwner().GetLogin()
	repoName := event.GetRepo().GetName()
	baseSHA := event.GetPullRequest().GetBase().GetSHA()
	baseBranch := event.GetPullRequest().GetBase().GetRef()
	headSHA := event.GetPullRequest().GetHead().GetSHA()
	headBranch := event.GetPullRequest().GetHead().GetRef()

	ghPR, err := g.Repo().GithubPullRequest().ReadGithubPullRequestByGithubID(team.ID, event.GetPullRequest().GetID())

	if err != nil {
		return err
	}

	ghPR.GithubPullRequestTitle = event.GetPullRequest().GetTitle()
	ghPR.GithubPullRequestHeadBranch = event.GetPullRequest().GetHead().GetRef()
	ghPR.GithubPullRequestBaseBranch = event.GetPullRequest().GetBase().GetRef()
	ghPR.GithubPullRequestState = event.GetPullRequest().GetState()

	ghPR, err = g.Repo().GithubPullRequest().UpdateGithubPullRequest(ghPR)

	if err != nil {
		return err
	}

	// determine all modules that should trigger based on this PR
	mods, err := g.Repo().Module().ListGithubRepositoryModules(team.ID, owner, repoName)

	if err != nil {
		return err
	}

	// if there are no modules, continue
	if len(mods) == 0 {
		return nil
	}

	ghClients := make(map[string]*githubsdk.Client)

	errs := make([]error, 0)

	// create comments corresponding to each module
	for _, mod := range mods {
		// TODO: check that there's not an EXISTING APPLY which is IN PROGRESS
		// check if there's an existing run for that specific commit SHA. If so, don't queue another run
		// existingRun, _ := g.Repo().Module().ReadModuleRunByGithubSHA(mod.ID, sha)

		// if existingRun != nil {
		// 	continue
		// }

		// TODO: check if we should actually process this module

		client, ok := ghClients[mod.DeploymentConfig.GithubAppInstallationID]

		if !ok {
			client, err = github.GetGithubAppClientFromModule(g.Config(), mod)

			if err != nil {
				errs = append(errs, err)
				continue
			}

			ghClients[mod.DeploymentConfig.GithubAppInstallationID] = client
		}

		// check if we should actually process this module
		shouldTrigger, msg, err := g.shouldTrigger(
			client, mod,
			models.ModuleRunKindPlan,
			owner,
			repoName,
			baseBranch,
			baseSHA,
			headSHA,
		)

		if err != nil {
			return err
		} else if !shouldTrigger {
			g.Config().Logger.Debug().Msgf("did not trigger a module run: %s", msg)
			continue
		}

		checkResp, _, err := client.Checks.CreateCheckRun(
			context.Background(),
			owner,
			repoName,
			githubsdk.CreateCheckRunOptions{
				Name:    fmt.Sprintf("Hatchet apply for %s", mod.DeploymentConfig.ModulePath),
				HeadSHA: headSHA,
			},
		)

		if err != nil {
			return fmt.Errorf("error creating new github check run for owner: %s repo %s branch: %s. Error: %w",
				owner, repoName, headBranch, err)
		}

		run := &models.ModuleRun{
			ModuleID: mod.ID,
			Status:   models.ModuleRunStatusInProgress,
			Kind:     models.ModuleRunKindApply,
			ModuleRunConfig: models.ModuleRunConfig{
				TriggerKind:     models.ModuleRunTriggerKindGithub,
				GithubCheckID:   checkResp.GetID(),
				GithubCommitSHA: headSHA,
			},
		}

		run, err = g.Repo().Module().CreateModuleRun(run)

		if err != nil {
			return err
		}

		err = g.Config().DefaultProvisioner.RunApply(&provisioner.ProvisionOpts{
			Team:       team,
			Module:     mod,
			ModuleRun:  run,
			TokenOpts:  *g.Config().TokenOpts,
			Repository: g.Repo(),
			ServerURL:  g.Config().ServerRuntimeConfig.ServerURL,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

type eventData struct {
	repoOwner, repoName, baseSHA, headSHA, branch string
}

func (g *GithubIncomingWebhookHandler) newPlanFromPR(
	team *models.Team,
	mod *models.Module,
	ghPR *models.GithubPullRequest,
	eventData *eventData,
	ghClients map[string]*githubsdk.Client,
) error {
	commentBody := "## Hatchet Plan\nRunning `terraform plan`..."

	// check if there's an existing run for that specific commit SHA. If so, don't queue another run
	existingRun, _ := g.Repo().Module().ReadModuleRunByGithubSHA(mod.ID, eventData.headSHA)

	if existingRun != nil {
		return nil
	}

	client, ok := ghClients[mod.DeploymentConfig.GithubAppInstallationID]
	var err error

	if !ok {
		client, err = github.GetGithubAppClientFromModule(g.Config(), mod)

		if err != nil {
			return err
		}

		ghClients[mod.DeploymentConfig.GithubAppInstallationID] = client
	}

	// check if we should actually process this module
	shouldTrigger, msg, err := g.shouldTrigger(
		client, mod,
		models.ModuleRunKindPlan,
		eventData.repoOwner,
		eventData.repoName,
		eventData.branch,
		eventData.baseSHA,
		eventData.headSHA,
	)

	if err != nil {
		return err
	} else if !shouldTrigger {
		g.Config().Logger.Debug().Msgf("did not trigger a module run: %s", msg)
		return nil
	}

	checkResp, _, err := client.Checks.CreateCheckRun(
		context.Background(),
		eventData.repoOwner,
		eventData.repoName,
		githubsdk.CreateCheckRunOptions{
			Name:    fmt.Sprintf("Hatchet plan for %s", mod.DeploymentConfig.ModulePath),
			HeadSHA: eventData.headSHA,
		},
	)

	if err != nil {
		return fmt.Errorf("error creating new github check run for owner: %s repo %s prNumber: %d. Error: %w",
			eventData.repoOwner, eventData.repoName, ghPR.GithubPullRequestNumber, err)
	}

	commentResp, _, err := client.Issues.CreateComment(
		context.Background(),
		eventData.repoOwner,
		eventData.repoName,
		int(ghPR.GithubPullRequestNumber),
		&githubsdk.IssueComment{
			Body: &commentBody,
		},
	)

	if err != nil {
		return fmt.Errorf("error creating new github comment for owner: %s repo %s prNumber: %d. Error: %w",
			eventData.repoOwner, eventData.repoName, ghPR.GithubPullRequestNumber, err)
	}

	// create comment in database
	_, err = g.Repo().GithubPullRequest().CreateGithubPullRequestComment(&models.GithubPullRequestComment{
		GithubPullRequestID: ghPR.ID,
		ModuleID:            mod.ID,
		GithubCommentID:     *commentResp.ID,
	})

	if err != nil {
		return fmt.Errorf("error saving github comment for owner: %s repo %s prNumber: %d. Error: %w",
			eventData.repoOwner, eventData.repoName, ghPR.GithubPullRequestNumber, err)
	}

	run := &models.ModuleRun{
		ModuleID: mod.ID,
		Status:   models.ModuleRunStatusInProgress,
		Kind:     models.ModuleRunKindPlan,
		ModuleRunConfig: models.ModuleRunConfig{
			TriggerKind:     models.ModuleRunTriggerKindGithub,
			GithubCheckID:   checkResp.GetID(),
			GithubCommentID: commentResp.GetID(),
			GithubCommitSHA: eventData.headSHA,
		},
	}

	run, err = g.Repo().Module().CreateModuleRun(run)

	if err != nil {
		return err
	}

	err = g.Config().DefaultProvisioner.RunPlan(&provisioner.ProvisionOpts{
		Team:       team,
		Module:     mod,
		ModuleRun:  run,
		TokenOpts:  *g.Config().TokenOpts,
		Repository: g.Repo(),
		ServerURL:  g.Config().ServerRuntimeConfig.ServerURL,
	})

	if err != nil {
		return err
	}

	return nil
}

func (g *GithubIncomingWebhookHandler) shouldTrigger(
	client *githubsdk.Client,
	mod *models.Module,
	kind models.ModuleRunKind,
	repoOwner, repoName, baseBranch, baseCommit, headCommit string,
) (bool, string, error) {
	// get files for pull request
	commitsRes, _, err := client.Repositories.CompareCommits(
		context.Background(),
		repoOwner,
		repoName,
		baseCommit,
		headCommit,
		&githubsdk.ListOptions{},
	)

	if err != nil {
		return false, "", err
	}

	fileNames := make([]string, 0)

	for _, file := range commitsRes.Files {
		fileNames = append(fileNames, file.GetFilename())
	}

	res, msg := runmanager.Trigger(mod, kind, &runmanager.TriggerInput{
		BaseBranch: baseBranch,
		Files:      fileNames,
	})

	return res, msg, nil
}
