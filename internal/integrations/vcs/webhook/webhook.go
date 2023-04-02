package webhook

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/monitors"
	"github.com/hatchet-dev/hatchet/internal/queuemanager"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/runmanager"
	"github.com/hatchet-dev/hatchet/internal/runutils"
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
)

type WebhookHandler struct {
	repo   repository.Repository
	config *server.Config
}

func NewWebhookHandler(repo repository.Repository, config *server.Config) *WebhookHandler {
	return &WebhookHandler{
		repo:   repo,
		config: config,
	}
}

func (w *WebhookHandler) ProcessPullRequestMerged(
	team *models.Team,
	mod *models.Module,
	vcsRepo vcs.VCSRepository,
	pr vcs.VCSRepositoryPullRequest,
) error {
	// check if we should actually process this module
	shouldTrigger, msg, err := w.ShouldTrigger(
		mod,
		models.ModuleRunKindApply,
		pr,
		vcsRepo,
	)

	if err != nil {
		return err
	} else if !shouldTrigger {
		w.config.Logger.Debug().Msgf("did not trigger a module run: %s", msg)

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
		LogLocation: w.config.DefaultLogStore.GetID(),
	}

	desc, err := runutils.GenerateRunDescription(w.config, mod, run, run.Status, nil)

	if err != nil {
		return err
	}

	run.StatusDescription = desc

	run, err = w.repo.Module().CreateModuleRun(run)

	if err != nil {
		return err
	}

	// get all monitors for this run
	runMonitors, err := monitors.GetAllMonitorsForModuleRun(w.repo, mod.TeamID, run)

	if err != nil {
		return err
	}

	run, err = w.repo.Module().AppendModuleRunMonitors(run, runMonitors)

	if err != nil {
		return err
	}

	err = w.config.ModuleRunQueueManager.Enqueue(mod, run, &queuemanager.LockOpts{
		LockID:   pr.GetHeadBranch(),
		LockKind: models.ModuleLockKindVCSBranch,
	})

	if err != nil {
		return err
	}

	err = dispatcher.DispatchModuleRunQueueChecker(w.config.TemporalClient, &modulequeuechecker.CheckQueueInput{
		TeamID:   mod.TeamID,
		ModuleID: mod.ID,
	})

	if err != nil {
		return err
	}

	return err
}

func (w *WebhookHandler) ProcessPullRequestOpened(
	team *models.Team,
	mod *models.Module,
	vcsRepo vcs.VCSRepository,
	pr vcs.VCSRepositoryPullRequest,
) error {
	return w.NewPlanFromPR(team, mod, vcsRepo, pr)
}

func (w *WebhookHandler) NewPlanFromPR(
	team *models.Team,
	mod *models.Module,
	vcsRepo vcs.VCSRepository,
	pr vcs.VCSRepositoryPullRequest,
) error {
	commentBody := "## Hatchet Plan\nRunning `terraform plan`..."

	// check if there's an existing plan for that specific commit SHA. If so, don't queue another run
	planKind := models.ModuleRunKindPlan
	existingRun, _ := w.repo.Module().ListModuleRunsByVCSSHA(mod.ID, pr.GetHeadSHA(), &planKind)

	if existingRun != nil && len(existingRun) > 0 {
		return nil
	}

	// check if we should actually process this module
	shouldTrigger, msg, err := w.ShouldTrigger(
		mod,
		models.ModuleRunKindApply,
		pr,
		vcsRepo,
	)

	if err != nil {
		return err
	} else if !shouldTrigger {
		w.config.Logger.Debug().Msgf("did not trigger a module run: %s", msg)

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
		LogLocation: w.config.DefaultLogStore.GetID(),
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

	desc, err := runutils.GenerateRunDescription(w.config, mod, run, run.Status, nil)

	if err != nil {
		return err
	}

	run.StatusDescription = desc

	run, err = w.repo.Module().CreateModuleRun(run)

	if err != nil {
		return err
	}

	// get all monitors for this run
	runMonitors, err := monitors.GetAllMonitorsForModuleRun(w.repo, mod.TeamID, run)

	if err != nil {
		return err
	}

	run, err = w.repo.Module().AppendModuleRunMonitors(run, runMonitors)

	if err != nil {
		return err
	}

	err = w.config.ModuleRunQueueManager.Enqueue(mod, run, &queuemanager.LockOpts{
		LockID:   pr.GetHeadBranch(),
		LockKind: models.ModuleLockKindVCSBranch,
	})

	if err != nil {
		return err
	}

	return dispatcher.DispatchModuleRunQueueChecker(w.config.TemporalClient, &modulequeuechecker.CheckQueueInput{
		TeamID:   mod.TeamID,
		ModuleID: mod.ID,
	})
}

func (w *WebhookHandler) ShouldTrigger(
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

func (w *WebhookHandler) ProcessPullRequestEdited(
	team *models.Team,
	mod *models.Module,
	vcsRepo vcs.VCSRepository,
	pr vcs.VCSRepositoryPullRequest,
) error {
	// if the pr has been closed, determine if the head branch holds the lock on
	// any module. if so, remove the lock
	if pr.GetState() == "closed" {
		err := w.config.ModuleRunQueueManager.FlushQueue(mod, &queuemanager.LockOpts{
			LockID:   pr.GetHeadBranch(),
			LockKind: models.ModuleLockKindVCSBranch,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
