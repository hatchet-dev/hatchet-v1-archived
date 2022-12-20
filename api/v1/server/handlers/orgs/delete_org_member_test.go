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

func TestDeleteOrgMemberSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		assert.Equal(t, http.StatusAccepted, rr.Result().StatusCode, "status should be accepted")

		// query the database directly to ensure that user was removed
		orgMembers, _, err := config.DB.Repository.Org().ListOrgMembersByOrgID(testutils.OrgModels[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, 1, len(orgMembers))

		return nil
	}, &apitest.APITesterOpts{
		Method:      "DELETE",
		Route:       fmt.Sprintf("/api/v1/organizations/{%s}/members/{%s}", types.URLParamOrgID, types.URLParamOrgMemberID),
		HandlerInit: orgs.NewOrgDeleteMemberHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
			apitest.WithOrg(0),
			apitest.WithOrgMember(0),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamOrg(0),
		},
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}
