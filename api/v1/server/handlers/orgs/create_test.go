package orgs_test

import (
	"encoding/json"
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

func TestCreateOrgSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		gotOrgResponse := &types.CreateOrganizationResponse{}

		err := json.NewDecoder(rr.Body).Decode(gotOrgResponse)

		if err != nil {
			t.Fatal(err)
		}

		// assert that the uuid on the pat field is valid
		assert.True(t, uuidutils.IsValidUUID(gotOrgResponse.ID))

		assert.Equal(t, "My Org 1", gotOrgResponse.DisplayName)
		assert.Equal(t, "user1@gmail.com", gotOrgResponse.Owner.Email)

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  "/api/v1/organizations",
		RequestObj: types.CreateOrganizationRequest{
			DisplayName: "My Org 1",
		},
		HandlerInit: orgs.NewOrgCreateHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
	}, testutils.InitUsers)
}

func TestCreateOrgMissingDisplayName(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		apitest.AssertResponseError(t, rr, http.StatusBadRequest, &types.APIErrors{
			Errors: []types.APIError{
				{
					Code:        types.ErrCodeBadRequest,
					Description: "validation failed on field 'DisplayName' on condition 'required'",
				},
			},
		})

		return nil
	}, &apitest.APITesterOpts{
		Method:      "POST",
		Route:       "/api/v1/organizations",
		RequestObj:  types.CreateOrganizationRequest{},
		HandlerInit: orgs.NewOrgCreateHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
		},
	}, testutils.InitUsers)
}
