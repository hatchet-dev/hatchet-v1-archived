package policies

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/mitchellh/mapstructure"
)

type OrgPolicyInput struct {
	Endpoint EndpointInput `mapstructure:"endpoint"`
}

type EndpointInput struct {
	Verb string `mapstructure:"verb"`

	// The resources aren't typed because mapstructure was converting the typed values
	// to uppercase. Rather than investigate further, this was a simpler approach.
	Resources []map[string]interface{} `mapstructure:"resources"`
}

func GetInputFromRequest(r *http.Request) map[string]interface{} {
	endpointMeta, _ := r.Context().Value(endpoint.EndpointMetadataCtxKey).(*endpoint.EndpointMetadata)
	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)

	structuredRes := &OrgPolicyInput{}

	structuredRes.Endpoint = EndpointInput{
		Verb:      string(endpointMeta.Verb),
		Resources: make([]map[string]interface{}, 0),
	}

	for scopeName, scope := range reqScopes {
		structuredRes.Endpoint.Resources = append(structuredRes.Endpoint.Resources, map[string]interface{}{
			"scope":  string(scopeName),
			"verb":   string(scope.Verb),
			"target": scope.ResourceID,
		})
	}

	res := make(map[string]interface{})

	mapstructure.Decode(structuredRes, &res)

	return res
}
