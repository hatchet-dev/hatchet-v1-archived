package endpoint

import (
	"context"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

type EndpointMiddleware struct {
	config       *server.Config
	endpointMeta *EndpointMetadata
}

func NewEndpointMiddleware(
	config *server.Config,
	endpoint *EndpointMetadata,
) *EndpointMiddleware {
	return &EndpointMiddleware{config, endpoint}
}

func (p *EndpointMiddleware) Middleware(next http.Handler) http.Handler {
	return &EndpointMiddlewareHandler{next, p.config, p.endpointMeta}
}

type EndpointMiddlewareHandler struct {
	next         http.Handler
	config       *server.Config
	endpointMeta *EndpointMetadata
}

func (h *EndpointMiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// add the endpoint metadata to the request context
	ctx := NewEndpointMetaCtx(r.Context(), h.endpointMeta)
	r = r.Clone(ctx)

	// get the full map of scopes to resource actions
	reqScopes, reqErr := getRequestActionForEndpoint(r, h.endpointMeta)

	if reqErr != nil {
		apierrors.HandleAPIError(h.config.Logger, h.config.ErrorAlerter, w, r, reqErr, true)
		return
	}

	// add the set of resource ids to the request context
	ctx = NewRequestScopeCtx(r.Context(), reqScopes)
	r = r.Clone(ctx)
	h.next.ServeHTTP(w, r)
}

func NewEndpointMetaCtx(ctx context.Context, endpointMeta *EndpointMetadata) context.Context {
	return context.WithValue(ctx, EndpointMetadataCtxKey, endpointMeta)
}

func NewRequestScopeCtx(ctx context.Context, reqScopes map[types.PermissionScope]*RequestAction) context.Context {
	return context.WithValue(ctx, RequestScopeCtxKey, reqScopes)
}

func getRequestActionForEndpoint(
	r *http.Request,
	endpointMeta *EndpointMetadata,
) (res map[types.PermissionScope]*RequestAction, reqErr apierrors.RequestError) {
	res = make(map[types.PermissionScope]*RequestAction)

	var resourceID string

	// iterate through scopes, attach policies as needed
	for _, scope := range endpointMeta.Scopes {
		// find the resource ID and create the resource
		switch scope {
		case types.OrgScope:
			resourceID, reqErr = handlerutils.GetURLParamString(r, types.URLParamOrgID)
		case types.OrgOwnerScope:
			// in this case, there is no unique resource ID for the "owner" concept, so we simply attach a
			// resource ID of "owner" to the request
			resourceID = "owner"
		case types.OrgMemberScope:
			resourceID, reqErr = handlerutils.GetURLParamString(r, types.URLParamOrgMemberID)
		case types.TeamScope:
			resourceID, reqErr = handlerutils.GetURLParamString(r, types.URLParamTeamID)
		case types.TeamMemberScope:
			resourceID, reqErr = handlerutils.GetURLParamString(r, types.URLParamTeamMemberID)
		}

		if reqErr != nil {
			return nil, reqErr
		}

		res[scope] = &RequestAction{
			Verb:       endpointMeta.Verb,
			ResourceID: resourceID,
		}
	}

	return res, nil
}
