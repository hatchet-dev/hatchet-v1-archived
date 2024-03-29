package router

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/github_app"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/notifications"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/teams"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// swagger:parameters listTeamMembers addTeamMember updateTeam deleteTeam createModule listModules createMonitor listMonitors
type teamPathParams struct {
	// The team id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Team string `json:"team_id"`
}

// swagger:parameters deleteTeamMember
type teamMemberPathParams struct {
	// The team member id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	TeamMember string `json:"team_member_id"`
}

// swagger:parameters githubIncomingWebhook
type githubIncomingWebhookParams struct {
	// The team id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Team string `json:"team_id"`

	// The incoming webhook id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	WebhookID string `json:"github_incoming_webhook_id"`
}

// swagger:parameters getNotification
type notifPathParams struct {
	// The team id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Team string `json:"team_id"`

	// The notification id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	Notification string `json:"notification_id"`
}

func NewTeamRouteRegisterer(children ...*router.Registerer) *router.Registerer {
	return &router.Registerer{
		GetRoutes: GetTeamRoutes,
		Children:  children,
	}
}

func GetTeamRoutes(
	r chi.Router,
	config *server.Config,
	basePath *endpoint.Path,
	factory endpoint.APIEndpointFactory,
	children ...*router.Registerer,
) []*router.Route {
	routes := make([]*router.Route, 0)

	// POST /api/v1/organizations/{org_id}/teams -> teams.NewTeamCreateHandler
	// swagger:operation POST /api/v1/organizations/{org_id}/teams createTeam
	//
	// ### Description
	//
	// Creates a new team, with the authenticated user set as a team admin.
	//
	// ---
	// produces:
	// - application/json
	// summary: Create a new team.
	// tags:
	// - Teams
	// parameters:
	//   - name: org_id
	//   - in: body
	//     name: CreateTeamRequest
	//     description: The team to create
	//     schema:
	//       $ref: '#/definitions/CreateTeamRequest'
	// responses:
	//   '201':
	//     description: Successfully created the team
	//     schema:
	//       $ref: '#/definitions/CreateTeamResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	createTeamEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}/teams", types.URLParamOrgID),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
			},
		},
	)

	createTeamHandler := teams.NewTeamCreateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: createTeamEndpoint,
		Handler:  createTeamHandler,
		Router:   r,
	})

	// GET /api/v1/organizations/{org_id}/teams -> teams.NewTeamListHandler
	// swagger:operation GET /api/v1/organizations/{org_id}/teams listTeams
	//
	// ### Description
	//
	// Lists teams for an organization.
	//
	// ---
	// produces:
	// - application/json
	// summary: List teams.
	// tags:
	// - Teams
	// parameters:
	//   - name: org_id
	// responses:
	//   '200':
	//     description: Successfully listed teams
	//     schema:
	//       $ref: '#/definitions/ListTeamsResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	listTeamsEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbList,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/organizations/{%s}/teams", string(types.URLParamOrgID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.OrgScope,
			},
		},
	)

	listTeamsHandler := teams.NewTeamListHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: listTeamsEndpoint,
		Handler:  listTeamsHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/members -> teams.NewTeamListMemberHandler
	// swagger:operation GET /api/v1/teams/{team_id}/members listTeamMembers
	//
	// ### Description
	//
	// Lists team members for a team.
	//
	// ---
	// produces:
	// - application/json
	// summary: List team members
	// tags:
	// - Teams
	// parameters:
	//   - name: team_id
	// responses:
	//   '200':
	//     description: Successfully listed team members
	//     schema:
	//       $ref: '#/definitions/ListTeamMembersResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	listTeamMembersEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbList,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/teams/{%s}/members", string(types.URLParamTeamID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
			},
		},
	)

	listTeamMembersHandler := teams.NewTeamListMemberHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: listTeamMembersEndpoint,
		Handler:  listTeamMembersHandler,
		Router:   r,
	})

	// POST /api/v1/teams/{team_id}/members -> teams.NewTeamAddMemberHandler
	// swagger:operation POST /api/v1/teams/{team_id}/members addTeamMember
	//
	// ### Description
	//
	// Add a team member from the organization members to the team.
	//
	// ---
	// produces:
	// - application/json
	// summary: Add team member
	// tags:
	// - Teams
	// parameters:
	//   - name: team_id
	//   - in: body
	//     name: TeamAddMemberRequest
	//     description: The team member to add
	//     schema:
	//       $ref: '#/definitions/TeamAddMemberRequest'
	// responses:
	//   '200':
	//     description: Successfully added team member
	//     schema:
	//       $ref: '#/definitions/TeamAddMemberResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	addTeamMembersEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/teams/{%s}/members", string(types.URLParamTeamID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
			},
		},
	)

	addTeamMembersHandler := teams.NewTeamAddMemberHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: addTeamMembersEndpoint,
		Handler:  addTeamMembersHandler,
		Router:   r,
	})

	// POST /api/v1/teams/{team_id} -> teams.NewTeamUpdateHandler
	// swagger:operation POST /api/v1/teams/{team_id} updateTeam
	//
	// ### Description
	//
	// Updates a team.
	//
	// ---
	// produces:
	// - application/json
	// summary: Update team
	// tags:
	// - Teams
	// parameters:
	//   - name: team_id
	//   - in: body
	//     name: TeamUpdateRequest
	//     description: The team parameters to update
	//     schema:
	//       $ref: '#/definitions/TeamUpdateRequest'
	// responses:
	//   '200':
	//     description: Successfully updated the team
	//     schema:
	//       $ref: '#/definitions/TeamUpdateResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	updateTeamEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/teams/{%s}", string(types.URLParamTeamID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
			},
		},
	)

	updateTeamHandler := teams.NewTeamUpdateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: updateTeamEndpoint,
		Handler:  updateTeamHandler,
		Router:   r,
	})

	// DELETE /api/v1/teams/{team_id} -> teams.NewTeamDeleteHandler
	// swagger:operation DELETE /api/v1/teams/{team_id} deleteTeam
	//
	// ### Description
	//
	// Delete a team. This operation cannot be undone.
	//
	// ---
	// produces:
	// - application/json
	// summary: Delete team.
	// tags:
	// - Teams
	// parameters:
	//   - name: team_id
	// responses:
	//   '200':
	//     description: Successfully triggered team deletion.
	//     schema:
	//       $ref: '#/definitions/DeleteTeamResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	deleteTeamEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbDelete,
			Method: types.HTTPVerbDelete,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/teams/{%s}", string(types.URLParamTeamID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
			},
		},
	)

	deleteTeamHandler := teams.NewTeamDeleteHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: deleteTeamEndpoint,
		Handler:  deleteTeamHandler,
		Router:   r,
	})

	// DELETE /api/v1/teams/{team_id}/members/{team_member_id} -> teams.NewTeamRemoveMemberHandler
	// swagger:operation DELETE /api/v1/teams/{team_id}/members/{team_member_id} deleteTeamMember
	//
	// ### Description
	//
	// Delete a team member.
	//
	// ---
	// produces:
	// - application/json
	// summary: Delete team member
	// tags:
	// - Teams
	// parameters:
	//   - name: team_id
	// responses:
	//   '201':
	//     description: Successfully triggered team member deletion.
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
	deleteTeamMemberEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbDelete,
			Method: types.HTTPVerbDelete,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/teams/{%s}/members/{%s}", string(types.URLParamTeamID), string(types.URLParamTeamMemberID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.TeamMemberScope,
			},
		},
	)

	deleteTeamMemberHandler := teams.NewTeamRemoveMemberHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: deleteTeamMemberEndpoint,
		Handler:  deleteTeamMemberHandler,
		Router:   r,
	})

	// POST /api/v1/teams/{team_id}/github_incoming/{github_incoming_webhook_id} -> teams.NewGithubIncomingWebhookHandler
	// swagger:operation POST /api/v1/teams/{team_id}/github_incoming/{github_incoming_webhook_id} githubIncomingWebhook
	//
	// ### Description
	//
	// Github incoming webhook handler.
	//
	// ---
	// produces:
	// - application/json
	// summary: Github incoming webhook endpoint
	// tags:
	// - Teams
	// parameters:
	//   - name: team_id
	//   - name: github_incoming_webhook_id
	// responses:
	//   '200':
	//     description: Successfully responded to webhook
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	githubIncomingWebhookEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/teams/{%s}/github_incoming/{%s}", string(types.URLParamTeamID), string(types.URLParamGithubWebhookID)),
			},
			// Incoming webhook has no permission scopes because it uses its own auth mechanism
			Scopes: []types.PermissionScope{},
		},
	)

	githubIncomingWebhookHandler := github_app.NewGithubIncomingWebhookHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: githubIncomingWebhookEndpoint,
		Handler:  githubIncomingWebhookHandler,
		Router:   r,
	})

	// GET /api/v1/teams/{team_id}/notifications/{notification_id} -> notifications.NewNotificationGetHandler
	// swagger:operation GET /api/v1/teams/{team_id}/notifications/{notification_id} getNotification
	//
	// ### Description
	//
	// Gets a notification by id.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get notification
	// tags:
	// - Teams
	// parameters:
	//   - name: team_id
	// responses:
	//   '200':
	//     description: Successfully got the notification
	//     schema:
	//       $ref: '#/definitions/GetNotificationResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	getNotificationEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/teams/{%s}/notifications/{%s}", string(types.URLParamTeamID), string(types.URLParamNotificationID)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
				types.TeamScope,
				types.NotificationScope,
			},
		},
	)

	getNotificationHandler := notifications.NewNotificationGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getNotificationEndpoint,
		Handler:  getNotificationHandler,
		Router:   r,
	})

	return routes
}
