# {{classname}}

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateMonitor**](MonitorsApi.md#CreateMonitor) | **Post** /api/v1/teams/{team_id}/monitors | Create Monitor
[**DeleteMonitor**](MonitorsApi.md#DeleteMonitor) | **Delete** /api/v1/teams/{team_id}/monitors/{monitor_id} | Delete Monitor
[**GetMonitor**](MonitorsApi.md#GetMonitor) | **Get** /api/v1/teams/{team_id}/monitors/{monitor_id} | Get Monitor
[**ListMonitorResults**](MonitorsApi.md#ListMonitorResults) | **Get** /api/v1/teams/{team_id}/monitor_results | List Monitor Results
[**ListMonitors**](MonitorsApi.md#ListMonitors) | **Get** /api/v1/teams/{team_id}/monitors | List Monitors
[**UpdateMonitor**](MonitorsApi.md#UpdateMonitor) | **Post** /api/v1/teams/{team_id}/monitors/{monitor_id} | Update Monitor

# **CreateMonitor**
> CreateMonitorResponse CreateMonitor(ctx, teamId, optional)
Create Monitor

Creates a new monitor.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
 **optional** | ***MonitorsApiCreateMonitorOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a MonitorsApiCreateMonitorOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of CreateMonitorRequest**](CreateMonitorRequest.md)| The monitor to create | 

### Return type

[**CreateMonitorResponse**](CreateMonitorResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteMonitor**
> DeleteMonitorResponse DeleteMonitor(ctx, teamId, monitorId)
Delete Monitor

Deletes a monitor.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
  **monitorId** | **string**| The monitor id | 

### Return type

[**DeleteMonitorResponse**](DeleteMonitorResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMonitor**
> GetMonitorResponse GetMonitor(ctx, teamId, monitorId)
Get Monitor

Gets a monitor by id.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
  **monitorId** | **string**| The monitor id | 

### Return type

[**GetMonitorResponse**](GetMonitorResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListMonitorResults**
> ListMonitorResultsResponse ListMonitorResults(ctx, optional)
List Monitor Results

Lists monitor results for a given team, optionally filtered by module or monitor id.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***MonitorsApiListMonitorResultsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a MonitorsApiListMonitorResultsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **optional.Int64**| The page to query for | 
 **moduleMonitorId** | **optional.String**| The monitor id to filter by | 
 **moduleId** | **optional.String**| The module id to filter by | 

### Return type

[**ListMonitorResultsResponse**](ListMonitorResultsResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListMonitors**
> ListMonitorsResponse ListMonitors(ctx, teamId, optional)
List Monitors

Lists monitors for a given team.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
 **optional** | ***MonitorsApiListMonitorsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a MonitorsApiListMonitorsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **optional.Int64**| The page to query for | 

### Return type

[**ListMonitorsResponse**](ListMonitorsResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateMonitor**
> UpdateMonitorResponse UpdateMonitor(ctx, teamId, monitorId, optional)
Update Monitor

Updates a monitor.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
  **monitorId** | **string**| The monitor id | 
 **optional** | ***MonitorsApiUpdateMonitorOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a MonitorsApiUpdateMonitorOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **body** | [**optional.Interface of CreateMonitorRequest**](CreateMonitorRequest.md)| The monitor to update | 

### Return type

[**UpdateMonitorResponse**](UpdateMonitorResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

