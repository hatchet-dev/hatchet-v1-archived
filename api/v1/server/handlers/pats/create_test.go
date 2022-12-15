package pats_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hatchet-dev/hatchet/api/serverutils/apitest"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/pats"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/auth/token"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCreatePATSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		gotPATResp := &types.CreatePATResponse{}

		err := json.NewDecoder(rr.Body).Decode(gotPATResp)

		if err != nil {
			t.Fatal(err)
		}

		// assert that the uuid on the pat field is valid
		assert.True(t, uuidutils.IsValidUUID(gotPATResp.PersonalAccessToken.ID))

		// check that the generated token is valid
		isValid, err := token.IsPATValid(gotPATResp.Token, config.DB.Repository.PersonalAccessToken(), config.TokenOpts)

		assert.Nil(t, err, "error should be nil")
		assert.True(t, isValid, "token should be valid")

		AssertPATsEqual(t, &types.PersonalAccessToken{
			DisplayName: "my-pat-test",
		}, &gotPATResp.PersonalAccessToken)

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  "/api/v1/users/current/settings/pats",
		RequestObj: types.CreatePATRequest{
			DisplayName: "my-pat-test",
		},
		HandlerInit: pats.NewPATCreateHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
	}, testutils.InitUsers)
}

func TestCreatePATDuplicate(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertResponseError(t, rr, http.StatusBadRequest, &types.APIErrors{
			Errors: []types.APIError{
				{
					Code:        types.ErrCodeBadRequest,
					Description: fmt.Sprintf("Personal access token already exists with display_name %s for this user", testutils.PATModels[0].DisplayName),
				},
			},
		})

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  "/api/v1/users/current/settings/pats",
		RequestObj: types.CreatePATRequest{
			DisplayName: testutils.PATModels[0].DisplayName,
		},
		HandlerInit: pats.NewPATCreateHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
	}, testutils.InitUsers, testutils.InitPATs)
}

func TestCreatePATMissingDisplayName(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertResponseError(t, rr, http.StatusBadRequest, &types.APIErrors{
			Errors: []types.APIError{
				{
					Code:        types.ErrCodeBadRequest,
					Description: "validation failed on field 'DisplayName' on condition 'required'",
				},
			},
		})

		return nil
	}, &apitest.APITesterOpts{
		Method:      "POST",
		Route:       "/api/v1/users/current/settings/pats",
		RequestObj:  types.CreatePATRequest{},
		HandlerInit: pats.NewPATCreateHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
	}, testutils.InitUsers)
}
