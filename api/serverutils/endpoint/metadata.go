package endpoint

import (
	"strings"

	"github.com/hatchet-dev/hatchet/api/v1/types"
)

// Path represents the path of an endpoint
type Path struct {
	// The parent path
	Parent *Path

	// The relative path from the parent
	RelativePath string
}

func (p *Path) GetPathString() (res string) {
	pathStrArr := make([]string, 0)

	currPath := p

	for currPath.Parent != nil {
		pathStrArr = append(pathStrArr, currPath.RelativePath)

		currPath = currPath.Parent
	}

	// append one more since we don't process the last parent
	pathStrArr = append(pathStrArr, currPath.RelativePath)

	pLen := len(pathStrArr)

	reversedPathStrArr := make([]string, pLen)

	for i, path := range pathStrArr {
		reversedPathStrArr[pLen-i-1] = path
	}

	return strings.Join(reversedPathStrArr, "")
}

// A context lookup key for determining which scopes are part of this endpoint
const RequestScopeCtxKey = "requestscopes"

// A context lookup key for retrieving the endpoint metadata
const EndpointMetadataCtxKey = "endpointmetadata"

// A context lookup key for attaching a websocket read writer to context
const RequestCtxWebsocketKey = "websocket"

// EndpointMetadata contains information about the endpoint, such as the API verb
// and scopes reqiured to call this endpoint successfully.
type EndpointMetadata struct {
	// The API verb that this endpoint should use
	Verb types.APIVerb

	// The underlying HTTP method that this endpoint should use
	Method types.HTTPVerb

	// The constructed path for the endpoint
	Path *Path

	// A list of scopes required for the endpoint
	Scopes []types.PermissionScope

	// Whether the endpoint upgrades to a websocket
	IsWebsocket bool

	// Whether the endpoint should check for a usage limit
	CheckUsage bool

	// The usage metric that the request should check for, if CheckUsage
	UsageMetric types.UsageMetric

	// AllowUnverifiedEmails controls whether or not the endpoint should permit authenticated
	// users with unverified email addresses
	AllowUnverifiedEmails bool
}

type RequestAction struct {
	Verb       types.APIVerb
	ResourceID string
}
