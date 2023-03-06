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
	configLoader := &loader.EnvConfigLoader{}
	sc, err = configLoader.LoadServerConfigFromEnv()

	if err != nil {
		fmt.Printf("Fatal: could not load server config: %v", err)
		os.Exit(1)
	}
}
