package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hatchet-dev/hatchet/api/serverutils/apitest"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/users"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestLoginUserSuccessful(t *testing.T) {
	declaredUser := testutils.UserModels[0]

	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		gotUser := &types.User{}

		err := json.NewDecoder(rr.Body).Decode(gotUser)

		if err != nil {
			t.Fatal(err)
		}

		// assert that the uuid is valid
		assert.True(t, uuidutils.IsValidUUID(gotUser.ID))

		expUser := &types.User{
			APIResourceMeta: gotUser.APIResourceMeta,
			DisplayName:     declaredUser.DisplayName,
			Email:           declaredUser.Email,
			EmailVerified:   declaredUser.EmailVerified,
		}

		// assert the rest of the user object is valid
		assert.Equal(t, expUser, gotUser, "user should be equal")

		return nil
	}, &apitest.APITesterOpts{
		Method:      "POST",
		Route:       "/api/v1/users/login",
		HandlerInit: users.NewUserLoginHandler,
		RequestObj: &types.LoginUserRequest{
			Email:    "user1@gmail.com",
			Password: "Abcdefgh123",
		},
	}, testutils.InitUsers)
}

func TestLoginUserBadEmail(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertResponseError(t, rr, http.StatusUnauthorized, &types.APIErrors{
			Errors: []types.APIError{types.InvalidEmailOrPassword},
		})

		return nil
	}, &apitest.APITesterOpts{
		Method:      "POST",
		Route:       "/api/v1/users/login",
		HandlerInit: users.NewUserLoginHandler,
		RequestObj: &types.LoginUserRequest{
			Email:    "not-registered@gmail.com",
			Password: "Abcdefgh123",
		},
	}, testutils.InitUsers)
}

func TestLoginUserBadPassword(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertResponseError(t, rr, http.StatusUnauthorized, &types.APIErrors{
			Errors: []types.APIError{types.InvalidEmailOrPassword},
		})

		return nil
	}, &apitest.APITesterOpts{
		Method:      "POST",
		Route:       "/api/v1/users/login",
		HandlerInit: users.NewUserLoginHandler,
		RequestObj: &types.LoginUserRequest{
			Email:    "user1@gmail.com",
			Password: "Abcdefgh123!!",
		},
	}, testutils.InitUsers)
}
