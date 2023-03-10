//go:generate swagger generate spec

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/staticfileserver"
	"github.com/hatchet-dev/hatchet/api/v1/server/pb"
	"github.com/hatchet-dev/hatchet/api/v1/server/router"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/temporal/dispatcher"
	"github.com/hatchet-dev/hatchet/internal/temporal/server"
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

	configLoader := &loader.ConfigLoader{}

	sc, err := configLoader.LoadServerConfig()

	if err != nil {
		fmt.Fprintf(os.Stdout, "Fatal: could not load server config: %v\n", err)
		os.Exit(1)
	}

	appRouter := router.NewAPIRouter(sc)

	address := fmt.Sprintf(":%d", sc.ServerRuntimeConfig.Port)

	sc.Logger.Info().Msgf("Starting server %v", address)

	if sc.ServerRuntimeConfig.RunTemporalServer {
		startTemporalServerOrDie(configLoader)
	}

	if sc.ServerRuntimeConfig.RunBackgroundWorker {
		startBackgroundWorkerOrDie(configLoader)
	}

	if sc.ServerRuntimeConfig.RunRunnerWorker {
		startRunnerWorkerOrDie(configLoader)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProvisionerServer(grpcServer, pgrpc.NewProvisionerServer(sc))

	http2Server := &http2.Server{}

	runStaticServer := sc.ServerRuntimeConfig.RunStaticFileServer
	var staticServer *chi.Mux

	if runStaticServer {

		staticServer = staticfileserver.NewStaticFileServer(sc.ServerRuntimeConfig.StaticFileServerPath)
	}

	s := &http.Server{
		Addr: address,
		Handler: h2c.NewHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			hasAPIV1Prefix := strings.HasPrefix(request.URL.Path, "/api/v1")

			if request.ProtoMajor != 2 && hasAPIV1Prefix {
				appRouter.ServeHTTP(writer, request)
				return
			}

			if strings.Contains(request.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(writer, request)
				return
			}

			if hasAPIV1Prefix || !runStaticServer {
				appRouter.ServeHTTP(writer, request)
				return
			}

			if runStaticServer {
				staticServer.ServeHTTP(writer, request)
				return
			}
		}), http2Server),
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		sc.Logger.Fatal().Err(err).Msg("Server startup failed")
	}
}

func startBackgroundWorkerOrDie(configLoader *loader.ConfigLoader) {
	bwc, err := configLoader.LoadBackgroundWorkerConfig()

	if err != nil {
		fmt.Fprintf(os.Stdout, "Fatal: could not load background worker config: %v\n", err)
		os.Exit(1)
	}

	err = worker.StartBackgroundWorker(bwc)

	if err != nil {
		fmt.Fprintf(os.Stdout, "Fatal: could not start worker: %v\n", err)
		os.Exit(1)
	}

	err = dispatcher.DispatchBackgroundTasks(bwc.TemporalClient)

	if err != nil {
		fmt.Fprintf(os.Stdout, "Fatal: could not dispatch background workflows: %v\n", err)
		os.Exit(1)
	}
}

func startRunnerWorkerOrDie(configLoader *loader.ConfigLoader) {
	rwc, err := configLoader.LoadRunnerWorkerConfig()

	if err != nil {
		fmt.Fprintf(os.Stdout, "Fatal: could not load runner worker config: %v\n", err)
		os.Exit(1)
	}

	err = worker.StartRunnerWorker(rwc, false)

	if err != nil {
		fmt.Fprintf(os.Stdout, "Fatal: could not start worker: %v\n", err)
		os.Exit(1)
	}
}

func startTemporalServerOrDie(configLoader *loader.ConfigLoader) {
	tc, err := configLoader.LoadTemporalConfig()

	if err != nil {
		fmt.Fprintf(os.Stdout, "Fatal: could not load temporal config: %v\n", err)
		os.Exit(1)
	}

	s, err := server.NewTemporalServer(tc)

	if err != nil {
		fmt.Fprintf(os.Stdout, "Fatal: could not get temporal server: %v\n", err)
		os.Exit(1)
	}

	sui, err := server.NewUIServer(tc.ConfigFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: could not create ui server: %v\n", err)
		os.Exit(1)
	}

	go func() {
		if err := sui.Start(); err != nil {
			panic(err)
		}
	}()

	go func() {
		err = s.Start()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal: unable to start temporal server. Error: %v\n", err)
			panic(err)
		}
	}()
}
