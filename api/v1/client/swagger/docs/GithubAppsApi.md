# {{classname}}

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**FinishGithubAppOAuth**](GithubAppsApi.md#FinishGithubAppOAuth) | **Get** /api/v1/oauth/github_app/callback | Start Github App OAuth
[**GithubAppWebhook**](GithubAppsApi.md#GithubAppWebhook) | **Post** /api/v1/webhooks/github_app | Github App Webhook
[**InstallGithubApp**](GithubAppsApi.md#InstallGithubApp) | **Get** /api/v1/github_app/install | Install Github App
[**ListGithubRepoBranches**](GithubAppsApi.md#ListGithubRepoBranches) | **Get** /api/v1/github_app/{github_app_installation_id}/repos/{github_repo_owner}/{github_repo_name}/branches | List Github Branches
[**ListGithubRepos**](GithubAppsApi.md#ListGithubRepos) | **Get** /api/v1/github_app/{github_app_installation_id}/repos | List Github Repos
[**StartGithubAppOAuth**](GithubAppsApi.md#StartGithubAppOAuth) | **Get** /api/v1/oauth/github_app | Start Github App OAuth

# **FinishGithubAppOAuth**
> FinishGithubAppOAuth(ctx, )
Start Github App OAuth

Finishes the OAuth flow to authenticate with a Github App.

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GithubAppWebhook**
> GithubAppWebhook(ctx, )
Github App Webhook

Implements a Github App webhook.

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **InstallGithubApp**
> InstallGithubApp(ctx, )
Install Github App

Redirects the user to Github to install the Github App.

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListGithubRepoBranches**
> []GithubBranch ListGithubRepoBranches(ctx, )
List Github Branches

Lists the Github repo branches.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]GithubBranch**](array.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListGithubRepos**
> []GithubRepo ListGithubRepos(ctx, githubAppInstallationId)
List Github Repos

Lists the Github repos that the github app installation has access to.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **githubAppInstallationId** | **string**| The github app installation id | 

### Return type

[**[]GithubRepo**](array.md)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **StartGithubAppOAuth**
> StartGithubAppOAuth(ctx, )
Start Github App OAuth

Starts the OAuth flow to authenticate with a Github App.

### Required Parameters
This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

