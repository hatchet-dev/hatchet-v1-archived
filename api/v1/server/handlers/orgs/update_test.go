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

func TestUpdateOrgSuccessful(t *testing.T) {
	apitest.RunAPITest(t, func(config *server.Config, rr *httptest.ResponseRecorder, req *http.Request) error {
		updateOrgResponse := &types.UpdateOrgResponse{}

		err := json.NewDecoder(rr.Body).Decode(updateOrgResponse)

		if err != nil {
			t.Fatal(err)
		}

		// assert that the uuid on the pat field is valid
		assert.True(t, uuidutils.IsValidUUID(updateOrgResponse.ID))

		assert.Equal(t, "My Org 1 - Rename", updateOrgResponse.DisplayName)
		assert.Equal(t, "user1@gmail.com", updateOrgResponse.Owner.Email)

		return nil
	}, &apitest.APITesterOpts{
		Method: "POST",
		Route:  fmt.Sprintf("/api/v1/organizations/{%s}", string(types.URLParamOrgID)),
		RequestObj: types.UpdateOrgRequest{
			DisplayName: "My Org 1 - Rename",
		},
		HandlerInit: orgs.NewOrgUpdateHandler,
		CtxGenerators: []apitest.GenerateRequestCtx{
			apitest.WithAuthUser(0),
			apitest.WithOrg(0),
		},
		URLGenerators: []apitest.GenerateURLParam{
			apitest.WithURLParamOrg(0),
		},
	}, testutils.InitUsers, testutils.InitOrgs)
}
