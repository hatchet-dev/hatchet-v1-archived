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
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUpdateOrgMemberPoliciesSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		updateOrgMemberPoliciesResp := &types.UpdateOrgMemberPoliciesResponse{}

		err := json.NewDecoder(rr.Body).Decode(updateOrgMemberPoliciesResp)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "member", updateOrgMemberPoliciesResp.OrgPolicies[0].Name)

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  fmt.Sprintf("/api/v1/organizations/{%s}/members/{%s}/update_policies", string(types.URLParamOrgID), string(types.URLParamOrgMemberID)),
		RequestObj: types.UpdateOrgMemberPoliciesRequest{
			Policies: []types.OrganizationPolicyReference{
				{
					Name: string(models.PresetPolicyNameMember),
				},
			},
		},
		HandlerInit: orgs.NewOrgUpdateMemberPoliciesHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
			apitest.WithOrg(0),
			apitest.WithAuthOrgMember(0, testutils.UserModels[0].Email),
			apitest.WithOrgMember(0),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamOrg(0),
		},
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}
