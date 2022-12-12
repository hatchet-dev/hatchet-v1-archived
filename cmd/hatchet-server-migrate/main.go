//go:generate swagger generate spec

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/migrate"
)

// Version will be linked by an ldflag during build
var Version string = "dev-ce"

func main() {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "print version and exit")
	flag.Parse()

	// Exit safely when version is used
	if versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	configLoader := &loader.EnvConfigLoader{}
	dc, err := configLoader.LoadDatabaseConfig()

	if err != nil {
		fmt.Printf("Fatal: could not load database config: %v", err)
		os.Exit(1)
	}

	err = migrate.AutoMigrate(dc.GormDB, true)

	if err != nil {
		fmt.Printf("Fatal: could not run auto migration: %v", err)
		os.Exit(1)
	}
}
