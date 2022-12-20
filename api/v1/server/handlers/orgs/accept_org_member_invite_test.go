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

func TestAcceptOrgMemberInviteSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		// ensure that response is 200
		assert.Equal(t, http.StatusOK, rr.Result().StatusCode, "status should be 200")

		// query the database directly to ensure that user was added
		orgMembers, _, err := config.DB.Repository.Org().ListOrgMembersByOrgID(testutils.OrgModels[0].ID)

		if err != nil {
			t.Fatalf("%v", err)
		}

		assert.Equal(t, 2, len(orgMembers))

		var newOrgMember *models.OrganizationMember

		for _, orgMember := range orgMembers {
			if orgMember.User.Email == testutils.UserModels[2].Email {
				newOrgMember = orgMember
			}
		}

		assert.NotNil(t, newOrgMember)
		assert.Equal(t, testutils.UserModels[2].ID, newOrgMember.UserID, "user ids should be equal")

		return nil
	}, &apitest.APITesterOpts{
		Method:      "POST",
		Route:       fmt.Sprintf("/api/v1/invites/{%s}/{%s}", types.URLParamOrgMemberInviteID, types.URLParamOrgMemberInviteTok),
		HandlerInit: orgs.NewOrgAcceptMemberInviteHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(2),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamInviteID(0, testutils.OrgInviteLinks[0].InviteeEmail),
			apitest.WithURLParamInviteTok(0, testutils.OrgInviteLinks[0].InviteeEmail),
		},
	}, testutils.InitUsers, testutils.InitOrgs, testutils.InitOrgInviteLinks)
}
