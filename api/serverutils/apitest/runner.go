// go:build test
package apitest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/internal/config/database"
	"github.com/hatchet-dev/hatchet/internal/config/loader"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/config/shared"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
)

type APITester struct {
	conf *server.Config
}

type HandlerInitFunc func(
	config *server.Config,
	decoderValidator handlerutils.RequestDecoderValidator,
	writer handlerutils.ResultWriter,
) http.Handler

type APITesterOpts struct {
	Method, Route string
	RequestObj    interface{}
	HandlerInit   HandlerInitFunc
	CtxGenerators []GenerateRequestCtx
}

type APITestFunc func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error

func RunAPITest(t *testing.T, test APITestFunc, opts *APITesterOpts, initMethods ...testutils.InitDataFunc) {
	t.Helper()

	// initialize the server config
	apiTester := new(APITester)
	loader := new(loader.EnvConfigLoader)

	// load the shared configuration (server config has to be loaded after database config)
	sharedConfig, err := loader.LoadSharedConfigFromConfigFile(&shared.ConfigFile{
		Debug: true,
	})

	if err != nil {
		t.Fatalf("%v\n", err)
	}

	// wrap the database test
	testutils.RunTestWithDatabase(t, func(conf *database.Config) error {
		apiTester.conf, err = loader.LoadServerConfigFromConfigFile(&server.ConfigFile{
			Port:             8080,
			BasicAuthEnabled: true,
			CookieName:       "hatchet",
			CookieSecrets: []string{
				"random_hash_key_",
				"random_block_key",
			},
		}, conf, sharedConfig)

		if err != nil {
			t.Fatalf("%v\n", err)
		}

		// create a new response recorder
		req, rr := GetRequestAndRecorder(
			t,
			opts.Method,
			opts.Route,
			opts.RequestObj,
		)

		if opts.CtxGenerators != nil {
			for _, gen := range opts.CtxGenerators {
				ctx := req.Context()
				key, val := gen(apiTester.conf)
				ctx = context.WithValue(ctx, key, val)
				req = req.WithContext(ctx)
			}
		}

		handler := opts.HandlerInit(
			apiTester.conf,
			handlerutils.NewDefaultRequestDecoderValidator(apiTester.conf.Logger, apiTester.conf.ErrorAlerter),
			handlerutils.NewDefaultResultWriter(apiTester.conf.Logger, apiTester.conf.ErrorAlerter),
		)

		handler.ServeHTTP(rr, req)

		err = test(apiTester.conf, rr, req)

		if err != nil {
			t.Fatalf("%v\n", err)
		}

		return nil
	}, initMethods...)
}
