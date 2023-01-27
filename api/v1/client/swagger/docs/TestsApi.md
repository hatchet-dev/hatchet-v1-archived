# {{classname}}

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetModuleTarballURL**](TestsApi.md#GetModuleTarballURL) | **Get** /api/v1/teams/{team_id}/modules/{module_id}/tarball_url | Get Module Tarball URL

# **GetModuleTarballURL**
> GetModuleTarballUrlResponse GetModuleTarballURL(ctx, teamId, moduleId)
Get Module Tarball URL

Gets the Github tarball URL for the module.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
  **moduleId** | **string**| The module id | 

### Return type

[**GetModuleTarballUrlResponse**](GetModuleTarballURLResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

