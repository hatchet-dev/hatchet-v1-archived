package router

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/orgs"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// swagger:parameters getOrganization createOrgMemberInvite updateOrgOwner updateOrganization leaveOrg createTeam listTeams
type orgPathParams struct {
	// The org id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Org string `json:"org_id"`
}

// swagger:parameters deleteOrgMember updateOrgMemberPolicies getOrgMember
type orgMemberPathParams struct {
	// The org id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Org string `json:"org_id"`

	// The org member id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	OrgMember string `json:"org_member_id"`
}

// swagger:parameters acceptOrgMemberInvite
type inviteAcceptPathParams struct {
	// The member invite id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	OrgMemberInviteID string `json:"org_member_invite_id"`

	// The member invite token (sensitive)
	// in: path
	// required: true
	// example: abcdefgh...
	OrgMemberInviteTok string `json:"org_member_invite_tok"`
}

func NewOrgRouteRegisterer(children ...*router.Registerer) *router.Registerer {
	return &router.Registerer{
		GetRoutes: GetOrganizationRoutes,
		Children:  children,
	}
}

func GetOrganizationRoutes(
	r chi.Router,
	config *server.Config,
	basePath *endpoint.Path,
	factory endpoint.APIEndpointFactory,
	children ...*router.Registerer,
) []*router.Route {
	routes := make([]*router.Route, 0)

	// POST /api/v1/organizations -> orgs.NewOrgCreateHandler
	// swagger:operation POST /api/v1/organizations createOrganization
	//
	// ### Description
	//
	// Creates a new organization, with the authenticated user set as the organization owner.
	//
	// ---
	// produces:
	// - application/json
	// summary: Create a new organization
	// tags:
	// - Organizations
	// parameters:
	//   - in: body
	//     name: CreateOrganizationRequest
	//     description: The organization to create
	//     schema:
	//       $ref: '#/definitions/CreateOrganizationRequest'
	// responses:
	//   '201':
	//     description: Successfully created the organization
	//     schema:
	//       $ref: '#/definitions/CreateOrganizationResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	//   '405':
	//     description: This endpoint is not supported on this Hatchet instance.
	//     schema:
	//       $ref: '#/definitions/APIErrorNotSupportedExample'
	createOrgEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/organizations",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	createOrgHandler := orgs.NewOrgCreateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: createOrgEndpoint,
		Handler:  createOrgHandler,
		Router:   r,
	})

	// GET /api/v1/organizations/{org_id} -> orgs.NewOrgGetHandler
	// swagger:operation GET /api/v1/organizations/{org_id} getOrganization
	//
	// ### Description
	//
	// Retrieves an organization by the `org_id`.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get an organization
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	// responses:
	//   '200':
	//     description: Successfully got the organization
	//     schema:
	//       $ref: '#/definitions/GetOrganizationResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	//   '405':
	//     description: This endpoint is not supported on this Hatchet instance.
	//     schema:
	//       $ref: '#/definitions/APIErrorNotSupportedExample'
	getOrgEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}", string(types.URLParamOrgID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
			},
		},
	)

	getOrgHandler := orgs.NewOrgGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getOrgEndpoint,
		Handler:  getOrgHandler,
		Router:   r,
	})

	// POST /api/v1/organizations/{org_id} -> orgs.NewOrgUpdateHandler
	// swagger:operation POST /api/v1/organizations/{org_id} updateOrganization
	//
	// ### Description
	//
	// Updates organization metadata.
	//
	// ---
	// produces:
	// - application/json
	// summary: Update an organization
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	//   - in: body
	//     name: UpdateOrgRequest
	//     description: The values to update
	//     schema:
	//       $ref: '#/definitions/UpdateOrgRequest'
	// responses:
	//   '200':
	//     description: Successfully updated the organization
	//     schema:
	//       $ref: '#/definitions/UpdateOrgResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	updateOrgEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}", string(types.URLParamOrgID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
			},
		},
	)

	updateOrgHandler := orgs.NewOrgUpdateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: updateOrgEndpoint,
		Handler:  updateOrgHandler,
		Router:   r,
	})

	// POST /api/v1/organizations/{org_id}/members -> orgs.NewOrgCreateMemberInviteHandler
	// swagger:operation POST /api/v1/organizations/{org_id}/members createOrgMemberInvite
	//
	// ### Description
	//
	// Creates a new invite for an organization member.
	//
	// ---
	// produces:
	// - application/json
	// summary: Create a member invite
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	//   - in: body
	//     name: CreateOrgMemberInviteRequest
	//     description: The org member to create
	//     schema:
	//       $ref: '#/definitions/CreateOrgMemberInviteRequest'
	// responses:
	//   '201':
	//     description: Successfully created the invite.
	//     schema:
	//       $ref: '#/definitions/CreateOrgMemberInviteResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	//   '405':
	//     description: This endpoint is not supported on this Hatchet instance.
	//     schema:
	//       $ref: '#/definitions/APIErrorNotSupportedExample'
	createOrgMemberInviteEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}/members", string(types.URLParamOrgID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
			},
		},
	)

	createOrgMemberInviteHandler := orgs.NewOrgCreateMemberInviteHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: createOrgMemberInviteEndpoint,
		Handler:  createOrgMemberInviteHandler,
		Router:   r,
	})

	// POST /api/v1/invites/{org_member_invite_id}/{org_member_invite_tok} -> orgs.NewOrgAcceptMemberInviteHandler
	// swagger:operation POST /api/v1/invites/{org_member_invite_id}/{org_member_invite_tok} acceptOrgMemberInvite
	//
	// ### Description
	//
	// Accept an invite for an organization.
	//
	// ---
	// produces:
	// - application/json
	// summary: Accept an organization invite.
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_member_invite_id
	//   - name: org_member_invite_tok
	// responses:
	//   '200':
	//     description: Successfully accepted the invite.
	//     schema:
	//       $ref: '#/definitions/EmptyResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	//   '405':
	//     description: This endpoint is not supported on this Hatchet instance.
	//     schema:
	//       $ref: '#/definitions/APIErrorNotSupportedExample'
	acceptOrgMemberInviteEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/invites/{%s}/{%s}", string(types.URLParamOrgMemberInviteID), string(types.URLParamOrgMemberInviteTok)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	acceptOrgMemberInviteHandler := orgs.NewOrgAcceptMemberInviteHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: acceptOrgMemberInviteEndpoint,
		Handler:  acceptOrgMemberInviteHandler,
		Router:   r,
	})

	// GET /api/v1/organizations/{org_id}/members -> orgs.NewOrgListMembersHandler
	// swagger:operation GET /api/v1/organizations/{org_id}/members listOrgMembers
	//
	// ### Description
	//
	// Lists organization members.
	//
	// ---
	// produces:
	// - application/json
	// summary: List organization members
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	// responses:
	//   '200':
	//     description: Successfully listed members
	//     schema:
	//       $ref: '#/definitions/ListOrgMembersResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	listOrgMembersEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbList,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}/members", string(types.URLParamOrgID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
			},
		},
	)

	listOrgMembersHandler := orgs.NewOrgListMembersHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: listOrgMembersEndpoint,
		Handler:  listOrgMembersHandler,
		Router:   r,
	})

	// GET /api/v1/organizations/{org_id}/members/{org_member_id} -> orgs.NewOrgGetMemberHandler
	// swagger:operation GET /api/v1/organizations/{org_id}/members/{org_member_id} getOrgMember
	//
	// ### Description
	//
	// Get organization member. Only admins and owner can read full member data.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get organization member.
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	//   - name: org_member_id
	// responses:
	//   '200':
	//     description: Successfully got member
	//     schema:
	//       $ref: '#/definitions/GetOrgMemberResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	getOrgMemberEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}/members/{%s}", string(types.URLParamOrgID), string(types.URLParamOrgMemberID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
				types.OrgMemberScope,
			},
		},
	)

	getOrgMemberHandler := orgs.NewOrgGetMemberHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getOrgMemberEndpoint,
		Handler:  getOrgMemberHandler,
		Router:   r,
	})

	// POST /api/v1/organizations/{org_id}/change_owner -> orgs.NewOrgListMembersHandler
	// swagger:operation POST /api/v1/organizations/{org_id}/change_owner updateOrgOwner
	//
	// ### Description
	//
	// Update organization owner. Only owners may update organization owners. The previous owner will become
	// an admin (and can subsequently be removed from the organization, if required).
	//
	// ---
	// produces:
	// - application/json
	// summary: Update organization owner.
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	// responses:
	//   '200':
	//     description: Successfully changed organization owner.
	//     schema:
	//       $ref: '#/definitions/EmptyResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	updateOrgOwnerEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}/change_owner", string(types.URLParamOrgID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
				types.OrgOwnerScope,
			},
		},
	)

	updateOrgOwnerHandler := orgs.NewOrgUpdateOwnerHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: updateOrgOwnerEndpoint,
		Handler:  updateOrgOwnerHandler,
		Router:   r,
	})

	// POST /api/v1/organizations/{org_id}/members/{org_member_id}/update_policies -> orgs.NewOrgUpdateMemberPoliciesHandler
	// swagger:operation POST /api/v1/organizations/{org_id}/members/{org_member_id}/update_policies updateOrgMemberPolicies
	//
	// ### Description
	//
	// Update an organization member's policies.
	//
	// ---
	// produces:
	// - application/json
	// summary: Update organization member policies.
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	//   - name: org_member_id
	// responses:
	//   '202':
	//     description: Successfully updated organization member policies.
	//     schema:
	//       $ref: '#/definitions/UpdateOrgMemberPoliciesResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	updateOrgMemberPoliciesEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}/members/{%s}/update_policies", string(types.URLParamOrgID), string(types.URLParamOrgMemberID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
				types.OrgMemberScope,
			},
		},
	)

	updateOrgMemberPoliciesHandler := orgs.NewOrgUpdateMemberPoliciesHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: updateOrgMemberPoliciesEndpoint,
		Handler:  updateOrgMemberPoliciesHandler,
		Router:   r,
	})

	// POST /api/v1/organizations/{org_id}/leave -> orgs.NewOrgLeaveHandler
	// swagger:operation POST /api/v1/organizations/{org_id}/leave leaveOrg
	//
	// ### Description
	//
	// Leave an organization. The currently authenticated user will leave this organization.
	// Owners cannot leave an organization without changing the owner first.
	//
	// ---
	// produces:
	// - application/json
	// summary: Leave an organization
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	// responses:
	//   '202':
	//     description: Successfully triggered organization member removal.
	//     schema:
	//       $ref: '#/definitions/EmptyResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	leaveOrgEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}/leave", string(types.URLParamOrgID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
			},
		},
	)

	leaveOrgHandler := orgs.NewOrgLeaveHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: leaveOrgEndpoint,
		Handler:  leaveOrgHandler,
		Router:   r,
	})

	// DELETE /api/v1/organizations/{org_id}/members/{org_member_id} -> orgs.NewOrgDeleteMemberHandler
	// swagger:operation DELETE /api/v1/organizations/{org_id}/members/{org_member_id} deleteOrgMember
	//
	// ### Description
	//
	// Delete an organization member. Only admins can delete an organization member. Owners cannot be
	// removed from the organization, the owner must be transferred before the organization owner can
	// be removed.
	//
	// ---
	// produces:
	// - application/json
	// summary: Delete organization member.
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	// responses:
	//   '200':
	//     description: Successfully triggered organization member deletion.
	//     schema:
	//       $ref: '#/definitions/EmptyResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	deleteOrgMemberEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbDelete,
			Method: types.HTTPVerbDelete,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}/members/{%s}", string(types.URLParamOrgID), string(types.URLParamOrgMemberID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
				types.OrgMemberScope,
			},
		},
	)

	deleteOrgMemberHandler := orgs.NewOrgDeleteMemberHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: deleteOrgMemberEndpoint,
		Handler:  deleteOrgMemberHandler,
		Router:   r,
	})

	// DELETE /api/v1/organizations/{org_id} -> orgs.NewOrgDeleteHandler
	// swagger:operation DELETE /api/v1/organizations/{org_id} deleteOrg
	//
	// ### Description
	//
	// Delete an organization. Only owners can delete organizations.
	//
	// ---
	// produces:
	// - application/json
	// summary: Delete organization.
	// tags:
	// - Organizations
	// parameters:
	//   - name: org_id
	// responses:
	//   '200':
	//     description: Successfully triggered organization deletion.
	//     schema:
	//       $ref: '#/definitions/DeleteOrganizationResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	deleteOrgEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbDelete,
			Method: types.HTTPVerbDelete,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}", string(types.URLParamOrgID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
				types.OrgOwnerScope,
			},
		},
	)

	deleteOrgHandler := orgs.NewOrgDeleteHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: deleteOrgEndpoint,
		Handler:  deleteOrgHandler,
		Router:   r,
	})

	return routes
}
