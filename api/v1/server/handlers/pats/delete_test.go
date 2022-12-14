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

func TestDeletePATSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		gotPATResp := &types.DeletePATResponse{}

		err := json.NewDecoder(rr.Body).Decode(gotPATResp)

		if err != nil {
			t.Fatal(err)
		}

		AssertPATsEqual(t, testutils.PATModels[0].ToAPIType(), (*types.PersonalAccessToken)(gotPATResp))

		return nil
	}, &apitest.APITesterOpts{
		Method:      "DELETE",
		Route:       fmt.Sprintf("/api/v1/users/current/settings/pats/{%s}", string(types.PersonalAccessTokenURLParam)),
		HandlerInit: pats.NewPATDeleteHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamPAT(0),
		},
	}, testutils.InitUsers, testutils.InitPATs)
}

func TestDeletePATInvalidTokenID(t *testing.T) {
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
		Method:      "DELETE",
		Route:       fmt.Sprintf("/api/v1/users/current/settings/pats/{%s}", string(types.PersonalAccessTokenURLParam)),
		HandlerInit: pats.NewPATDeleteHandler,
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
