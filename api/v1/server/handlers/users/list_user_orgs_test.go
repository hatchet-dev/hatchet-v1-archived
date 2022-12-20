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
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestListUserOrgsSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		gotResponse := &types.ListUserOrgsResponse{}

		err := json.NewDecoder(rr.Body).Decode(gotResponse)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, &types.PaginationResponse{
			NumPages:    1,
			CurrentPage: 0,
			NextPage:    0,
		}, gotResponse.Pagination, "pagination should be equal")

		assert.Equal(t, 1, len(gotResponse.Rows), "length of rows should be 1")

		assert.Equal(t, testutils.OrgModels[0].DisplayName, gotResponse.Rows[0].DisplayName, "org display names should be equal")

		return nil
	}, &apitest.APITesterOpts{
		Method:      "GET",
		Route:       "/api/v1/users/current/organizations",
		HandlerInit: users.NewUserOrgListHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
	}, testutils.InitUsers, testutils.InitOrgs)
}
