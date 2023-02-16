package worker

import (
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/temporal"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/logflusher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulerunner"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/queuechecker"
	"go.temporal.io/sdk/worker"
)

func NewBackgroundWorker(config *server.Config) error {
	backgroundWorker := worker.New(config.TemporalClient.GetClient(), temporal.BackgroundQueueName, worker.Options{})

	lf := logflusher.NewLogFlusher(&logflusher.LogFlusherOpts{
		LogStore:   config.DefaultLogStore,
		FileStore:  config.DefaultFileStore,
		Repository: config.DB.Repository,
	})

	backgroundWorker.RegisterWorkflow(lf.FlushLogs)
	backgroundWorker.RegisterActivity(lf.Flush)

	mqc := modulequeuechecker.NewModuleQueueChecker(config.ModuleRunQueueManager, config.DB.Repository)
	qc := queuechecker.NewQueueChecker(config.DB.Repository, mqc)

	backgroundWorker.RegisterWorkflow(mqc.ScheduleFromQueue)
	backgroundWorker.RegisterWorkflow(qc.CheckQueues)

	return backgroundWorker.Start()
}

func NewModuleRunWorker(config *server.Config) error {
	moduleRunWorker := worker.New(config.TemporalClient.GetClient(), temporal.ModuleRunSchedulerQueueName, worker.Options{
		DisableWorkflowWorker: true,
	})

	moduleRunWorker.RegisterActivity(modulerunner.Run)

	return moduleRunWorker.Start()
}
