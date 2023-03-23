//go:build test

package apitest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/testutils"
)

// GenerateRequestCtx outputs key, value pairs to be set in context. Most handlers require additional
// scopes passed in context that can be referenced via their scope keys.
type GenerateRequestCtx func(t *testing.T, config *server.Config) (interface{}, interface{})

// GenerateURLParam takes in the server config and outputs URL params. It's meant to populate
// URL params that are only available from the init data.
type GenerateURLParam func(t *testing.T, config *server.Config, currParams map[string]string) map[string]string

// GenerateBodyParam returns a set of key-value string pairs. This is primarily meant for inserting
// dynamic IDs into the request body.
type GenerateBodyParam func(t *testing.T, config *server.Config) (string, string)

func WithAuthUser(initDataIndex uint) GenerateRequestCtx {
	return func(t *testing.T, config *server.Config) (interface{}, interface{}) {
		return types.UserScope, testutils.UserModels[initDataIndex]
	}
}

func WithOrg(initDataIndex uint) GenerateRequestCtx {
	return func(t *testing.T, config *server.Config) (interface{}, interface{}) {
		return types.OrgScope, testutils.OrgModels[initDataIndex]
	}
}

func WithAuthOrgMember(orgInitDataIndex uint, orgMemberEmail string) GenerateRequestCtx {
	return func(t *testing.T, config *server.Config) (interface{}, interface{}) {
		authOrgMember, err := config.DB.Repository.Org().ReadOrgMemberByUserOrInviteeEmail(testutils.OrgModels[orgInitDataIndex].ID, orgMemberEmail, false)

		if err != nil {
			panic(err)
		}

		return types.OrgMemberLookupKey, authOrgMember
	}
}

func WithOrgMember(initDataIndex uint) GenerateRequestCtx {
	return func(t *testing.T, config *server.Config) (interface{}, interface{}) {
		return types.OrgMemberScope, testutils.OrgAdditionalMembers[initDataIndex]
	}
}

func WithURLParamInviteID(orgInitDataIndex uint, inviteeEmail string) GenerateURLParam {
	return func(t *testing.T, config *server.Config, currParams map[string]string) map[string]string {
		orgMember, err := config.DB.Repository.Org().ReadOrgMemberByUserOrInviteeEmail(testutils.OrgModels[orgInitDataIndex].ID, inviteeEmail, false)

		if err != nil {
			panic(err)
		}

		currParams[string(types.URLParamOrgMemberInviteID)] = orgMember.InviteLink.ID

		return currParams
	}
}

func WithURLParamInviteTok(orgInitDataIndex uint, inviteeEmail string) GenerateURLParam {
	return func(t *testing.T, config *server.Config, currParams map[string]string) map[string]string {
		orgMember, _ := config.DB.Repository.Org().ReadOrgMemberByUserOrInviteeEmail(testutils.OrgModels[orgInitDataIndex].ID, inviteeEmail, false)
		invite := &orgMember.InviteLink

		invite.Decrypt(config.DB.GetEncryptionKey())

		currParams[string(types.URLParamOrgMemberInviteTok)] = testutils.OrgInviteLinksUnencryptedTok[invite.ID]

		return currParams
	}
}

func WithURLParamPAT(initDataIndex uint) GenerateURLParam {
	return func(t *testing.T, config *server.Config, currParams map[string]string) map[string]string {
		currParams[string(types.PersonalAccessTokenURLParam)] = testutils.PATModels[initDataIndex].ID

		return currParams
	}
}

func WithURLParamOrg(initDataIndex uint) GenerateURLParam {
	return func(t *testing.T, config *server.Config, currParams map[string]string) map[string]string {
		currParams[string(types.URLParamOrgID)] = testutils.OrgModels[initDataIndex].ID

		return currParams
	}
}

func WithBodyParamUserID(initDataIndex uint, varName string) GenerateBodyParam {
	return func(t *testing.T, config *server.Config) (string, string) {
		return varName, testutils.UserModels[initDataIndex].ID
	}
}

func WithBodyParamOrgMemberID(orgInitDataIndex uint, email string, varName string) GenerateBodyParam {
	return func(t *testing.T, config *server.Config) (string, string) {
		orgMember, err := config.DB.Repository.Org().ReadOrgMemberByUserOrInviteeEmail(testutils.OrgModels[orgInitDataIndex].ID, email, false)

		if err != nil {
			panic(err)
		}

		return varName, orgMember.ID
	}
}

func GetRequestAndRecorder(t *testing.T, method, route string, requestObj interface{}, config *server.Config, bodyGens ...GenerateBodyParam) (*http.Request, *httptest.ResponseRecorder) {
	var reader io.Reader = nil

	if requestObj != nil {
		data, err := json.Marshal(requestObj)

		if err != nil {
			t.Fatal(err)
		}

		// yes, this involves a lot of []byte -> string -> []byte conversions.
		// when this becomes an issue this should be made more efficient.
		for _, gen := range bodyGens {
			key, val := gen(t, config)

			// replace keys with generated results
			dataStr := strings.Replace(string(data), key, val, -1)

			data = []byte(dataStr)
		}

		reader = strings.NewReader(string(data))
	}

	// method and route don't actually matter since this is meant to test handlers
	req, err := http.NewRequest(method, route, reader)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	return req, rr
}

func WithURLParams(t *testing.T, req *http.Request, params map[string]string) *http.Request {
	rctx := chi.NewRouteContext()
	routeParams := &chi.RouteParams{}

	for key, val := range params {
		routeParams.Add(key, val)
	}

	rctx.URLParams = *routeParams

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	return req
}

type failingDecoderValidator struct {
	config *server.Config
}

func (f *failingDecoderValidator) DecodeAndValidate(
	w http.ResponseWriter,
	r *http.Request,
	v interface{},
) (ok bool) {
	apierrors.HandleAPIError(f.config.Logger, f.config.ErrorAlerter, w, r, apierrors.NewErrInternal(fmt.Errorf("fake error")), true)
	return false
}

func (f *failingDecoderValidator) DecodeAndValidateQueryOnly(
	w http.ResponseWriter,
	r *http.Request,
	v interface{},
) (ok bool) {
	apierrors.HandleAPIError(f.config.Logger, f.config.ErrorAlerter, w, r, apierrors.NewErrInternal(fmt.Errorf("fake error")), true)
	return false
}

func (f *failingDecoderValidator) DecodeAndValidateNoWrite(
	r *http.Request,
	v interface{},
) error {
	return fmt.Errorf("fake error")
}

func NewFailingDecoderValidator(config *server.Config) handlerutils.RequestDecoderValidator {
	return &failingDecoderValidator{config}
}
