package modulerunner

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

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

func (mqc *ModuleRunner) RunMonitor(ctx workflow.Context, input MonitorInput) (string, error) {
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

	// TODO: output should be a monitor result, not a run output. This will be cased on and piped into an alerting system
	var runOutput string

	runErr := workflow.ExecuteActivity(ctx, "Monitor", input).Get(ctx, &runOutput)

	if runErr != nil {
		return "", runErr
	}

	return runOutput, nil
}
