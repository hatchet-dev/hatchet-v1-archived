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
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
)

func TestRevokePATSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		gotPATResp := &types.RevokePATResponse{}

		err := json.NewDecoder(rr.Body).Decode(gotPATResp)

		if err != nil {
			t.Fatal(err)
		}

		expResp := testutils.PATModels[0].ToAPIType()

		expResp.Revoked = true

		AssertPATsEqual(t, expResp, (*types.PersonalAccessToken)(gotPATResp))

		return nil
	}, &apitest.APITesterOpts{
		Method:      "POST",
		Route:       fmt.Sprintf("/api/v1/users/current/settings/pats/{%s}/revoke", string(types.PersonalAccessTokenURLParam)),
		HandlerInit: pats.NewPATRevokeHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamPAT(0),
		},
	}, testutils.InitUsers, testutils.InitPATs)
}

func TestRevokePATInvalidTokenID(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertResponseError(t, rr, http.StatusNotFound, &types.APIErrors{
			Errors: []types.APIError{
				{
					Code:        types.ErrCodeNotFound,
					Description: types.GenericResourceNotFound,
				},
			},
		})

		return nil
	}, &apitest.APITesterOpts{
		Method:      "POST",
		Route:       fmt.Sprintf("/api/v1/users/current/settings/pats/{%s}/revoke", string(types.PersonalAccessTokenURLParam)),
		HandlerInit: pats.NewPATGetHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
		URLGenerators: []apitest.GenerateURLParam{
			func(t *testing.T, config *server.Config, currParams map[string]string) map[string]string {
				currParams[fmt.Sprintf("%s", string(types.PersonalAccessTokenURLParam))] = "not-a-token-id"
				return currParams
			},
		},
	}, testutils.InitUsers, testutils.InitPATs)
}
