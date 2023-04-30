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
	"github.com/hatchet-dev/hatchet/internal/integrations/vcs/webhook"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"

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

	fmt.Println("HERE 0")

	if err != nil {
		if errors.Is(err, repository.RepositoryErrorNotFound) {
			g.HandleAPIError(w, r, apierrors.NewErrForbidden(
				fmt.Errorf("github webhook with id %s not found in team %s", gwID, team.ID),
			))

			return
		}

		fmt.Println("THIS IS THE ERROR")

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

	webhookHandler := webhook.NewWebhookHandler(g.Repo(), g.Config())

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
			return webhookHandler.ProcessPullRequestOpened(team, mod, vcsRepo, vcsPR)
		case "edited":
			return webhookHandler.ProcessPullRequestEdited(team, mod, vcsRepo, vcsPR)
		case "closed":
			if event.GetPullRequest().GetMerged() {
				return webhookHandler.ProcessPullRequestMerged(team, mod, vcsRepo, vcsPR)
			} else {
				return webhookHandler.ProcessPullRequestEdited(team, mod, vcsRepo, vcsPR)
			}
		}
	}

	return nil
}
