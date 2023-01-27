# {{classname}}

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddTeamMember**](TeamsApi.md#AddTeamMember) | **Post** /api/v1/teams/{team_id}/members | Add team member
[**CreateTeam**](TeamsApi.md#CreateTeam) | **Post** /api/v1/organizations/{org_id}/teams | Create a new team.
[**DeleteTeam**](TeamsApi.md#DeleteTeam) | **Delete** /api/v1/teams/{team_id} | Delete team.
[**DeleteTeamMember**](TeamsApi.md#DeleteTeamMember) | **Delete** /api/v1/teams/{team_id}/members/{team_member_id} | Delete team member
[**ListTeamMembers**](TeamsApi.md#ListTeamMembers) | **Get** /api/v1/teams/{team_id}/members | List team members
[**ListTeams**](TeamsApi.md#ListTeams) | **Get** /api/v1/organizations/{org_id}/teams | List teams.
[**UpdateTeam**](TeamsApi.md#UpdateTeam) | **Post** /api/v1/teams/{team_id} | Update team

# **AddTeamMember**
> TeamAddMemberResponse AddTeamMember(ctx, teamId, optional)
Add team member

Add a team member from the organization members to the team.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
 **optional** | ***TeamsApiAddTeamMemberOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TeamsApiAddTeamMemberOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of AddTeamMemberRequest**](AddTeamMemberRequest.md)| The team member to add | 

### Return type

[**TeamAddMemberResponse**](TeamAddMemberResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateTeam**
> CreateTeamResponse CreateTeam(ctx, orgId, optional)
Create a new team.

Creates a new team, with the authenticated user set as a team admin.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **orgId** | **string**| The org id | 
 **optional** | ***TeamsApiCreateTeamOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TeamsApiCreateTeamOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of CreateTeamRequest**](CreateTeamRequest.md)| The team to create | 

### Return type

[**CreateTeamResponse**](CreateTeamResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteTeam**
> DeleteTeamResponse DeleteTeam(ctx, teamId)
Delete team.

Delete a team. This operation cannot be undone.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 

### Return type

[**DeleteTeamResponse**](DeleteTeamResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteTeamMember**
> EmptyResponse DeleteTeamMember(ctx, teamId)
Delete team member

Delete a team member.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team member id | 

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListTeamMembers**
> ListTeamMembersResponse ListTeamMembers(ctx, optional)
List team members

Lists team members for a team.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***TeamsApiListTeamMembersOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TeamsApiListTeamMembersOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **teamId** | **optional.Int64**| The page to query for | 

### Return type

[**ListTeamMembersResponse**](ListTeamMembersResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListTeams**
> ListTeamsResponse ListTeams(ctx, optional)
List teams.

Lists teams for an organization.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***TeamsApiListTeamsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TeamsApiListTeamsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **orgId** | **optional.Int64**| The page to query for | 

### Return type

[**ListTeamsResponse**](ListTeamsResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateTeam**
> TeamUpdateResponse UpdateTeam(ctx, teamId, optional)
Update team

Updates a team.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
 **optional** | ***TeamsApiUpdateTeamOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a TeamsApiUpdateTeamOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of CreateTeamRequest**](CreateTeamRequest.md)| The team parameters to update | 

### Return type

[**TeamUpdateResponse**](TeamUpdateResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

