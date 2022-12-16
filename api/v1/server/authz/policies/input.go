package policies

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/mitchellh/mapstructure"
)

type OrgPolicyInput struct {
	Endpoint EndpointInput `json:"endpoint"`
}

type EndpointInput struct {
	Verb      string                  `json:"verb"`
	Resources []EndpointResourceInput `json:"resources"`
}

type EndpointResourceInput struct {
	Verb   string      `json:"verb"`
	Scope  string      `json:"scope"`
	Target interface{} `json:"target"`
}

func GetInputFromRequest(r *http.Request) map[string]interface{} {
	endpointMeta, _ := r.Context().Value(endpoint.EndpointMetadataCtxKey).(*endpoint.EndpointMetadata)
	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)

	structuredRes := &OrgPolicyInput{}

	structuredRes.Endpoint = EndpointInput{
		Verb:      string(endpointMeta.Verb),
		Resources: make([]EndpointResourceInput, 0),
	}

	for scopeName, scope := range reqScopes {
		structuredRes.Endpoint.Resources = append(structuredRes.Endpoint.Resources, EndpointResourceInput{
			Scope:  string(scopeName),
			Verb:   string(scope.Verb),
			Target: scope.ResourceID,
		})
	}

	// marshal to json
	res := make(map[string]interface{})

	mapstructure.Decode(structuredRes, res)

	return res
}
