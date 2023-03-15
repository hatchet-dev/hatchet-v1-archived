package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hatchet-dev/hatchet/cmd/cmdutils"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/temporal/worker"
)

// Version will be linked by an ldflag during build
var Version string = "v0.1.0-alpha.0"

func main() {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "print version and exit")
	flag.Parse()

	// Exit safely when version is used
	if versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	configLoader := &loader.ConfigLoader{}
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
}
