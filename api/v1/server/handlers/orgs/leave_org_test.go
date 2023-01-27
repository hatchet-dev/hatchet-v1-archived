package orgs_test

import (
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

func TestLeaveOrg(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		assert.Equal(t, rr.Result().StatusCode, http.StatusOK)

		// ensure only the original owner is listed
		orgMembers, _, err := config.DB.Repository.Org().ListOrgMembersByOrgID(testutils.OrgModels[0].ID, false)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, 1, len(orgMembers), "only a single member left")
		assert.Equal(t, testutils.UserModels[0].Email, orgMembers[0].User.Email, "single member is the original owner")

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",

		Route:       fmt.Sprintf("/api/v1/organizations/{%s}/leave", types.URLParamOrgID),
		HandlerInit: orgs.NewOrgLeaveHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(2),
			apitest.WithOrg(0),
			apitest.WithAuthOrgMember(0, testutils.UserModels[2].Email),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamOrg(0),
		},
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}
