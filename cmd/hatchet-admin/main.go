package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hatchet-dev/hatchet/cmd/hatchet-admin/cli"
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

	cli.Execute()
}
