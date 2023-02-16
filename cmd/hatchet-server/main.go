//go:generate swagger generate spec

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/hatchet-dev/hatchet/api/v1/server/pb"
	"github.com/hatchet-dev/hatchet/api/v1/server/router"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"
	"github.com/hatchet-dev/hatchet/internal/temporal/worker"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	pgrpc "github.com/hatchet-dev/hatchet/api/v1/server/grpc"
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

	if sc.ServerRuntimeConfig.RunWorkers {
		wc, err := configLoader.LoadWorkerConfigFromEnv()

		if err != nil {
			fmt.Printf("Fatal: could not load worker config: %v", err)
			os.Exit(1)
		}

		err = worker.NewWorker(&worker.WorkerOpts{
			ServerConfig:         sc,
			WorkerConfig:         wc,
			RegisterBackground:   true,
			RegisterModuleRunner: true,
		})

		if err != nil {
			fmt.Printf("Fatal: could not start worker: %v", err)
			os.Exit(1)
		}
	}

	err = dispatcher.DispatchBackgroundTasks(sc.TemporalClient.GetClient())

	if err != nil {
		fmt.Printf("Fatal: could not dispatch background workflows: %v", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProvisionerServer(grpcServer, pgrpc.NewProvisionerServer(sc))

	http2Server := &http2.Server{}
	s := &http.Server{
		Addr: address,
		Handler: h2c.NewHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.ProtoMajor != 2 {
				appRouter.ServeHTTP(writer, request)
				return
			}

			if strings.Contains(request.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(writer, request)
				return
			}

			appRouter.ServeHTTP(writer, request)
		}), http2Server),
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		sc.Logger.Fatal().Err(err).Msg("Server startup failed")
	}
}
