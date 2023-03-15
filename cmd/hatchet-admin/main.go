package main

import (
	"github.com/hatchet-dev/hatchet/cmd/hatchet-admin/cli"
)

func main() {
	// var versionFlag bool
	// flag.BoolVar(&versionFlag, "version", false, "print version and exit")
	// flag.Parse()

	// // Exit safely when version is used
	// if versionFlag {
	// 	fmt.Println(Version)
	// 	os.Exit(0)
	// }

	cli.Execute()
}
