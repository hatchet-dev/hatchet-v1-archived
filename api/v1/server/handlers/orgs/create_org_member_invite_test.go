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
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrgMemberInviteSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		gotResponse := &types.CreateOrgMemberInviteResponse{}

		err := json.NewDecoder(rr.Body).Decode(gotResponse)

		if err != nil {
			t.Fatal(err)
		}

		// assert that the uuid on the pat field is valid
		assert.True(t, uuidutils.IsValidUUID(gotResponse.ID))

		// assert that the user field is empty, as the invite is not accepted yet
		assert.Equal(t, "", gotResponse.User.DisplayName)
		assert.Equal(t, "", gotResponse.User.Email)
		assert.False(t, gotResponse.InviteAccepted, "invite accepted should be false")
		assert.Equal(t, 1, len(gotResponse.OrgPolicies))

		// assert that invite field has correct email
		assert.Equal(t, testutils.UserModels[1].Email, gotResponse.Invite.InviteeEmail, "invitee email should match second user")

		// assert that the org policy is an admin
		assert.Equal(t, "admin", gotResponse.OrgPolicies[0].Name)

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  fmt.Sprintf("/api/v1/organizations/{%s}/members", types.URLParamOrgID),
		RequestObj: types.CreateOrgMemberInviteRequest{
			InviteeEmail: testutils.UserModels[1].Email,
			InviteePolicies: []types.OrganizationPolicyReference{
				{
					Name: "admin",
				},
			},
		},
		HandlerInit: orgs.NewOrgCreateMemberInviteHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
			apitest.WithOrg(0),
		},
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestCreateOrgMemberInviteAlreadyExists(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertResponseError(t, rr, http.StatusBadRequest, &types.APIErrors{
			Errors: []types.APIError{
				{
					Code:        types.ErrCodeBadRequest,
					Description: "There is already an organization member with this email address.",
				},
			},
		})

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  fmt.Sprintf("/api/v1/organizations/{%s}/members", types.URLParamOrgID),
		RequestObj: types.CreateOrgMemberInviteRequest{
			InviteeEmail: testutils.UserModels[0].Email,
			InviteePolicies: []types.OrganizationPolicyReference{
				{
					Name: "admin",
				},
			},
		},
		HandlerInit: orgs.NewOrgCreateMemberInviteHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
			apitest.WithOrg(0),
		},
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestCreateOrgMemberInviteInvalidEmail(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertResponseError(t, rr, http.StatusBadRequest, &types.APIErrors{
			Errors: []types.APIError{
				{
					Code:        types.ErrCodeBadRequest,
					Description: "Invalid email address",
				},
			},
		})

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  fmt.Sprintf("/api/v1/organizations/{%s}/members", types.URLParamOrgID),
		RequestObj: types.CreateOrgMemberInviteRequest{
			InviteeEmail: "invalidemail",
			InviteePolicies: []types.OrganizationPolicyReference{
				{
					Name: "admin",
				},
			},
		},
		HandlerInit: orgs.NewOrgCreateMemberInviteHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
			apitest.WithOrg(0),
		},
	}, testutils.InitUsers, testutils.InitOrgs)
}

func TestCreateOrgMemberInviteNoPolicies(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertResponseError(t, rr, http.StatusBadRequest, &types.APIErrors{
			Errors: []types.APIError{
				{
					Code:        types.ErrCodeBadRequest,
					Description: "validation failed on field 'InviteePolicies' on condition 'min' [ 1 ]: got invalid type",
				},
			},
		})

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  fmt.Sprintf("/api/v1/organizations/{%s}/members", types.URLParamOrgID),
		RequestObj: types.CreateOrgMemberInviteRequest{
			InviteeEmail:    testutils.UserModels[1].Email,
			InviteePolicies: []types.OrganizationPolicyReference{},
		},
		HandlerInit: orgs.NewOrgCreateMemberInviteHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
			apitest.WithOrg(0),
		},
	}, testutils.InitUsers, testutils.InitOrgs)
}
