package logflusher

import (
	"fmt"
	"time"

	"github.com/hatchet-dev/hatchet/internal/repository"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type FlushInput struct {
	TeamID      string
	ModuleID    string
	ModuleRunID string
}

type FlushLogsInput struct {
}

func (lf *LogFlusher) FlushLogs(ctx workflow.Context, input FlushLogsInput) (string, error) {
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    100 * time.Second,
		MaximumAttempts:    0,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retrypolicy,
	}

	// find all module runs with the given default log location that are completed
	runs, _, err := lf.repo.Module().ListCompletedModuleRunsByLogLocation(lf.ls.GetID(), repository.WithLimit(50))

	if err != nil {
		return "", fmt.Errorf("could not list module runs: %s", err.Error())
	}

	// TODO: paginate
	for _, run := range runs {
		// stream logs into the file
		if run.TeamID != "" {
			ctx = workflow.WithActivityOptions(ctx, options)

			var flushOutput string

			flushErr := workflow.ExecuteActivity(ctx, lf.Flush, FlushInput{
				TeamID:      run.TeamID,
				ModuleID:    run.ModuleID,
				ModuleRunID: run.ID,
			}).Get(ctx, &flushOutput)

			if flushErr != nil {
				return "", flushErr
			}
		}
	}

	return "success", nil
}
