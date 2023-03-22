package main

import (
	"fmt"
	"os"

	"github.com/hatchet-dev/hatchet/cmd/cmdutils"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
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
		interruptChan := cmdutils.InterruptChan()
		rwc, err := configLoader.LoadRunnerWorkerConfig()

		if err != nil {
			fmt.Printf("Fatal: could not load runner worker config: %v\n", err)
			os.Exit(1)
		}

		err = worker.StartRunnerWorker(rwc, true, interruptChan)

		if err != nil {
			fmt.Printf("Fatal: could not start worker: %v\n", err)
			os.Exit(1)
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
