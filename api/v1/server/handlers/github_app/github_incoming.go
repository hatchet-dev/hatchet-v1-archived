package github_app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs/github"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/monitors"
	"github.com/hatchet-dev/hatchet/internal/queuemanager"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/runmanager"
	"github.com/hatchet-dev/hatchet/internal/runutils"
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"

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
	}
}

func (g *GithubIncomingWebhookHandler) processPullRequestEvent(team *models.Team, event *githubsdk.PullRequestEvent, r *http.Request) error {
	// convert event to a vcs.VCSRepositoryPullRequest
	vcsPR := github.ToVCSRepositoryPullRequest(*event.GetRepo().GetOwner().Login, event.GetRepo().GetName(), event.GetPullRequest())

	// call create or update on the PR
	// get the VCSRepository from the repo name + owner (without a module)
	mods, err := g.Repo().Module().ListVCSRepositoryModules(team.ID, vcsPR.GetRepoOwner(), vcsPR.GetRepoName())

	if err != nil {
		return err
	}

	// if there are no modules, continue
	if len(mods) == 0 {
		return nil
	}

	for _, mod := range mods {
		vcsRepo, err := vcs.GetVCSRepositoryFromModule(g.Config().VCSProviders, mod)

		if err != nil {
			err = multierror.Append(err)
			continue
		}

		err = vcsRepo.CreateOrUpdatePRInDatabase(mod.TeamID, vcsPR)

		if err != nil {
			err = multierror.Append(err)
			continue
		}

		// case on the event action
		switch *event.Action {
		case "opened", "reopened", "synchronize":
			return g.VCS_processPullRequestOpened(team, mod, vcsRepo, vcsPR)
		case "edited":
			return g.VCS_processPullRequestEdited(team, mod, vcsRepo, vcsPR)
		case "closed":
			if event.GetPullRequest().GetMerged() {
				return g.VCS_processPullRequestMerged(team, mod, vcsRepo, vcsPR)
			} else {
				return g.VCS_processPullRequestEdited(team, mod, vcsRepo, vcsPR)
			}
		}
	}

	return nil
}

// TODO: anything in "VCS" should be moved to the VCS package. These are shared methods.
func (g *GithubIncomingWebhookHandler) VCS_processPullRequestMerged(
	team *models.Team,
	mod *models.Module,
	vcsRepo vcs.VCSRepository,
	pr vcs.VCSRepositoryPullRequest,
) error {
	// check if we should actually process this module
	shouldTrigger, msg, err := g.VCS_shouldTrigger(
		mod,
		models.ModuleRunKindApply,
		pr,
		vcsRepo,
	)

	if err != nil {
		return err
	} else if !shouldTrigger {
		g.Config().Logger.Debug().Msgf("did not trigger a module run: %s", msg)

		return nil
	}

	run := &models.ModuleRun{
		ModuleID: mod.ID,
		Status:   models.ModuleRunStatusQueued,
		Kind:     models.ModuleRunKindApply,
		ModuleRunConfig: models.ModuleRunConfig{
			TriggerKind:            models.ModuleRunTriggerKindVCS,
			GitCommitSHA:           pr.GetHeadSHA(),
			ModuleValuesVersionID:  mod.CurrentModuleValuesVersionID,
			ModuleEnvVarsVersionID: mod.CurrentModuleEnvVarsVersionID,
		},
		LogLocation: g.Config().DefaultLogStore.GetID(),
	}

	desc, err := runutils.GenerateRunDescription(g.Config(), mod, run, run.Status, nil)

	if err != nil {
		return err
	}

	run.StatusDescription = desc

	run, err = g.Repo().Module().CreateModuleRun(run)

	if err != nil {
		return err
	}

	// get all monitors for this run
	runMonitors, err := monitors.GetAllMonitorsForModuleRun(g.Repo(), mod.TeamID, run)

	if err != nil {
		return err
	}

	run, err = g.Repo().Module().AppendModuleRunMonitors(run, runMonitors)

	if err != nil {
		return err
	}

	err = g.Config().ModuleRunQueueManager.Enqueue(mod, run, &queuemanager.LockOpts{
		LockID:   pr.GetHeadBranch(),
		LockKind: models.ModuleLockKindVCSBranch,
	})

	if err != nil {
		return err
	}

	err = dispatcher.DispatchModuleRunQueueChecker(g.Config().TemporalClient, &modulequeuechecker.CheckQueueInput{
		TeamID:   mod.TeamID,
		ModuleID: mod.ID,
	})

	if err != nil {
		return err
	}

	return err
}

func (g *GithubIncomingWebhookHandler) VCS_processPullRequestOpened(
	team *models.Team,
	mod *models.Module,
	vcsRepo vcs.VCSRepository,
	pr vcs.VCSRepositoryPullRequest,
) error {
	return g.VCS_newPlanFromPR(team, mod, vcsRepo, pr)
}

func (g *GithubIncomingWebhookHandler) VCS_newPlanFromPR(
	team *models.Team,
	mod *models.Module,
	vcsRepo vcs.VCSRepository,
	pr vcs.VCSRepositoryPullRequest,
) error {
	commentBody := "## Hatchet Plan\nRunning `terraform plan`..."

	// check if there's an existing plan for that specific commit SHA. If so, don't queue another run
	planKind := models.ModuleRunKindPlan
	existingRun, _ := g.Repo().Module().ListModuleRunsByVCSSHA(mod.ID, pr.GetHeadSHA(), &planKind)

	if existingRun != nil && len(existingRun) > 0 {
		return nil
	}

	// check if we should actually process this module
	shouldTrigger, msg, err := g.VCS_shouldTrigger(
		mod,
		models.ModuleRunKindApply,
		pr,
		vcsRepo,
	)

	if err != nil {
		return err
	} else if !shouldTrigger {
		g.Config().Logger.Debug().Msgf("did not trigger a module run: %s", msg)

		return nil
	}

	// check if module lock is held by a different module run
	locked := mod.LockID != "" && (mod.LockKind != models.ModuleLockKindVCSBranch || mod.LockID != pr.GetHeadBranch())

	if locked {
		commentBody = "## Hatchet Plan\nLock is currently held by a different PR. Queued..."
	}

	vcsCheck, err := vcsRepo.CreateCheckRun(pr, mod, vcs.VCSCheckRun{
		Name:   fmt.Sprintf("Hatchet plan for %s", mod.DeploymentConfig.ModulePath),
		Status: vcs.VCSCheckRunStatusInProgress,
	})

	if err != nil {
		return fmt.Errorf("error creating new github check run for owner: %s repo %s prNumber: %d. Error: %w",
			pr.GetRepoOwner(), pr.GetRepoName(), pr.GetPRNumber(), err)
	}

	vcsComment, err := vcsRepo.CreateOrUpdateComment(pr, mod, nil, commentBody)

	if err != nil {
		return fmt.Errorf("error creating new github comment for owner: %s repo %s prNumber: %d. Error: %w",
			pr.GetRepoOwner(), pr.GetRepoName(), pr.GetPRNumber(), err)
	}

	status := models.ModuleRunStatusQueued

	run := &models.ModuleRun{
		ModuleID:    mod.ID,
		Status:      status,
		Kind:        models.ModuleRunKindPlan,
		LogLocation: g.Config().DefaultLogStore.GetID(),
		ModuleRunConfig: models.ModuleRunConfig{
			TriggerKind: models.ModuleRunTriggerKindVCS,
			// TODO(abelanger5): store check id and comment id as part of run
			// GithubCheckID:          checkResp.GetID(),
			// GithubCommentID:        commentResp.GetID(),
			GitCommitSHA:           pr.GetHeadSHA(),
			ModuleValuesVersionID:  mod.CurrentModuleValuesVersionID,
			ModuleEnvVarsVersionID: mod.CurrentModuleEnvVarsVersionID,
		},
	}

	vcsRepo.PopulateModuleRun(run, pr.GetVCSID(), vcsCheck, vcsComment)

	desc, err := runutils.GenerateRunDescription(g.Config(), mod, run, run.Status, nil)

	if err != nil {
		return err
	}

	run.StatusDescription = desc

	run, err = g.Repo().Module().CreateModuleRun(run)

	if err != nil {
		return err
	}

	// get all monitors for this run
	runMonitors, err := monitors.GetAllMonitorsForModuleRun(g.Repo(), mod.TeamID, run)

	if err != nil {
		return err
	}

	run, err = g.Repo().Module().AppendModuleRunMonitors(run, runMonitors)

	if err != nil {
		return err
	}

	err = g.Config().ModuleRunQueueManager.Enqueue(mod, run, &queuemanager.LockOpts{
		LockID:   pr.GetHeadBranch(),
		LockKind: models.ModuleLockKindVCSBranch,
	})

	if err != nil {
		return err
	}

	return dispatcher.DispatchModuleRunQueueChecker(g.Config().TemporalClient, &modulequeuechecker.CheckQueueInput{
		TeamID:   mod.TeamID,
		ModuleID: mod.ID,
	})
}

func (g *GithubIncomingWebhookHandler) VCS_shouldTrigger(
	mod *models.Module,
	kind models.ModuleRunKind,
	pr vcs.VCSRepositoryPullRequest,
	vcsRepository vcs.VCSRepository,
) (bool, string, error) {
	triggerInput := &runmanager.TriggerInput{
		BaseBranch: pr.GetBaseBranch(),
	}

	// get files for pull request if this is a plan
	if kind == models.ModuleRunKindPlan {
		commitsRes, err := vcsRepository.CompareCommits(
			pr.GetBaseSHA(), pr.GetHeadSHA(),
		)

		if err != nil {
			return false, "", err
		}

		fileNames := make([]string, 0)

		for _, file := range commitsRes.GetFiles() {
			fileNames = append(fileNames, file.GetFilename())
		}

		triggerInput.Files = fileNames
	}

	res, msg := runmanager.Trigger(mod, kind, triggerInput)

	return res, msg, nil
}

func (g *GithubIncomingWebhookHandler) VCS_processPullRequestEdited(
	team *models.Team,
	mod *models.Module,
	vcsRepo vcs.VCSRepository,
	pr vcs.VCSRepositoryPullRequest,
) error {
	// if the pr has been closed, determine if the head branch holds the lock on
	// any module. if so, remove the lock
	if pr.GetState() == "closed" {
		err := g.Config().ModuleRunQueueManager.FlushQueue(mod, &queuemanager.LockOpts{
			LockID:   pr.GetHeadBranch(),
			LockKind: models.ModuleLockKindVCSBranch,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
