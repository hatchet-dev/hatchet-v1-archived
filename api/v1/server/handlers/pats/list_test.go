package pats_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hatchet-dev/hatchet/api/serverutils/apitest"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/pats"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestListPATsSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		listPATsResp := &types.ListPATsResponse{}

		err := json.NewDecoder(rr.Body).Decode(listPATsResp)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, &types.PaginationResponse{
			NumPages:    1,
			CurrentPage: 0,
			NextPage:    0,
		}, listPATsResp.Pagination, "pagination should be equal")

		assert.Equal(t, 1, len(listPATsResp.Rows), "length of rows should be 1")

		AssertPATsEqual(t, testutils.PATModels[0].ToAPIType(), (*types.PersonalAccessToken)(listPATsResp.Rows[0]))

		return nil
	}, &apitest.APITesterOpts{
		Method:      "GET",
		Route:       "/api/v1/users/current/settings/pats",
		HandlerInit: pats.NewPATListHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamPAT(0),
		},
	}, testutils.InitUsers, testutils.InitPATs)
}
