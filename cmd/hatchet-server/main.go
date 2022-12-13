//go:generate swagger generate spec

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/hatchet-dev/hatchet/api/v1/server/router"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
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
	sc, err := configLoader.LoadServerConfigFromEnv()

	if err != nil {
		fmt.Printf("Fatal: could not load server config: %v", err)
		os.Exit(1)
	}

	appRouter := router.NewAPIRouter(sc)

	address := fmt.Sprintf(":%d", sc.ServerRuntimeConfig.Port)

	sc.Logger.Info().Msgf("Starting server %v", address)

	s := &http.Server{
		Addr:    address,
		Handler: appRouter,
		// ReadTimeout:  config.ServerConf.TimeoutRead,
		// WriteTimeout: config.ServerConf.TimeoutWrite,
		// IdleTimeout:  config.ServerConf.TimeoutIdle,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		sc.Logger.Fatal().Err(err).Msg("Server startup failed")
	}
}
