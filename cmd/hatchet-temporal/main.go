package main

import (
	"fmt"
	goLog "log"
	"os"

	"github.com/hatchet-dev/hatchet/cmd/cmdutils"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/temporal/server"

	// Load sqlite storage driver
	_ "go.temporal.io/server/common/persistence/sql/sqlplugin/sqlite"
)

type uiConfig struct {
	Host                string
	Port                int
	TemporalGRPCAddress string
	EnableUI            bool
	CodecEndpoint       string
}

func main() {
	configLoader := &loader.ConfigLoader{}
	interruptChan := cmdutils.InterruptChan()
	tc, err := configLoader.LoadTemporalConfig()

	if err != nil {
		fmt.Printf("Fatal: could not load server config: %v\n", err)
		os.Exit(1)
	}

	s, err := server.NewTemporalServer(tc, interruptChan)

	if err != nil {
		goLog.Fatal(err)
	}

	sui, err := server.NewUIServer(tc.ConfigFile)

	if err != nil {
		goLog.Fatal(fmt.Sprintf("Unable to create UI server. Error: %v\n", err))
	}

	go func() {
		if err := sui.Start(); err != nil {
			panic(err)
		}
	}()

	err = s.Start()

	if err != nil {
		goLog.Fatal(fmt.Sprintf("Unable to start server. Error: %v\n", err))
	}

	return
}
