package router

import (
	"fmt"

	"github.com/go-chi/chi"
	"github.com/hatchet-dev/hatchet/api/serverutils/endpoint"
	"github.com/hatchet-dev/hatchet/api/serverutils/router"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/pats"
	"github.com/hatchet-dev/hatchet/api/v1/server/handlers/users"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// swagger:parameters getPersonalAccessToken revokePersonalAccessToken deletePersonalAccessToken
type patPathParams struct {
	// The personal access token id
	// in: path
	// required: true
	// example: 322346f9-54b4-497d-bc9a-c54b5aaa4400
	PAT string `json:"pat_id"`
}

func NewUserRouteRegisterer(children ...*router.Registerer) *router.Registerer {
	return &router.Registerer{
		GetRoutes: GetUserRoutes,
		Children:  children,
	}
}

func GetUserRoutes(
	r chi.Router,
	config *server.Config,
	basePath *endpoint.Path,
	factory endpoint.APIEndpointFactory,
	children ...*router.Registerer,
) []*router.Route {
	routes := make([]*router.Route, 0)

	if config.AuthConfig.BasicAuthEnabled {
		// POST /api/v1/users -> users.NewCreateUserHandler
		// swagger:operation POST /api/v1/users createUser
		//
		// ### Description
		//
		// Creates a new user via email and password-based authentication. This endpoint is only registered if the
		// environment variable `BASIC_AUTH_ENABLED` is set.
		//
		// ---
		// produces:
		// - application/json
		// summary: Create a new user
		// tags:
		// - Users
		// parameters:
		//   - in: body
		//     name: CreateUserRequest
		//     description: The user to create
		//     schema:
		//       $ref: '#/definitions/CreateUserRequest'
		// responses:
		//   '201':
		//     description: Successfully created the user
		//     schema:
		//       $ref: '#/definitions/CreateUserResponse'
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
		createUserEndpoint := factory.NewAPIEndpoint(
			&endpoint.EndpointMetadata{
				Verb:   types.APIVerbCreate,
				Method: types.HTTPVerbPost,
				Path: &endpoint.Path{
					Parent:       basePath,
					RelativePath: "/users",
				},
				Scopes: []types.PermissionScope{
					types.NoUserScope,
				},
			},
		)

		createUserHandler := users.NewUserCreateHandler(
			config,
			factory.GetDecoderValidator(),
			factory.GetResultWriter(),
		)

		routes = append(routes, &router.Route{
			Endpoint: createUserEndpoint,
			Handler:  createUserHandler,
			Router:   r,
		})

		// POST /api/v1/users/login -> users.NewUserLoginHandler
		// swagger:operation POST /api/v1/users/login loginUser
		//
		// ### Description
		//
		// Logs a user in via email and password-based authentication. This endpoint is only registered if the
		// environment variable `BASIC_AUTH_ENABLED` is set.
		//
		// ---
		// produces:
		// - application/json
		// summary: Login user
		// tags:
		// - Users
		// parameters:
		//   - in: body
		//     name: LoginUserRequest
		//     description: The credentials for basic login
		//     schema:
		//       $ref: '#/definitions/LoginUserRequest'
		// responses:
		//   '200':
		//     description: Successfully logged in
		//     schema:
		//       $ref: '#/definitions/LoginUserResponse'
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
		loginUserEndpoint := factory.NewAPIEndpoint(
			&endpoint.EndpointMetadata{
				Verb:   types.APIVerbCreate,
				Method: types.HTTPVerbPost,
				Path: &endpoint.Path{
					Parent:       basePath,
					RelativePath: "/users/login",
				},
				Scopes: []types.PermissionScope{
					types.NoUserScope,
				},
			},
		)

		loginUserHandler := users.NewUserLoginHandler(
			config,
			factory.GetDecoderValidator(),
			factory.GetResultWriter(),
		)

		routes = append(routes, &router.Route{
			Endpoint: loginUserEndpoint,
			Handler:  loginUserHandler,
			Router:   r,
		})
	}

	// POST /api/v1/users/current/reset_password_manual -> users.NewResetPasswordManualHandler
	// swagger:operation POST /api/v1/users/current/reset_password_manual resetPasswordManual
	//
	// ### Description
	//
	// Resets a password for a user using the old password as validation.
	//
	// ---
	// produces:
	// - application/json
	// summary: Reset password (manual)
	// tags:
	// - Users
	// parameters:
	//   - in: body
	//     name: ResetPasswordManualRequest
	//     description: The old password and new password
	//     schema:
	//       $ref: '#/definitions/ResetPasswordManualRequest'
	// responses:
	//   '200':
	//     description: Successfully reset password
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
	resetPasswordManualEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current/reset_password_manual",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	resetPasswordManualHandler := users.NewResetPasswordManualHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: resetPasswordManualEndpoint,
		Handler:  resetPasswordManualHandler,
		Router:   r,
	})

	// POST /api/v1/users/reset_password_email -> users.NewResetPasswordEmailHandler
	// swagger:operation POST /api/v1/users/reset_password_email resetPasswordEmail
	//
	// ### Description
	//
	// Resets a password for a user by sending them a verification email.
	//
	// ---
	// produces:
	// - application/json
	// summary: Reset password (email)
	// tags:
	// - Users
	// parameters:
	//   - in: body
	//     name: ResetPasswordEmailRequest
	//     description: Request for resetting a password over email
	//     schema:
	//       $ref: '#/definitions/ResetPasswordEmailRequest'
	// responses:
	//   '200':
	//     description: Successfully triggered password reset
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
	resetPasswordEmailEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/reset_password_email",
			},
			Scopes: []types.PermissionScope{
				types.NoUserScope,
			},
		},
	)

	resetPasswordEmailHandler := users.NewResetPasswordEmailHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: resetPasswordEmailEndpoint,
		Handler:  resetPasswordEmailHandler,
		Router:   r,
	})

	// POST /api/v1/users/reset_password_email/verify -> users.NewResetPasswordEmailVerifyHandler
	// swagger:operation POST /api/v1/users/reset_password_email/verify resetPasswordEmailVerify
	//
	// ### Description
	//
	// Verifies that the token id and token are valid for a given reset password request.
	//
	// ---
	// produces:
	// - application/json
	// summary: Verify password reset data
	// tags:
	// - Users
	// parameters:
	//   - in: body
	//     name: ResetPasswordEmailVerifyTokenRequest
	//     description: Token verification data
	//     schema:
	//       $ref: '#/definitions/ResetPasswordEmailVerifyTokenRequest'
	// responses:
	//   '200':
	//     description: Token is valid
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
	resetPasswordEmailVerifyEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/reset_password_email/verify",
			},
			Scopes: []types.PermissionScope{
				types.NoUserScope,
			},
		},
	)

	resetPasswordEmailVerifyHandler := users.NewResetPasswordEmailVerifyHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: resetPasswordEmailVerifyEndpoint,
		Handler:  resetPasswordEmailVerifyHandler,
		Router:   r,
	})

	// POST /api/v1/users/reset_password_email/finalize -> users.NewResetPasswordEmailFinalizeHandler
	// swagger:operation POST /api/v1/users/reset_password_email/finalize resetPasswordEmailFinalize
	//
	// ### Description
	//
	// Resets a user's password given a token-based reset password mechanism.
	//
	// ---
	// produces:
	// - application/json
	// summary: Reset password
	// tags:
	// - Users
	// parameters:
	//   - in: body
	//     name: ResetPasswordEmailFinalizeRequest
	//     description: Reset password data
	//     schema:
	//       $ref: '#/definitions/ResetPasswordEmailFinalizeRequest'
	// responses:
	//   '200':
	//     description: Password reset successfully
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
	resetPasswordEmailFinalizeEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/reset_password_email/finalize",
			},
			Scopes: []types.PermissionScope{
				types.NoUserScope,
			},
		},
	)

	resetPasswordEmailFinalizeHandler := users.NewResetPasswordEmailFinalizeHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: resetPasswordEmailFinalizeEndpoint,
		Handler:  resetPasswordEmailFinalizeHandler,
		Router:   r,
	})

	// POST /api/v1/users/current/verify_email/resend -> users.NewResendVerifyEmailHandler
	// swagger:operation POST /api/v1/users/current/verify_email/resend resendVerificationEmail
	//
	// ### Description
	//
	// Resends a verification email for the user.
	//
	// ---
	// produces:
	// - application/json
	// summary: Resend verification email.
	// tags:
	// - Users
	// responses:
	//   '200':
	//     description: Verification email resent successfully.
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
	resendVerificationEmailEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current/verify_email/resend",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
			AllowUnverifiedEmails: true,
		},
	)

	resendVerificationEmailHandler := users.NewResendVerifyEmailHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: resendVerificationEmailEndpoint,
		Handler:  resendVerificationEmailHandler,
		Router:   r,
	})

	// POST /api/v1/users/current/verify_email/finalize -> users.NewVerifyEmailHandler
	// swagger:operation POST /api/v1/users/current/verify_email/finalize verifyEmail
	//
	// ### Description
	//
	// Verifies a user's email via a token-based mechanism.
	//
	// ---
	// produces:
	// - application/json
	// summary: Verify email
	// tags:
	// - Users
	// parameters:
	//   - in: body
	//     name: VerifyEmailRequest
	//     description: Reset password data
	//     schema:
	//       $ref: '#/definitions/VerifyEmailRequest'
	// responses:
	//   '200':
	//     description: Email verified successfully
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
	verifyEmailEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current/verify_email/finalize",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
			AllowUnverifiedEmails: true,
		},
	)

	verifyEmailHandler := users.NewVerifyEmailHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: verifyEmailEndpoint,
		Handler:  verifyEmailHandler,
		Router:   r,
	})

	// POST /api/v1/users/logout -> users.NewUserLogoutHandler
	// swagger:operation POST /api/v1/users/logout logoutUser
	//
	// ### Description
	//
	// Logs a user out. This endpoint only performs an action if it's called with cookie-based authentication.
	//
	// ---
	// produces:
	// - application/json
	// summary: Logout user
	// tags:
	// - Users
	// responses:
	//   '200':
	//     description: Successfully logged out
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
	logoutUserEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/logout",
			},
			AllowUnverifiedEmails: true,
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	logoutUserHandler := users.NewUserLogoutHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: logoutUserEndpoint,
		Handler:  logoutUserHandler,
		Router:   r,
	})

	// GET /api/v1/users/current -> users.UserGetCurrentHandler
	// swagger:operation GET /api/v1/users/current getCurrentUser
	//
	// ### Description
	//
	// Retrieves the current user object based on the data passed in auth.
	//
	// ---
	// produces:
	// - application/json
	// summary: Retrieve the current user.
	// tags:
	// - Users
	// responses:
	//   '200':
	//     description: Successfully got the user
	//     schema:
	//       $ref: '#/definitions/GetUserResponse'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	getUserCurrentEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current",
			},
			AllowUnverifiedEmails: true,
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	getUserCurrentHandler := users.NewUserGetCurrentHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getUserCurrentEndpoint,
		Handler:  getUserCurrentHandler,
		Router:   r,
	})

	// POST /api/v1/users/current -> users.UserUpdateCurrentHandler
	// swagger:operation POST /api/v1/users/current updateCurrentUser
	//
	// ### Description
	//
	// Updates the current user.
	//
	// ---
	// produces:
	// - application/json
	// summary: Update the current user.
	// tags:
	// - Users
	// parameters:
	//   - in: body
	//     name: UpdateUserRequest
	//     description: The user parameters to update
	//     schema:
	//       $ref: '#/definitions/UpdateUserRequest'
	// responses:
	//   '200':
	//     description: Successfully updated the user.
	//     schema:
	//       $ref: '#/definitions/UpdateUserResponse'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	updateUserCurrentEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	updateUserCurrentHandler := users.NewUserUpdateCurrentHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: updateUserCurrentEndpoint,
		Handler:  updateUserCurrentHandler,
		Router:   r,
	})

	// DELETE /api/v1/users/current -> users.UserDeleteCurrentHandler
	// swagger:operation DELETE /api/v1/users/current deleteCurrentUser
	//
	// ### Description
	//
	// Deletes the current user.
	//
	// ---
	// produces:
	// - application/json
	// summary: Delete the current user.
	// tags:
	// - Users
	// responses:
	//   '202':
	//     description: Successfully triggered user deletion
	//     schema:
	//       $ref: '#/definitions/EmptyResponse'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	deleteUserCurrentEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbDelete,
			Method: types.HTTPVerbDelete,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current",
			},
			AllowUnverifiedEmails: true,
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	deleteUserCurrentHandler := users.NewUserDeleteCurrentHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: deleteUserCurrentEndpoint,
		Handler:  deleteUserCurrentHandler,
		Router:   r,
	})

	// POST /api/v1/users/current/settings/pats -> pats.NewPATCreateHandler
	// swagger:operation POST /api/v1/users/current/settings/pats createPersonalAccessToken
	//
	// ### Description
	//
	// Creates a new personal access token for a user.
	//
	// ---
	// produces:
	// - application/json
	// summary: Create a new personal access token
	// tags:
	// - Users
	// parameters:
	//   - in: body
	//     name: CreatePATRequest
	//     description: The personal access token to create
	//     schema:
	//       $ref: '#/definitions/CreatePATRequest'
	// responses:
	//   '201':
	//     description: Successfully created the personal access token
	//     schema:
	//       $ref: '#/definitions/CreatePATResponse'
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
	createPATEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbCreate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current/settings/pats",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	createPATHandler := pats.NewPATCreateHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: createPATEndpoint,
		Handler:  createPATHandler,
		Router:   r,
	})

	// GET /api/v1/users/current/settings/pats/{pat_id} -> pats.NewPATGetHandler
	// swagger:operation GET /api/v1/users/current/settings/pats/{pat_id} getPersonalAccessToken
	//
	// ### Description
	//
	// Gets a personal access token for a user, specified by the path param `pat_id`.
	//
	// ---
	// produces:
	// - application/json
	// summary: Get a personal access token
	// tags:
	// - Users
	// parameters:
	//   - name: pat_id
	// responses:
	//   '200':
	//     description: Successfully got the personal access token
	//     schema:
	//       $ref: '#/definitions/GetPATResponse'
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
	getPATEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbGet,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/users/current/settings/pats/{%s}", string(types.PersonalAccessTokenURLParam)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	getPATHandler := pats.NewPATGetHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: getPATEndpoint,
		Handler:  getPATHandler,
		Router:   r,
	})

	// GET /api/v1/users/current/settings/pats -> pats.NewPATListHandler
	// swagger:operation GET /api/v1/users/current/settings/pats listPersonalAccessTokens
	//
	// ### Description
	//
	// Lists personal access token for a user.
	//
	// ---
	// produces:
	// - application/json
	// summary: List personal access tokens.
	// tags:
	// - Users
	// parameters:
	//   - name: page
	// responses:
	//   '200':
	//     description: Successfully listed personal access tokens
	//     schema:
	//       $ref: '#/definitions/ListPATsResponse'
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
	listPATsEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbList,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current/settings/pats",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	listPATsHandler := pats.NewPATListHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: listPATsEndpoint,
		Handler:  listPATsHandler,
		Router:   r,
	})

	// POST /api/v1/users/current/settings/pats/{pat_id}/revoke -> pats.NewPATRevokeHandler
	// swagger:operation POST /api/v1/users/current/settings/pats/{pat_id}/revoke revokePersonalAccessToken
	//
	// ### Description
	//
	// Revokes the personal access token for the user
	//
	// ---
	// produces:
	// - application/json
	// summary: Revoke the personal access token.
	// tags:
	// - Users
	// parameters:
	//   - name: pat_id
	// responses:
	//   '200':
	//     description: Successfully revoked the personal access token
	//     schema:
	//       $ref: '#/definitions/RevokePATResponseExample'
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
	revokePATEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbUpdate,
			Method: types.HTTPVerbPost,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/users/current/settings/pats/{%s}/revoke", string(types.PersonalAccessTokenURLParam)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	revokePATHandler := pats.NewPATRevokeHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: revokePATEndpoint,
		Handler:  revokePATHandler,
		Router:   r,
	})

	// DELETE /api/v1/users/current/settings/pats/{pat_id} -> pats.NewPATDeleteHandler
	// swagger:operation DELETE /api/v1/users/current/settings/pats/{pat_id} deletePersonalAccessToken
	//
	// ### Description
	//
	// Deletes the personal access token for the user
	//
	// ---
	// produces:
	// - application/json
	// summary: Delete the personal access token.
	// tags:
	// - Users
	// parameters:
	//   - name: pat_id
	// responses:
	//   '200':
	//     description: Successfully deleted the personal access token
	//     schema:
	//       $ref: '#/definitions/DeletePATResponse'
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
	deletePATEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbDelete,
			Method: types.HTTPVerbDelete,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: fmt.Sprintf("/users/current/settings/pats/{%s}", string(types.PersonalAccessTokenURLParam)),
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	deletePATHandler := pats.NewPATDeleteHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: deletePATEndpoint,
		Handler:  deletePATHandler,
		Router:   r,
	})

	// GET /api/v1/users/current/organizations -> users.NewUserOrgListHandler
	// swagger:operation GET /api/v1/users/current/organizations listUserOrganizations
	//
	// ### Description
	//
	// Lists organizations for a user.
	//
	// ---
	// produces:
	// - application/json
	// summary: List user organizations
	// tags:
	// - Users
	// parameters:
	//   - name: page
	// responses:
	//   '200':
	//     description: Successfully listed organizations
	//     schema:
	//       $ref: '#/definitions/ListUserOrgsResponse'
	//   '400':
	//     description: A malformed or bad request
	//     schema:
	//       $ref: '#/definitions/APIErrorBadRequestExample'
	//   '403':
	//     description: Forbidden
	//     schema:
	//       $ref: '#/definitions/APIErrorForbiddenExample'
	listOrgsEndpoint := factory.NewAPIEndpoint(
		&endpoint.EndpointMetadata{
			Verb:   types.APIVerbList,
			Method: types.HTTPVerbGet,
			Path: &endpoint.Path{
				Parent:       basePath,
				RelativePath: "/users/current/organizations",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	listUserOrgsHandler := users.NewUserOrgListHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: listOrgsEndpoint,
		Handler:  listUserOrgsHandler,
		Router:   r,
	})

	// GET /api/v1/users/current/teams -> users.NewUserTeamListHandler
	// swagger:operation GET /api/v1/users/current/teams listUserTeams
	//
	// ### Description
	//
	// Lists teams for a user, optionally filtered by organization id.
	//
	// ---
	// produces:
	// - application/json
	// summary: List user teams
	// tags:
	// - Users
	// parameters:
	//   - name: page
	//   - name: organization_id
	// responses:
	//   '200':
	//     description: Successfully listed teams
	//     schema:
	//       $ref: '#/definitions/ListUserTeamsResponse'
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
				RelativePath: "/users/current/teams",
			},
			Scopes: []types.PermissionScope{
				types.UserScope,
			},
		},
	)

	listUserTeamsHandler := users.NewUserTeamListHandler(
		config,
		factory.GetDecoderValidator(),
		factory.GetResultWriter(),
	)

	routes = append(routes, &router.Route{
		Endpoint: listTeamsEndpoint,
		Handler:  listUserTeamsHandler,
		Router:   r,
	})

	return routes
}
