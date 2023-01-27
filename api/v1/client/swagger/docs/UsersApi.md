# {{classname}}

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreatePersonalAccessToken**](UsersApi.md#CreatePersonalAccessToken) | **Post** /api/v1/users/current/settings/pats | Create a new personal access token
[**CreateUser**](UsersApi.md#CreateUser) | **Post** /api/v1/users | Create a new user
[**DeleteCurrentUser**](UsersApi.md#DeleteCurrentUser) | **Delete** /api/v1/users/current | Delete the current user.
[**DeletePersonalAccessToken**](UsersApi.md#DeletePersonalAccessToken) | **Delete** /api/v1/users/current/settings/pats/{pat_id} | Delete the personal access token.
[**GetCurrentUser**](UsersApi.md#GetCurrentUser) | **Get** /api/v1/users/current | Retrieve the current user.
[**GetPersonalAccessToken**](UsersApi.md#GetPersonalAccessToken) | **Get** /api/v1/users/current/settings/pats/{pat_id} | Get a personal access token
[**ListGithubAppInstallations**](UsersApi.md#ListGithubAppInstallations) | **Get** /api/v1/users/current/github_app/installations | List Github App installations
[**ListPersonalAccessTokens**](UsersApi.md#ListPersonalAccessTokens) | **Get** /api/v1/users/current/settings/pats | List personal access tokens.
[**ListUserOrganizations**](UsersApi.md#ListUserOrganizations) | **Get** /api/v1/users/current/organizations | List user organizations
[**ListUserTeams**](UsersApi.md#ListUserTeams) | **Get** /api/v1/users/current/teams | List user teams
[**LoginUser**](UsersApi.md#LoginUser) | **Post** /api/v1/users/login | Login user
[**LogoutUser**](UsersApi.md#LogoutUser) | **Post** /api/v1/users/logout | Logout user
[**ResendVerificationEmail**](UsersApi.md#ResendVerificationEmail) | **Post** /api/v1/users/current/verify_email/resend | Resend verification email.
[**ResetPasswordEmail**](UsersApi.md#ResetPasswordEmail) | **Post** /api/v1/users/reset_password_email | Reset password (email)
[**ResetPasswordEmailFinalize**](UsersApi.md#ResetPasswordEmailFinalize) | **Post** /api/v1/users/reset_password_email/finalize | Reset password
[**ResetPasswordEmailVerify**](UsersApi.md#ResetPasswordEmailVerify) | **Post** /api/v1/users/reset_password_email/verify | Verify password reset data
[**ResetPasswordManual**](UsersApi.md#ResetPasswordManual) | **Post** /api/v1/users/current/reset_password_manual | Reset password (manual)
[**RevokePersonalAccessToken**](UsersApi.md#RevokePersonalAccessToken) | **Post** /api/v1/users/current/settings/pats/{pat_id}/revoke | Revoke the personal access token.
[**UpdateCurrentUser**](UsersApi.md#UpdateCurrentUser) | **Post** /api/v1/users/current | Update the current user.
[**VerifyEmail**](UsersApi.md#VerifyEmail) | **Post** /api/v1/users/current/verify_email/finalize | Verify email

# **CreatePersonalAccessToken**
> CreatePatResponse CreatePersonalAccessToken(ctx, optional)
Create a new personal access token

Creates a new personal access token for a user.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiCreatePersonalAccessTokenOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiCreatePersonalAccessTokenOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of CreatePersonalAccessTokenRequest**](CreatePersonalAccessTokenRequest.md)| The personal access token to create | 

### Return type

[**CreatePatResponse**](CreatePATResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateUser**
> CreateUserResponse CreateUser(ctx, optional)
Create a new user

Creates a new user via email and password-based authentication. This endpoint is only registered if the environment variable `BASIC_AUTH_ENABLED` is set.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiCreateUserOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiCreateUserOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of CreateUserRequest**](CreateUserRequest.md)| The user to create | 

### Return type

[**CreateUserResponse**](CreateUserResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteCurrentUser**
> EmptyResponse DeleteCurrentUser(ctx, )
Delete the current user.

Deletes the current user.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeletePersonalAccessToken**
> DeletePatResponse DeletePersonalAccessToken(ctx, patId)
Delete the personal access token.

Deletes the personal access token for the user

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **patId** | **string**| The personal access token id | 

### Return type

[**DeletePatResponse**](DeletePATResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCurrentUser**
> GetUserResponse GetCurrentUser(ctx, )
Retrieve the current user.

Retrieves the current user object based on the data passed in auth.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**GetUserResponse**](GetUserResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetPersonalAccessToken**
> GetPatResponse GetPersonalAccessToken(ctx, patId)
Get a personal access token

Gets a personal access token for a user, specified by the path param `pat_id`.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **patId** | **string**| The personal access token id | 

### Return type

[**GetPatResponse**](GetPATResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListGithubAppInstallations**
> ListGithubAppInstallationsResponse ListGithubAppInstallations(ctx, optional)
List Github App installations

Lists the github app installations for the currently authenticated user.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiListGithubAppInstallationsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiListGithubAppInstallationsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **optional.Int64**| The page to query for | 

### Return type

[**ListGithubAppInstallationsResponse**](ListGithubAppInstallationsResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListPersonalAccessTokens**
> ListPatsResponse ListPersonalAccessTokens(ctx, optional)
List personal access tokens.

Lists personal access token for a user.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiListPersonalAccessTokensOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiListPersonalAccessTokensOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **optional.Int64**| The page to query for | 

### Return type

[**ListPatsResponse**](ListPATsResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListUserOrganizations**
> ListUserOrgsResponse ListUserOrganizations(ctx, optional)
List user organizations

Lists organizations for a user.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiListUserOrganizationsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiListUserOrganizationsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **optional.Int64**| The page to query for | 

### Return type

[**ListUserOrgsResponse**](ListUserOrgsResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListUserTeams**
> ListUserTeamsResponse ListUserTeams(ctx, optional)
List user teams

Lists teams for a user, optionally filtered by organization id.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiListUserTeamsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiListUserTeamsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **optional.Int64**| The page to query for | 
 **organizationId** | **optional.String**| the id of the organization to filter by (optional) | 

### Return type

[**ListUserTeamsResponse**](ListUserTeamsResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **LoginUser**
> LoginUserResponse LoginUser(ctx, optional)
Login user

Logs a user in via email and password-based authentication. This endpoint is only registered if the environment variable `BASIC_AUTH_ENABLED` is set.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiLoginUserOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiLoginUserOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of LoginUserRequest**](LoginUserRequest.md)| The credentials for basic login | 

### Return type

[**LoginUserResponse**](LoginUserResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **LogoutUser**
> EmptyResponse LogoutUser(ctx, )
Logout user

Logs a user out. This endpoint only performs an action if it's called with cookie-based authentication.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ResendVerificationEmail**
> EmptyResponse ResendVerificationEmail(ctx, )
Resend verification email.

Resends a verification email for the user.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ResetPasswordEmail**
> EmptyResponse ResetPasswordEmail(ctx, optional)
Reset password (email)

Resets a password for a user by sending them a verification email.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiResetPasswordEmailOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiResetPasswordEmailOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of ResetPasswordEmailRequest**](ResetPasswordEmailRequest.md)| Request for resetting a password over email | 

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ResetPasswordEmailFinalize**
> EmptyResponse ResetPasswordEmailFinalize(ctx, optional)
Reset password

Resets a user's password given a token-based reset password mechanism.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiResetPasswordEmailFinalizeOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiResetPasswordEmailFinalizeOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of ResetPasswordEmailFinalizeRequest**](ResetPasswordEmailFinalizeRequest.md)| Reset password data | 

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ResetPasswordEmailVerify**
> EmptyResponse ResetPasswordEmailVerify(ctx, optional)
Verify password reset data

Verifies that the token id and token are valid for a given reset password request.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiResetPasswordEmailVerifyOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiResetPasswordEmailVerifyOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of ResetPasswordEmailVerifyRequest**](ResetPasswordEmailVerifyRequest.md)| Token verification data | 

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ResetPasswordManual**
> EmptyResponse ResetPasswordManual(ctx, optional)
Reset password (manual)

Resets a password for a user using the old password as validation.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiResetPasswordManualOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiResetPasswordManualOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of ResetPasswordManualRequest**](ResetPasswordManualRequest.md)| The old password and new password | 

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RevokePersonalAccessToken**
> RevokePatResponseExample RevokePersonalAccessToken(ctx, patId)
Revoke the personal access token.

Revokes the personal access token for the user

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **patId** | **string**| The personal access token id | 

### Return type

[**RevokePatResponseExample**](RevokePATResponseExample.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateCurrentUser**
> UpdateUserResponse UpdateCurrentUser(ctx, optional)
Update the current user.

Updates the current user.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiUpdateCurrentUserOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiUpdateCurrentUserOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of UpdateOrganizationRequest**](UpdateOrganizationRequest.md)| The user parameters to update | 

### Return type

[**UpdateUserResponse**](UpdateUserResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **VerifyEmail**
> EmptyResponse VerifyEmail(ctx, optional)
Verify email

Verifies a user's email via a token-based mechanism.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***UsersApiVerifyEmailOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a UsersApiVerifyEmailOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of VerifyEmailRequest**](VerifyEmailRequest.md)| Reset password data | 

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

