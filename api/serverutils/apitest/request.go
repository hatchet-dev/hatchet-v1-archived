//go:build test

package apitest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
)

func WithAuthUser(initDataIndex uint) GenerateRequestCtx {
	return func(t *testing.T, config *server.Config) (interface{}, interface{}) {
		return types.UserScope, testutils.InitDataAll.Users[initDataIndex]
	}
}

func WithURLParamPAT(initDataIndex uint) GenerateURLParam {
	return func(t *testing.T, config *server.Config, currParams map[string]string) map[string]string {
		currParams[string(types.PersonalAccessTokenURLParam)] = testutils.InitDataAll.PATs[initDataIndex].ID

		return currParams
	}
}

type GenerateRequestCtx func(t *testing.T, config *server.Config) (interface{}, interface{})

// GenerateURLParam takes in the server config and outputs URL params. It's meant to populate
// URL params that are only available from the init data.
type GenerateURLParam func(t *testing.T, config *server.Config, currParams map[string]string) map[string]string

func GetRequestAndRecorder(t *testing.T, method, route string, requestObj interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var reader io.Reader = nil

	if requestObj != nil {
		data, err := json.Marshal(requestObj)

		if err != nil {
			t.Fatal(err)
		}

		reader = strings.NewReader(string(data))
	}

	// method and route don't actually matter since this is meant to test handlers
	req, err := http.NewRequest(method, route, reader)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	return req, rr
}

func WithURLParams(t *testing.T, req *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	routeParams := &chi.RouteParams{}

	for key, val := range params {
		routeParams.Add(key, val)
	}

	rctx.URLParams = *routeParams

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	return req
}

type failingDecoderValidator struct {
	config *server.Config
}

func (f *failingDecoderValidator) DecodeAndValidate(
	w http.ResponseWriter,
	r *http.Request,
	v interface{},
) (ok bool) {
	apierrors.HandleAPIError(f.config.Logger, f.config.ErrorAlerter, w, r, apierrors.NewErrInternal(fmt.Errorf("fake error")), true)
	return false
}

func (f *failingDecoderValidator) DecodeAndValidateNoWrite(
	r *http.Request,
	v interface{},
) error {
	return fmt.Errorf("fake error")
}

func NewFailingDecoderValidator(config *server.Config) handlerutils.RequestDecoderValidator {
	return &failingDecoderValidator{config}
}
