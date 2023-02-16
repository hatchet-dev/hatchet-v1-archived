package worker

import (
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/temporal/enums"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/logflusher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulerunner"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/queuechecker"
	"go.temporal.io/sdk/worker"

	hatchetworker "github.com/hatchet-dev/hatchet/internal/config/worker"
)

type WorkerOpts struct {
	RegisterBackground   bool
	RegisterModuleRunner bool

	// TODO: switch this to the worker config
	ServerConfig *server.Config
	WorkerConfig *hatchetworker.Config
}

func NewWorker(opts *WorkerOpts) error {
	// TODO: queue name shouldn't always be background
	hatchetWorker := worker.New(opts.WorkerConfig.TemporalClient.GetClient(), enums.BackgroundQueueName, worker.Options{})

	if opts.RegisterBackground {
		sc := opts.ServerConfig

		lf := logflusher.NewLogFlusher(&logflusher.LogFlusherOpts{
			LogStore:   sc.DefaultLogStore,
			FileStore:  sc.DefaultFileStore,
			Repository: sc.DB.Repository,
		})

		hatchetWorker.RegisterWorkflow(lf.FlushLogs)
		hatchetWorker.RegisterActivity(lf.Flush)

		mqc := modulequeuechecker.NewModuleQueueChecker(sc.ModuleRunQueueManager, sc.DB, *sc.TokenOpts, sc.ServerRuntimeConfig.ServerURL)
		qc := queuechecker.NewQueueChecker(sc.DB.Repository, mqc)

		hatchetWorker.RegisterWorkflow(mqc.ScheduleFromQueue)
		hatchetWorker.RegisterWorkflow(qc.CheckQueues)
	}

	if opts.RegisterModuleRunner {

		mr := modulerunner.NewModuleRunner(opts.WorkerConfig)

		hatchetWorker.RegisterWorkflow(mr.Provision)
		hatchetWorker.RegisterActivity(mr.Run)
	}

	return hatchetWorker.Start()
}
