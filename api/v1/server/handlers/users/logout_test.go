package users_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hatchet-dev/hatchet/api/serverutils/apitest"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/users"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

// Note this is a very basic test to ensure that the status is correct. The login/logout testing
// relies on integration + e2e testing.
func TestLogoutUserSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		assert.Equal(t, http.StatusOK, rr.Result().StatusCode, "status should be ok")

		return nil
	}, &apitest.APITesterOpts{
		Method:      "POST",
		Route:       "/api/v1/users/logout",
		HandlerInit: users.NewUserLogoutHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
	}, testutils.InitUsers)
}
