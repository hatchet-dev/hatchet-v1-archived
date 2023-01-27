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
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUpdateOrgOwnerSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		assert.Equal(t, http.StatusOK, rr.Result().StatusCode, "status should be ok")

		// query the database directly to ensure that user was added
		orgMembers, _, err := config.DB.Repository.Org().ListOrgMembersByOrgID(testutils.OrgModels[0].ID, false)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, 2, len(orgMembers))

		var newOwnerMember *models.OrganizationMember
		var adminMember *models.OrganizationMember

		for _, orgMember := range orgMembers {
			pName := orgMember.OrgPolicies[0].PolicyName

			if pName == string(models.PresetPolicyNameOwner) {
				newOwnerMember = orgMember
			} else if pName == string(models.PresetPolicyNameAdmin) {
				adminMember = orgMember
			}
		}

		assert.NotNil(t, newOwnerMember)
		assert.Equal(t, testutils.UserModels[2].ID, newOwnerMember.UserID, "new owner should be second user")

		assert.NotNil(t, adminMember)
		assert.Equal(t, testutils.UserModels[0].ID, adminMember.UserID, "new admin should be previous owner")

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  fmt.Sprintf("/api/v1/organizations/{%s}/change_owner", string(types.URLParamOrgID)),
		RequestObj: types.UpdateOrgOwnerRequest{
			NewOwnerMemberID: "body_param_new_owner_member_id",
		},
		HandlerInit: orgs.NewOrgUpdateOwnerHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
			apitest.WithOrg(0),
			apitest.WithAuthOrgMember(0, testutils.UserModels[0].Email),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamOrg(0),
		},
		BodyGenerators: []apitest.GenerateBodyParam{
			apitest.WithBodyParamOrgMemberID(0, testutils.UserModels[2].Email, "body_param_new_owner_member_id"),
		},
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}

func TestUpdateOrgOwnerSameOwner(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertResponseError(t, rr, http.StatusBadRequest, &types.APIErrors{
			Errors: []types.APIError{
				{
					Code:        types.ErrCodeBadRequest,
					Description: "new owner member id must be distinct from previous owner member id",
				},
			},
		})

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  fmt.Sprintf("/api/v1/organizations/{%s}/change_owner", string(types.URLParamOrgID)),
		RequestObj: types.UpdateOrgOwnerRequest{
			NewOwnerMemberID: "body_param_new_owner_member_id",
		},
		HandlerInit: orgs.NewOrgUpdateOwnerHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
			apitest.WithOrg(0),
			apitest.WithAuthOrgMember(0, testutils.UserModels[0].Email),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamOrg(0),
		},
		BodyGenerators: []apitest.GenerateBodyParam{
			apitest.WithBodyParamOrgMemberID(0, testutils.UserModels[0].Email, "body_param_new_owner_member_id"),
		},
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgAdditionalMemberAdmin)
}
