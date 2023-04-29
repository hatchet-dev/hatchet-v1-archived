package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/temporal"
	"github.com/hatchet-dev/hatchet/internal/temporal/enums"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/logflusher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulequeuechecker"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/modulerunner"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/monitordispatcher"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/notifier"
	"github.com/hatchet-dev/hatchet/internal/temporal/workflows/queuechecker"
	"go.temporal.io/sdk/worker"

	"github.com/hatchet-dev/hatchet/internal/config/database"
	hatchetworker "github.com/hatchet-dev/hatchet/internal/config/worker"
	workerconfig "github.com/hatchet-dev/hatchet/internal/config/worker"
)

func StartBackgroundWorker(config *hatchetworker.BackgroundConfig, interruptCh <-chan interface{}) error {
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

	mqc := modulequeuechecker.NewModuleQueueChecker(config.ModuleRunQueueManager, config.DB, *config.TokenOpts, config.ServerURL, config.BroadcastGRPCAddress)
	qc := queuechecker.NewQueueChecker(config.DB.Repository, mqc)
	md := monitordispatcher.NewMonitorDispatcher(config.DefaultLogStore, config.DB, *config.TokenOpts, config.ServerURL, config.BroadcastGRPCAddress)

	backgroundWorker.RegisterWorkflow(mqc.ScheduleFromQueue)
	backgroundWorker.RegisterWorkflow(qc.CheckQueues)
	backgroundWorker.RegisterWorkflow(md.DispatchMonitors)

	go func() {
		<-interruptCh
		backgroundWorker.Stop()
	}()

	return backgroundWorker.Start()
}

type teams struct {
	teams map[string]string

	mu sync.Mutex
}

func StartRunnerWorkerCentralized(rwc *workerconfig.RunnerConfig, dc *database.Config, interruptChan <-chan interface{}, blocking bool) error {
	ts := &teams{
		teams: map[string]string{},
	}

	// spawn a go process that lists teams in the database periodically
	ticker := time.NewTicker(15 * time.Second)

	runFunc := func() error {
		for {
			select {
			case <-interruptChan:
				return nil
			case <-ticker.C:
				var teams []*models.Team

				if err := dc.GormDB.Find(&teams).Error; err != nil {
					fmt.Printf("Fatal: could not list teams: %v\n", err)
					return err
				}

				for _, team := range teams {
					if _, exists := ts.teams[team.ID]; !exists {
						ts.mu.Lock()
						ts.teams[team.ID] = team.ID
						ts.mu.Unlock()

						fmt.Printf("adding new runner worker for team %s\n", team.ID)

						_rwc := *rwc
						var err error

						newOpts := rwc.TemporalClient.GetOpts()

						newOpts.Namespace = team.ID

						_rwc.TemporalClient, err = temporal.NewTemporalClient(newOpts)

						if err != nil {
							return fmt.Errorf("Fatal: could not create new temporal client: %v\n", err)
						}

						// create a new runner worker
						err = startRunnerWorker(&_rwc, false, interruptChan)

						if err != nil {
							return fmt.Errorf("Fatal: could not start runner worker: %v\n", err)
						}

						fmt.Printf("successfully added new runner worker for team %s\n", team.ID)
					}
				}
			}
		}

	}

	if blocking {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil

}

func StartRunnerWorkerDecentralized(rwc *workerconfig.RunnerConfig, interruptChan <-chan interface{}, blocking bool) error {
	runFunc := func() error {
		err := startRunnerWorker(rwc, true, interruptChan)

		if err != nil {
			fmt.Printf("Fatal: could not start worker: %v\n", err)
			return err
		}

		return nil
	}

	if blocking {
		return runFunc()
	} else {
		go runFunc()
	}

	return nil
}

func startRunnerWorker(config *hatchetworker.RunnerConfig, blocking bool, interruptChan <-chan interface{}) error {
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
		return runnerWorker.Run(interruptChan)
	}

	go func() {
		<-interruptChan
		runnerWorker.Stop()
	}()

	return runnerWorker.Start()
}
