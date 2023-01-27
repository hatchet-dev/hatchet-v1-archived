# {{classname}}

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AcceptOrgMemberInvite**](OrganizationsApi.md#AcceptOrgMemberInvite) | **Post** /api/v1/invites/{org_member_invite_id}/{org_member_invite_tok} | Accept an organization invite.
[**CreateOrgMemberInvite**](OrganizationsApi.md#CreateOrgMemberInvite) | **Post** /api/v1/organizations/{org_id}/members | Create a member invite
[**CreateOrganization**](OrganizationsApi.md#CreateOrganization) | **Post** /api/v1/organizations | Create a new organization
[**DeleteOrg**](OrganizationsApi.md#DeleteOrg) | **Delete** /api/v1/organizations/{org_id} | Delete organization.
[**DeleteOrgMember**](OrganizationsApi.md#DeleteOrgMember) | **Delete** /api/v1/organizations/{org_id}/members/{org_member_id} | Delete organization member.
[**GetOrgMember**](OrganizationsApi.md#GetOrgMember) | **Get** /api/v1/organizations/{org_id}/members/{org_member_id} | Get organization member.
[**GetOrganization**](OrganizationsApi.md#GetOrganization) | **Get** /api/v1/organizations/{org_id} | Get an organization
[**LeaveOrg**](OrganizationsApi.md#LeaveOrg) | **Post** /api/v1/organizations/{org_id}/leave | Leave an organization
[**ListOrgMembers**](OrganizationsApi.md#ListOrgMembers) | **Get** /api/v1/organizations/{org_id}/members | List organization members
[**UpdateOrgMemberPolicies**](OrganizationsApi.md#UpdateOrgMemberPolicies) | **Post** /api/v1/organizations/{org_id}/members/{org_member_id}/update_policies | Update organization member policies.
[**UpdateOrgOwner**](OrganizationsApi.md#UpdateOrgOwner) | **Post** /api/v1/organizations/{org_id}/change_owner | Update organization owner.
[**UpdateOrganization**](OrganizationsApi.md#UpdateOrganization) | **Post** /api/v1/organizations/{org_id} | Update an organization

# **AcceptOrgMemberInvite**
> EmptyResponse AcceptOrgMemberInvite(ctx, orgMemberInviteId, orgMemberInviteTok)
Accept an organization invite.

Accept an invite for an organization.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **orgMemberInviteId** | **string**| The member invite id | 
  **orgMemberInviteTok** | **string**| The member invite token (sensitive) | 

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateOrgMemberInvite**
> CreateOrgMemberInviteResponse CreateOrgMemberInvite(ctx, orgId, optional)
Create a member invite

Creates a new invite for an organization member.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **orgId** | **string**| The org id | 
 **optional** | ***OrganizationsApiCreateOrgMemberInviteOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a OrganizationsApiCreateOrgMemberInviteOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of CreateOrgMemberInviteRequest**](CreateOrgMemberInviteRequest.md)| The org member to create | 

### Return type

[**CreateOrgMemberInviteResponse**](CreateOrgMemberInviteResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateOrganization**
> CreateOrganizationResponse CreateOrganization(ctx, optional)
Create a new organization

Creates a new organization, with the authenticated user set as the organization owner.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***OrganizationsApiCreateOrganizationOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a OrganizationsApiCreateOrganizationOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**optional.Interface of CreateOrganizationRequest**](CreateOrganizationRequest.md)| The organization to create | 

### Return type

[**CreateOrganizationResponse**](CreateOrganizationResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteOrg**
> DeleteOrganizationResponse DeleteOrg(ctx, )
Delete organization.

Delete an organization. Only owners can delete organizations.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**DeleteOrganizationResponse**](DeleteOrganizationResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteOrgMember**
> EmptyResponse DeleteOrgMember(ctx, orgId)
Delete organization member.

Delete an organization member. Only admins can delete an organization member. Owners cannot be removed from the organization, the owner must be transferred before the organization owner can be removed.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **orgId** | **string**| The org id | 

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetOrgMember**
> GetOrgMemberResponse GetOrgMember(ctx, orgId, orgMemberId)
Get organization member.

Get organization member. Only admins and owner can read full member data.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **orgId** | **string**| The org id | 
  **orgMemberId** | **string**| The org member id | 

### Return type

[**GetOrgMemberResponse**](GetOrgMemberResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetOrganization**
> GetOrganizationResponse GetOrganization(ctx, orgId)
Get an organization

Retrieves an organization by the `org_id`.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **orgId** | **string**| The org id | 

### Return type

[**GetOrganizationResponse**](GetOrganizationResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **LeaveOrg**
> EmptyResponse LeaveOrg(ctx, orgId)
Leave an organization

Leave an organization. The currently authenticated user will leave this organization. Owners cannot leave an organization without changing the owner first.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **orgId** | **string**| The org id | 

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListOrgMembers**
> ListOrgMembersResponse ListOrgMembers(ctx, optional)
List organization members

Lists organization members.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***OrganizationsApiListOrgMembersOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a OrganizationsApiListOrgMembersOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **orgId** | **optional.Int64**| The page to query for | 

### Return type

[**ListOrgMembersResponse**](ListOrgMembersResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateOrgMemberPolicies**
> UpdateOrgMemberPoliciesResponse UpdateOrgMemberPolicies(ctx, orgId, orgMemberId)
Update organization member policies.

Update an organization member's policies.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **orgId** | **string**| The org id | 
  **orgMemberId** | **string**| The org member id | 

### Return type

[**UpdateOrgMemberPoliciesResponse**](UpdateOrgMemberPoliciesResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateOrgOwner**
> EmptyResponse UpdateOrgOwner(ctx, orgId)
Update organization owner.

Update organization owner. Only owners may update organization owners. The previous owner will become an admin (and can subsequently be removed from the organization, if required).

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **orgId** | **string**| The org id | 

### Return type

[**EmptyResponse**](EmptyResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateOrganization**
> UpdateOrgResponse UpdateOrganization(ctx, orgId, optional)
Update an organization

Updates organization metadata.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **orgId** | **string**| The org id | 
 **optional** | ***OrganizationsApiUpdateOrganizationOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a OrganizationsApiUpdateOrganizationOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of UpdateOrganizationRequest**](UpdateOrganizationRequest.md)| The values to update | 

### Return type

[**UpdateOrgResponse**](UpdateOrgResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

