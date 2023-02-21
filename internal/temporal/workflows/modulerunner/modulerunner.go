package modulerunner

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type RunModuleInput struct {
	TeamID, ModuleID string
}

func (mqc *ModuleRunner) Provision(ctx workflow.Context, input RunInput) (string, error) {
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

	ctx = workflow.WithActivityOptions(ctx, options)

	var runOutput string

	runErr := workflow.ExecuteActivity(ctx, "Run", input).Get(ctx, &runOutput)

	if runErr != nil {
		return "", runErr
	}

	return runOutput, nil
}
