package policies

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/mitchellh/mapstructure"
)

type PolicyInput struct {
	Endpoint EndpointInput `mapstructure:"endpoint"`

	User UserInput `mapstructure:"user"`
}

type EndpointInput struct {
	Verb string `mapstructure:"verb"`

	// The resources aren't typed because mapstructure was converting the typed values
	// to uppercase. Rather than investigate further, this was a simpler approach.
	Resources []map[string]interface{} `mapstructure:"resources"`
}

type UserInput struct {
	UserAccountKind string `mapstructure:"user_account_kind"`
}

func GetInputFromRequest(r *http.Request) map[string]interface{} {
	user, _ := r.Context().Value(types.UserScope).(*models.User)
	endpointMeta, _ := r.Context().Value(endpoint.EndpointMetadataCtxKey).(*endpoint.EndpointMetadata)
	reqScopes, _ := r.Context().Value(endpoint.RequestScopeCtxKey).(map[types.PermissionScope]*endpoint.RequestAction)

	structuredRes := &PolicyInput{
		User: UserInput{
			UserAccountKind: string(user.UserAccountKind),
		},
	}

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
