package users_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hatchet-dev/hatchet/api/serverutils/apitest"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/users"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSuccessful(t *testing.T) {
	declaredUser := testutils.UserModels[0]

	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertStatusCode(t, rr, http.StatusAccepted)

		// make sure the user is deleted by checking the database
		_, failingErr := config.DB.Repository.User().ReadUserByEmail(declaredUser.Email)

		assert.NotNil(t, failingErr, "err is not nil")
		assert.ErrorIs(t, repository.RepositoryErrorNotFound, failingErr)

		return nil
	}, &apitest.APITesterOpts{
		Method:      "DELETE",
		Route:       "/api/v1/users/current",
		HandlerInit: users.NewUserDeleteCurrentHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
	}, testutils.InitUsers)
}
