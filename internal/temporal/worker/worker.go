package worker

import (
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/temporal"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/logflusher"
	"go.temporal.io/sdk/worker"
)

func NewWorker(config *server.Config) error {
	backgroundWorker := worker.New(config.TemporalClient.GetClient(), temporal.BackgroundQueueName, worker.Options{})

	lf := logflusher.NewLogFlusher(&logflusher.LogFlusherOpts{
		LogStore:   config.DefaultLogStore,
		FileStore:  config.DefaultFileStore,
		Repository: config.DB.Repository,
	})

	backgroundWorker.RegisterWorkflow(lf.FlushLogs)
	backgroundWorker.RegisterActivity(lf.Flush)

	return backgroundWorker.Start()
}
