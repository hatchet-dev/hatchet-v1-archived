package cli

import (
	"fmt"
	"os"

	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

var sc *server.Config

func init() {
	var err error
	configLoader := &loader.ConfigLoader{}
	sc, err = configLoader.LoadServerConfig()

	if err != nil {
		fmt.Printf("Fatal: could not load server config: %v\n", err)
		os.Exit(1)
	}
}
