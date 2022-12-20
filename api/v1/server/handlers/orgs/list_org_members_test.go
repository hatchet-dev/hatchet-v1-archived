package orgs_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hatchet-dev/hatchet/api/serverutils/apitest"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/orgs"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestListOrgMembersSingleOwner(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		listOrgMembersResponse := &types.ListOrgMembersResponse{}

		err := json.NewDecoder(rr.Body).Decode(listOrgMembersResponse)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, &types.PaginationResponse{
			NumPages:    1,
			CurrentPage: 0,
			NextPage:    0,
		}, listOrgMembersResponse.Pagination, "pagination should be equal")

		assert.Equal(t, 1, len(listOrgMembersResponse.Rows), "length of rows should be 1")

		// ensure that the first user is the owner
		assert.Equal(t, testutils.UserModels[0].Email, listOrgMembersResponse.Rows[0].User.Email, "only org member should be owner")
		assert.Equal(t, "owner", listOrgMembersResponse.Rows[0].OrgPolicies[0].Name, "policy should be owner")

		return nil
	}, &apitest.APITesterOpts{
		Method:      "GET",
		Route:       fmt.Sprintf("/api/v1/organizations/{%s}/members", types.URLParamOrgID),
		HandlerInit: orgs.NewOrgListMembersHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
			apitest.WithOrg(0),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamOrg(0),
		},
	}, testutils.InitUsers, testutils.InitOrgs)
}
