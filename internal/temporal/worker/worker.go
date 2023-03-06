package worker

import (
	"github.com/hatchet-dev/hatchet/internal/temporal/enums"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/logflusher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulerunner"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/monitordispatcher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/notifier"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/queuechecker"
	"go.temporal.io/sdk/worker"

	hatchetworker "github.com/hatchet-dev/hatchet/internal/config/worker"
)

func StartBackgroundWorker(config *hatchetworker.BackgroundConfig) error {
	tc, err := config.TemporalClient.GetClient(enums.BackgroundQueueName)

	if err != nil {
		return err
	}

	backgroundWorker := worker.New(tc, enums.BackgroundQueueName, worker.Options{})

	lf := logflusher.NewLogFlusher(&logflusher.LogFlusherOpts{
		LogStore:   config.DefaultLogStore,
		FileStore:  config.DefaultFileStore,
		Repository: config.DB.Repository,
	})

	backgroundWorker.RegisterWorkflow(lf.FlushLogs)
	backgroundWorker.RegisterActivity(lf.Flush)

	notifier := notifier.NewNotifier(config)

	backgroundWorker.RegisterWorkflow(notifier.NotifyWorkflow)

	mqc := modulequeuechecker.NewModuleQueueChecker(config.ModuleRunQueueManager, config.DB, *config.TokenOpts, config.ServerURL)
	qc := queuechecker.NewQueueChecker(config.DB.Repository, mqc)
	md := monitordispatcher.NewMonitorDispatcher(config.DefaultLogStore, config.DB, *config.TokenOpts, config.ServerURL)

	backgroundWorker.RegisterWorkflow(mqc.ScheduleFromQueue)
	backgroundWorker.RegisterWorkflow(qc.CheckQueues)
	backgroundWorker.RegisterWorkflow(md.DispatchMonitors)

	return backgroundWorker.Start()
}

func StartRunnerWorker(config *hatchetworker.RunnerConfig, blocking bool) error {
	tc, err := config.TemporalClient.GetClient(enums.ModuleRunQueueName)

	if err != nil {
		return err
	}

	runnerWorker := worker.New(tc, enums.ModuleRunQueueName, worker.Options{})

	mr := modulerunner.NewModuleRunner(config)

	runnerWorker.RegisterWorkflow(mr.Provision)
	runnerWorker.RegisterActivity(mr.Run)

	// TODO: name of workflow vs activity is confusing
	runnerWorker.RegisterWorkflow(mr.RunMonitor)
	runnerWorker.RegisterActivity(mr.Monitor)

	if blocking {
		return runnerWorker.Run(worker.InterruptCh())
	}

	return runnerWorker.Start()
}
