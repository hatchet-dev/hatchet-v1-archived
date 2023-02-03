# {{classname}}

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateModule**](ModulesApi.md#CreateModule) | **Post** /api/v1/teams/{team_id}/modules | Create Module
[**CreateModuleRun**](ModulesApi.md#CreateModuleRun) | **Post** /api/v1/teams/{team_id}/modules/{module_id}/runs | Create Module Run
[**CreateTerraformPlan**](ModulesApi.md#CreateTerraformPlan) | **Post** /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/plan | Create Terraform plan
[**FinalizeModuleRun**](ModulesApi.md#FinalizeModuleRun) | **Post** /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/finalize | Finalize module run
[**GetModuleTarballURL**](ModulesApi.md#GetModuleTarballURL) | **Get** /api/v1/teams/{team_id}/modules/{module_id}/tarball_url | Get Module Tarball URL
[**GetTerraformState**](ModulesApi.md#GetTerraformState) | **Post** /api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate | Create or Update Terraform State
[**ListModuleRuns**](ModulesApi.md#ListModuleRuns) | **Get** /api/v1/teams/{team_id}/modules/{module_id}/runs | List Module Runs
[**ListModules**](ModulesApi.md#ListModules) | **Get** /api/v1/teams/{team_id}/modules | List Modules

# **CreateModule**
> CreateModuleResponse CreateModule(ctx, teamId, optional)
Create Module

Creates a new module.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
 **optional** | ***ModulesApiCreateModuleOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ModulesApiCreateModuleOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**optional.Interface of CreateModuleRequest**](CreateModuleRequest.md)| The module to create | 

### Return type

[**CreateModuleResponse**](CreateModuleResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateModuleRun**
> CreateModuleRun(ctx, teamId, moduleId)
Create Module Run

Creates a new module run.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
  **moduleId** | **string**| The module id | 

### Return type

 (empty response body)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateTerraformPlan**
> CreateTerraformPlan(ctx, body, teamId, moduleId, moduleRunId)
Create Terraform plan

Creates a `POST` request for a Terraform plan. **Should only be called by Terraform in automation.**

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CreateTerraformPlanRequest**](CreateTerraformPlanRequest.md)| The terraform plan contents | 
  **teamId** | **string**| The team id | 
  **moduleId** | **string**| The module id | 
  **moduleRunId** | **string**| The module run id | 

### Return type

 (empty response body)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **FinalizeModuleRun**
> FinalizeModuleRunResponse FinalizeModuleRun(ctx, body, teamId, moduleId, moduleRunId)
Finalize module run

Updates a module run with a finalized status.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**FinalizeModuleRunRequest**](FinalizeModuleRunRequest.md)| The module run status to update | 
  **teamId** | **string**| The team id | 
  **moduleId** | **string**| The module id | 
  **moduleRunId** | **string**| The module run id | 

### Return type

[**FinalizeModuleRunResponse**](FinalizeModuleRunResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetModuleTarballURL**
> GetModuleTarballUrlResponse GetModuleTarballURL(ctx, teamId, moduleId, optional)
Get Module Tarball URL

Gets the Github tarball URL for the module.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
  **moduleId** | **string**| The module id | 
 **optional** | ***ModulesApiGetModuleTarballURLOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ModulesApiGetModuleTarballURLOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **githubSha** | **optional.String**| the SHA to get the tarball from | 

### Return type

[**GetModuleTarballUrlResponse**](GetModuleTarballURLResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetTerraformState**
> GetTerraformState(ctx, teamId, moduleId, moduleRunId)
Create or Update Terraform State

Creates a `GET` request for Terraform state. **Should only be called by Terraform in automation.**

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
  **moduleId** | **string**| The module id | 
  **moduleRunId** | **string**| The module run id | 

### Return type

 (empty response body)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListModuleRuns**
> ListModuleRunsResponse ListModuleRuns(ctx, optional)
List Module Runs

Lists module runs for a given module id.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ModulesApiListModuleRunsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ModulesApiListModuleRunsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **teamId** | **optional.Int64**| The page to query for | 
 **moduleId** | **optional.String**| the status of the module run | 

### Return type

[**ListModuleRunsResponse**](ListModuleRunsResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListModules**
> ListModulesResponse ListModules(ctx, teamId, optional)
List Modules

Lists modules for a given team.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **teamId** | **string**| The team id | 
 **optional** | ***ModulesApiListModulesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ModulesApiListModulesOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **page** | **optional.Int64**| The page to query for | 

### Return type

[**ListModulesResponse**](ListModulesResponse.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

