package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hatchet-dev/hatchet/cmd/hatchet-runner/cli"
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

	cli.Execute()
}
