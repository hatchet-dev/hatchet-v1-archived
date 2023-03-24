package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hatchet-dev/hatchet/cmd/cmdutils"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	workerconfig "github.com/hatchet-dev/hatchet/internal/config/worker"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/temporal"
	"github.com/hatchet-dev/hatchet/internal/temporal/worker"
	"github.com/spf13/cobra"
)

// Version will be linked by an ldflag during build
var Version string = "v0.1.0-alpha.0"

var printVersion bool
var configDirectory string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hatchet-server",
	Short: "hatchet-server runs a Hatchet instance.",
	Run: func(cmd *cobra.Command, args []string) {
		if printVersion {
			fmt.Println(Version)
			os.Exit(0)
		}

		configLoader := loader.NewConfigLoader(Version, configDirectory)

		rwc, err := configLoader.LoadRunnerWorkerConfig()

		if err != nil {
			fmt.Printf("Fatal: could not load runner worker config: %v\n", err)
			os.Exit(1)
		}

		// if this is a centralized mechanism, load a database connection
		if rwc.ProvisionerMechanism == "centralized" {
			dc, err := configLoader.LoadDatabaseConfig()

			if err != nil {
				fmt.Printf("Fatal: could not load database config: %v\n", err)
				os.Exit(1)
			}

			startWorkerCentralized(rwc, dc)
		} else {
			startWorkerDecentralized(rwc)
		}
	},
}

func main() {
	rootCmd.PersistentFlags().BoolVar(
		&printVersion,
		"version",
		false,
		"print version and exit.",
	)

	rootCmd.PersistentFlags().StringVar(
		&configDirectory,
		"config",
		"",
		"The path the config folder.",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type teams struct {
	teams map[string]string

	mu sync.Mutex
}

func startWorkerCentralized(rwc *workerconfig.RunnerConfig, dc *database.Config) {
	interruptChan := cmdutils.InterruptChan()

	ts := &teams{
		teams: map[string]string{},
	}

	// spawn a go process that lists teams in the database periodically
	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-interruptChan:
			return
		case <-ticker.C:
			var teams []*models.Team

			if err := dc.GormDB.Find(&teams).Error; err != nil {
				fmt.Printf("Fatal: could not list teams: %v\n", err)
				os.Exit(1)
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
						fmt.Printf("Fatal: could not create new temporal client: %v\n", err)
						os.Exit(1)
					}

					// create a new runner worker
					err = worker.StartRunnerWorker(&_rwc, false, interruptChan)

					if err != nil {
						fmt.Printf("Fatal: could not start runner worker: %v\n", err)
						os.Exit(1)
					}

					fmt.Printf("successfully added new runner worker for team %s\n", team.ID)
				}
			}
		}
	}
}

func startWorkerDecentralized(rwc *workerconfig.RunnerConfig) {
	interruptChan := cmdutils.InterruptChan()

	err := worker.StartRunnerWorker(rwc, true, interruptChan)

	if err != nil {
		fmt.Printf("Fatal: could not start worker: %v\n", err)
		os.Exit(1)
	}
}
