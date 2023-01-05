package users_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hatchet-dev/hatchet/api/serverutils/apitest"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/users"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSuccessful(t *testing.T) {
	declaredUser := testutils.UserModels[0]

	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertStatusCode(t, rr, http.StatusOK)

		// make sure the user is updated by checking the database
		user, err := config.DB.Repository.User().ReadUserByEmail(declaredUser.Email)

		assert.Nil(t, err, "err is nil")
		assert.Equal(t, user.DisplayName, "My New Name")

		return nil
	}, &apitest.APITesterOpts{
		Method:      "POST",
		Route:       "/api/v1/users/current",
		HandlerInit: users.NewUserUpdateCurrentHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
		RequestObj: types.UpdateUserRequest{
			DisplayName: "My New Name",
		},
	}, testutils.InitUsers)
}
