/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

import {
  AddTeamMemberRequest,
  APIErrorBadRequestExample,
  APIErrorForbiddenExample,
  APIErrorNotSupportedExample,
  APIServerMetadata,
  CreateModuleRequest,
  CreateModuleResponse,
  CreateModuleRunRequest,
  CreateModuleRunResponse,
  CreateMonitorRequest,
  CreateMonitorResponse,
  CreateMonitorResultRequest,
  CreateOrganizationRequest,
  CreateOrganizationResponse,
  CreateOrgMemberInviteRequest,
  CreateOrgMemberInviteResponse,
  CreatePATResponse,
  CreatePersonalAccessTokenRequest,
  CreateTeamRequest,
  CreateTeamResponse,
  CreateTerraformPlanRequest,
  CreateUserRequest,
  CreateUserResponse,
  DeleteModuleResponse,
  DeleteMonitorResponse,
  DeleteOrganizationResponse,
  DeletePATResponse,
  DeleteTeamResponse,
  EmptyResponse,
  FinalizeModuleRunRequest,
  FinalizeModuleRunResponse,
  ForceUnlockModuleResponse,
  GetLogsResponse,
  GetModuleEnvVarsVersionResponse,
  GetModulePlanSummaryResponse,
  GetModuleResponse,
  GetModuleRunResponse,
  GetModuleRunTokenResponse,
  GetModuleTarballURLResponse,
  GetModuleValuesCurrentResponse,
  GetModuleValuesResponse,
  GetMonitorResponse,
  GetNotificationResponse,
  GetOrganizationResponse,
  GetOrgMemberResponse,
  GetPATResponse,
  GetUserResponse,
  ListGithubAppInstallationsResponse,
  ListGithubRepoBranchesResponse,
  ListGithubReposResponse,
  ListModuleRunsResponse,
  ListModulesResponse,
  ListMonitorResultsResponse,
  ListMonitorsResponse,
  ListNotificationsResponse,
  ListOrgMembersResponse,
  ListPATsResponse,
  ListTeamMembersResponse,
  ListTeamsResponse,
  ListUserOrgsResponse,
  ListUserTeamsResponse,
  LoginUserRequest,
  LoginUserResponse,
  ModuleRunKind,
  ModuleRunStatus,
  ResetPasswordEmailFinalizeRequest,
  ResetPasswordEmailRequest,
  ResetPasswordEmailVerifyRequest,
  ResetPasswordManualRequest,
  RevokePATResponseExample,
  TeamAddMemberResponse,
  TeamUpdateResponse,
  UpdateModuleRequest,
  UpdateModuleResponse,
  UpdateMonitorResponse,
  UpdateOrganizationRequest,
  UpdateOrgMemberPoliciesResponse,
  UpdateOrgResponse,
  UpdateUserResponse,
  VerifyEmailRequest,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Api<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
  /**
   * @description Lists the Github repos that the github app installation has access to.
   *
   * @tags Github Apps
   * @name ListGithubRepos
   * @summary List Github Repos
   * @request GET:/api/v1/github_app/{github_app_installation_id}/repos
   * @secure
   */
  listGithubRepos = (githubAppInstallationId: string, params: RequestParams = {}) =>
    this.request<ListGithubReposResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/github_app/${githubAppInstallationId}/repos`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Lists the Github repo branches.
   *
   * @tags Github Apps
   * @name ListGithubRepoBranches
   * @summary List Github Branches
   * @request GET:/api/v1/github_app/{github_app_installation_id}/repos/{github_repo_owner}/{github_repo_name}/branches
   * @secure
   */
  listGithubRepoBranches = (
    githubAppInstallationId: string,
    githubRepoOwner: string,
    githubRepoName: string,
    params: RequestParams = {},
  ) =>
    this.request<ListGithubRepoBranchesResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/github_app/${githubAppInstallationId}/repos/${githubRepoOwner}/${githubRepoName}/branches`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Redirects the user to Github to install the Github App.
   *
   * @tags Github Apps
   * @name InstallGithubApp
   * @summary Install Github App
   * @request GET:/api/v1/github_app/install
   * @secure
   */
  installGithubApp = (params: RequestParams = {}) =>
    this.request<any, void | APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>({
      path: `/api/v1/github_app/install`,
      method: "GET",
      secure: true,
      ...params,
    });
  /**
   * @description Accept an invite for an organization.
   *
   * @tags Organizations
   * @name AcceptOrgMemberInvite
   * @summary Accept an organization invite.
   * @request POST:/api/v1/invites/{org_member_invite_id}/{org_member_invite_tok}
   * @secure
   */
  acceptOrgMemberInvite = (orgMemberInviteId: string, orgMemberInviteTok: string, params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>({
      path: `/api/v1/invites/${orgMemberInviteId}/${orgMemberInviteTok}`,
      method: "POST",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Gets the metadata for the Hatchet instance.
   *
   * @tags Metadata
   * @name GetServerMetadata
   * @summary Get server metadata
   * @request GET:/api/v1/metadata
   * @secure
   */
  getServerMetadata = (params: RequestParams = {}) =>
    this.request<APIServerMetadata, any>({
      path: `/api/v1/metadata`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Starts the OAuth flow to authenticate with a Github App.
   *
   * @tags Github Apps
   * @name StartGithubAppOAuth
   * @summary Start Github App OAuth
   * @request GET:/api/v1/oauth/github_app
   * @secure
   */
  startGithubAppOAuth = (params: RequestParams = {}) =>
    this.request<any, void | APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>({
      path: `/api/v1/oauth/github_app`,
      method: "GET",
      secure: true,
      ...params,
    });
  /**
   * @description Finishes the OAuth flow to authenticate with a Github App.
   *
   * @tags Github Apps
   * @name FinishGithubAppOAuth
   * @summary Start Github App OAuth
   * @request GET:/api/v1/oauth/github_app/callback
   * @secure
   */
  finishGithubAppOAuth = (params: RequestParams = {}) =>
    this.request<any, void | APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>({
      path: `/api/v1/oauth/github_app/callback`,
      method: "GET",
      secure: true,
      ...params,
    });
  /**
   * @description Creates a new organization, with the authenticated user set as the organization owner.
   *
   * @tags Organizations
   * @name CreateOrganization
   * @summary Create a new organization
   * @request POST:/api/v1/organizations
   * @secure
   */
  createOrganization = (data?: CreateOrganizationRequest, params: RequestParams = {}) =>
    this.request<
      CreateOrganizationResponse,
      APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample
    >({
      path: `/api/v1/organizations`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Delete an organization. Only owners can delete organizations.
   *
   * @tags Organizations
   * @name DeleteOrg
   * @summary Delete organization.
   * @request DELETE:/api/v1/organizations/{org_id}
   * @secure
   */
  deleteOrg = (orgId: string, params: RequestParams = {}) =>
    this.request<DeleteOrganizationResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}`,
      method: "DELETE",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Retrieves an organization by the `org_id`.
   *
   * @tags Organizations
   * @name GetOrganization
   * @summary Get an organization
   * @request GET:/api/v1/organizations/{org_id}
   * @secure
   */
  getOrganization = (orgId: string, params: RequestParams = {}) =>
    this.request<
      GetOrganizationResponse,
      APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample
    >({
      path: `/api/v1/organizations/${orgId}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Updates organization metadata.
   *
   * @tags Organizations
   * @name UpdateOrganization
   * @summary Update an organization
   * @request POST:/api/v1/organizations/{org_id}
   * @secure
   */
  updateOrganization = (orgId: string, data?: UpdateOrganizationRequest, params: RequestParams = {}) =>
    this.request<UpdateOrgResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Update organization owner. Only owners may update organization owners. The previous owner will become an admin (and can subsequently be removed from the organization, if required).
   *
   * @tags Organizations
   * @name UpdateOrgOwner
   * @summary Update organization owner.
   * @request POST:/api/v1/organizations/{org_id}/change_owner
   * @secure
   */
  updateOrgOwner = (orgId: string, params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/change_owner`,
      method: "POST",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Leave an organization. The currently authenticated user will leave this organization. Owners cannot leave an organization without changing the owner first.
   *
   * @tags Organizations
   * @name LeaveOrg
   * @summary Leave an organization
   * @request POST:/api/v1/organizations/{org_id}/leave
   * @secure
   */
  leaveOrg = (orgId: string, params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/leave`,
      method: "POST",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Lists organization members.
   *
   * @tags Organizations
   * @name ListOrgMembers
   * @summary List organization members
   * @request GET:/api/v1/organizations/{org_id}/members
   * @secure
   */
  listOrgMembers = (
    orgId: string,
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      org_id?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<ListOrgMembersResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/members`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Creates a new invite for an organization member.
   *
   * @tags Organizations
   * @name CreateOrgMemberInvite
   * @summary Create a member invite
   * @request POST:/api/v1/organizations/{org_id}/members
   * @secure
   */
  createOrgMemberInvite = (orgId: string, data?: CreateOrgMemberInviteRequest, params: RequestParams = {}) =>
    this.request<
      CreateOrgMemberInviteResponse,
      APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample
    >({
      path: `/api/v1/organizations/${orgId}/members`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Delete an organization member. Only admins can delete an organization member. Owners cannot be removed from the organization, the owner must be transferred before the organization owner can be removed.
   *
   * @tags Organizations
   * @name DeleteOrgMember
   * @summary Delete organization member.
   * @request DELETE:/api/v1/organizations/{org_id}/members/{org_member_id}
   * @secure
   */
  deleteOrgMember = (orgId: string, orgMemberId: string, params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/members/${orgMemberId}`,
      method: "DELETE",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Get organization member. Only admins and owner can read full member data.
   *
   * @tags Organizations
   * @name GetOrgMember
   * @summary Get organization member.
   * @request GET:/api/v1/organizations/{org_id}/members/{org_member_id}
   * @secure
   */
  getOrgMember = (orgId: string, orgMemberId: string, params: RequestParams = {}) =>
    this.request<GetOrgMemberResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/members/${orgMemberId}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Update an organization member's policies.
   *
   * @tags Organizations
   * @name UpdateOrgMemberPolicies
   * @summary Update organization member policies.
   * @request POST:/api/v1/organizations/{org_id}/members/{org_member_id}/update_policies
   * @secure
   */
  updateOrgMemberPolicies = (orgId: string, orgMemberId: string, params: RequestParams = {}) =>
    this.request<UpdateOrgMemberPoliciesResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/members/${orgMemberId}/update_policies`,
      method: "POST",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Lists notifications for an organization, optionally filtered by team id.
   *
   * @tags Teams
   * @name ListNotifications
   * @summary List notifications
   * @request GET:/api/v1/organizations/{org_id}/notifications
   * @secure
   */
  listNotifications = (
    orgId: string,
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      org_id?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<ListNotificationsResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/notifications`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Lists teams for an organization.
   *
   * @tags Teams
   * @name ListTeams
   * @summary List teams.
   * @request GET:/api/v1/organizations/{org_id}/teams
   * @secure
   */
  listTeams = (
    orgId: string,
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      org_id?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<ListTeamsResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/teams`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Creates a new team, with the authenticated user set as a team admin.
   *
   * @tags Teams
   * @name CreateTeam
   * @summary Create a new team.
   * @request POST:/api/v1/organizations/{org_id}/teams
   * @secure
   */
  createTeam = (orgId: string, data?: CreateTeamRequest, params: RequestParams = {}) =>
    this.request<CreateTeamResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/teams`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Delete a team. This operation cannot be undone.
   *
   * @tags Teams
   * @name DeleteTeam
   * @summary Delete team.
   * @request DELETE:/api/v1/teams/{team_id}
   * @secure
   */
  deleteTeam = (teamId: string, params: RequestParams = {}) =>
    this.request<DeleteTeamResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}`,
      method: "DELETE",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Updates a team.
   *
   * @tags Teams
   * @name UpdateTeam
   * @summary Update team
   * @request POST:/api/v1/teams/{team_id}
   * @secure
   */
  updateTeam = (teamId: string, data?: CreateTeamRequest, params: RequestParams = {}) =>
    this.request<TeamUpdateResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Github incoming webhook handler.
   *
   * @tags Teams
   * @name GithubIncomingWebhook
   * @summary Github incoming webhook endpoint
   * @request POST:/api/v1/teams/{team_id}/github_incoming/{github_incoming_webhook_id}
   * @secure
   */
  githubIncomingWebhook = (teamId: string, githubIncomingWebhookId: string, params: RequestParams = {}) =>
    this.request<void, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/github_incoming/${githubIncomingWebhookId}`,
      method: "POST",
      secure: true,
      ...params,
    });
  /**
   * @description Lists team members for a team.
   *
   * @tags Teams
   * @name ListTeamMembers
   * @summary List team members
   * @request GET:/api/v1/teams/{team_id}/members
   * @secure
   */
  listTeamMembers = (
    teamId: string,
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      team_id?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<ListTeamMembersResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/members`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Add a team member from the organization members to the team.
   *
   * @tags Teams
   * @name AddTeamMember
   * @summary Add team member
   * @request POST:/api/v1/teams/{team_id}/members
   * @secure
   */
  addTeamMember = (teamId: string, data?: AddTeamMemberRequest, params: RequestParams = {}) =>
    this.request<TeamAddMemberResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/members`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Delete a team member.
   *
   * @tags Teams
   * @name DeleteTeamMember
   * @summary Delete team member
   * @request DELETE:/api/v1/teams/{team_id}/members/{team_member_id}
   * @secure
   */
  deleteTeamMember = (teamId: string, teamMemberId: string, params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/members/${teamMemberId}`,
      method: "DELETE",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Lists modules for a given team.
   *
   * @tags Modules
   * @name ListModules
   * @summary List Modules
   * @request GET:/api/v1/teams/{team_id}/modules
   * @secure
   */
  listModules = (
    teamId: string,
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      page?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<ListModulesResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Creates a new module.
   *
   * @tags Modules
   * @name CreateModule
   * @summary Create Module
   * @request POST:/api/v1/teams/{team_id}/modules
   * @secure
   */
  createModule = (teamId: string, data: CreateModuleRequest, params: RequestParams = {}) =>
    this.request<CreateModuleResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Triggers module deletion. First triggers destroy for the Terraform module.
   *
   * @tags Modules
   * @name DeleteModule
   * @summary Delete Module
   * @request DELETE:/api/v1/teams/{team_id}/modules/{module_id}
   * @secure
   */
  deleteModule = (teamId: string, moduleId: string, params: RequestParams = {}) =>
    this.request<DeleteModuleResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}`,
      method: "DELETE",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Gets a module by module id.
   *
   * @tags Modules
   * @name GetModule
   * @summary Get module
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}
   * @secure
   */
  getModule = (teamId: string, moduleId: string, params: RequestParams = {}) =>
    this.request<GetModuleResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Updates a module.
   *
   * @tags Modules
   * @name UpdateModule
   * @summary Update Module Run
   * @request POST:/api/v1/teams/{team_id}/modules/{module_id}
   * @secure
   */
  updateModule = (teamId: string, moduleId: string, data?: UpdateModuleRequest, params: RequestParams = {}) =>
    this.request<UpdateModuleResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Gets the module env vars version by id.
   *
   * @tags Modules
   * @name GetModuleEnvVars
   * @summary Get Module Env Vars
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}/env_vars/{module_env_vars_id}
   * @secure
   */
  getModuleEnvVars = (teamId: string, moduleId: string, moduleEnvVarsId: string, params: RequestParams = {}) =>
    this.request<GetModuleEnvVarsVersionResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/env_vars/${moduleEnvVarsId}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Force unlocks a module.
   *
   * @tags Modules
   * @name ForceUnlockModule
   * @summary Unlock Module
   * @request POST:/api/v1/teams/{team_id}/modules/{module_id}/force_unlock
   * @secure
   */
  forceUnlockModule = (teamId: string, moduleId: string, params: RequestParams = {}) =>
    this.request<ForceUnlockModuleResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/force_unlock`,
      method: "POST",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Lists module runs for a given module id.
   *
   * @tags Modules
   * @name ListModuleRuns
   * @summary List Module Runs
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}/runs
   * @secure
   */
  listModuleRuns = (
    teamId: string,
    moduleId: string,
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      page?: number;
      /** an optional list of statuses for the module run */
      status?: ModuleRunStatus[];
      /** an optional list of kinds for the module run */
      kind?: ModuleRunKind[];
    },
    params: RequestParams = {},
  ) =>
    this.request<ListModuleRunsResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Creates a new module run.
   *
   * @tags Modules
   * @name CreateModuleRun
   * @summary Create Module Run
   * @request POST:/api/v1/teams/{team_id}/modules/{module_id}/runs
   * @secure
   */
  createModuleRun = (teamId: string, moduleId: string, data: CreateModuleRunRequest, params: RequestParams = {}) =>
    this.request<CreateModuleRunResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Gets a module run.
   *
   * @tags Modules
   * @name GetModuleRun
   * @summary Get module run
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}
   * @secure
   */
  getModuleRun = (teamId: string, moduleId: string, moduleRunId: string, params: RequestParams = {}) =>
    this.request<GetModuleRunResponse, APIErrorBadRequestExample | APIErrorForbiddenExample | void>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs/${moduleRunId}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Updates a module run with a finalized status.
   *
   * @tags Modules
   * @name FinalizeModuleRun
   * @summary Finalize module run
   * @request POST:/api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/finalize
   * @secure
   */
  finalizeModuleRun = (
    teamId: string,
    moduleId: string,
    moduleRunId: string,
    data: FinalizeModuleRunRequest,
    params: RequestParams = {},
  ) =>
    this.request<FinalizeModuleRunResponse, APIErrorBadRequestExample | APIErrorForbiddenExample | void>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs/${moduleRunId}/finalize`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Gets a module run token for a local run.
   *
   * @tags Modules
   * @name GetModuleRunLocalToken
   * @summary Get Module Run Token (Local)
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/local_token
   * @secure
   */
  getModuleRunLocalToken = (teamId: string, moduleId: string, moduleRunId: string, params: RequestParams = {}) =>
    this.request<GetModuleRunTokenResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs/${moduleRunId}/local_token`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Gets the logs for a module run.
   *
   * @tags Modules
   * @name GetModuleRunLogs
   * @summary Get logs
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/logs
   * @secure
   */
  getModuleRunLogs = (teamId: string, moduleId: string, moduleRunId: string, params: RequestParams = {}) =>
    this.request<GetLogsResponse, APIErrorBadRequestExample | APIErrorForbiddenExample | void>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs/${moduleRunId}/logs`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Reports a monitor result.
   *
   * @tags Modules
   * @name CreateMonitorResult
   * @summary Create Monitor Result
   * @request POST:/api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/monitor_result
   * @secure
   */
  createMonitorResult = (
    teamId: string,
    moduleId: string,
    moduleRunId: string,
    data: CreateMonitorResultRequest,
    params: RequestParams = {},
  ) =>
    this.request<void, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs/${moduleRunId}/monitor_result`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Creates a `POST` request for a Terraform plan. **Should only be called by Terraform in automation.**
   *
   * @tags Modules
   * @name CreateTerraformPlan
   * @summary Create Terraform plan
   * @request POST:/api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/plan
   * @secure
   */
  createTerraformPlan = (
    teamId: string,
    moduleId: string,
    moduleRunId: string,
    data: CreateTerraformPlanRequest,
    params: RequestParams = {},
  ) =>
    this.request<void, APIErrorBadRequestExample | APIErrorForbiddenExample | void>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs/${moduleRunId}/plan`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      ...params,
    });
  /**
   * @description Gets the plan summary for a module run.
   *
   * @tags Modules
   * @name GetModuleRunPlanSummary
   * @summary Get plan summary
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/plan_summary
   * @secure
   */
  getModuleRunPlanSummary = (teamId: string, moduleId: string, moduleRunId: string, params: RequestParams = {}) =>
    this.request<GetModulePlanSummaryResponse, APIErrorBadRequestExample | APIErrorForbiddenExample | void>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs/${moduleRunId}/plan_summary`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Creates a `GET` request for Terraform state. **Should only be called by Terraform in automation.**
   *
   * @tags Modules
   * @name GetTerraformState
   * @summary Create or Update Terraform State
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate
   * @secure
   */
  getTerraformState = (teamId: string, moduleId: string, moduleRunId: string, params: RequestParams = {}) =>
    this.request<void, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs/${moduleRunId}/tfstate`,
      method: "GET",
      secure: true,
      ...params,
    });
  /**
   * @description Creates a `POST` request for Terraform state. **Should only be called by Terraform in automation.**
   *
   * @tags Modules
   * @name CreateTerraformState
   * @summary Create or Update Terraform State
   * @request POST:/api/v1/teams/{team_id}/modules/{module_id}/runs/{module_run_id}/tfstate
   * @secure
   */
  createTerraformState = (teamId: string, moduleId: string, moduleRunId: string, params: RequestParams = {}) =>
    this.request<void, APIErrorBadRequestExample | APIErrorForbiddenExample | void>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/runs/${moduleRunId}/tfstate`,
      method: "POST",
      secure: true,
      ...params,
    });
  /**
   * @description Gets the Github tarball URL for the module.
   *
   * @tags Modules
   * @name GetModuleTarballUrl
   * @summary Get Module Tarball URL
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}/tarball_url
   * @secure
   */
  getModuleTarballUrl = (
    teamId: string,
    moduleId: string,
    query?: {
      /**
       * the SHA to get the tarball from
       * name: github_sha
       */
      github_sha?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<GetModuleTarballURLResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/tarball_url`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Gets the current module values by ID.
   *
   * @tags Modules
   * @name GetModuleValues
   * @summary Get Module Values
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}/values/{module_values_id}
   * @secure
   */
  getModuleValues = (teamId: string, moduleId: string, moduleValuesId: string, params: RequestParams = {}) =>
    this.request<GetModuleValuesResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/values/${moduleValuesId}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Gets the current module values for the given module, by github reference or SHA.
   *
   * @tags Modules
   * @name GetCurrentModuleValues
   * @summary Get Current Module Values
   * @request GET:/api/v1/teams/{team_id}/modules/{module_id}/values/current
   * @secure
   */
  getCurrentModuleValues = (
    teamId: string,
    moduleId: string,
    query?: {
      /**
       * the SHA to get the module values file from
       * name: github_sha
       */
      github_sha?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<GetModuleValuesCurrentResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/modules/${moduleId}/values/current`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Lists monitor results for a given team, optionally filtered by module or monitor id.
   *
   * @tags Monitors
   * @name ListMonitorResults
   * @summary List Monitor Results
   * @request GET:/api/v1/teams/{team_id}/monitor_results
   * @secure
   */
  listMonitorResults = (
    teamId: string,
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      page?: number;
      /** The monitor id to filter by */
      module_monitor_id?: string;
      /** The module id to filter by */
      module_id?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<ListMonitorResultsResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/monitor_results`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Lists monitors for a given team.
   *
   * @tags Monitors
   * @name ListMonitors
   * @summary List Monitors
   * @request GET:/api/v1/teams/{team_id}/monitors
   * @secure
   */
  listMonitors = (
    teamId: string,
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      page?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<ListMonitorsResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/monitors`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Creates a new monitor.
   *
   * @tags Monitors
   * @name CreateMonitor
   * @summary Create Monitor
   * @request POST:/api/v1/teams/{team_id}/monitors
   * @secure
   */
  createMonitor = (teamId: string, data?: CreateMonitorRequest, params: RequestParams = {}) =>
    this.request<CreateMonitorResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/monitors`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Deletes a monitor.
   *
   * @tags Monitors
   * @name DeleteMonitor
   * @summary Delete Monitor
   * @request DELETE:/api/v1/teams/{team_id}/monitors/{monitor_id}
   * @secure
   */
  deleteMonitor = (teamId: string, monitorId: string, params: RequestParams = {}) =>
    this.request<DeleteMonitorResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/monitors/${monitorId}`,
      method: "DELETE",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Gets a monitor by id.
   *
   * @tags Monitors
   * @name GetMonitor
   * @summary Get Monitor
   * @request GET:/api/v1/teams/{team_id}/monitors/{monitor_id}
   * @secure
   */
  getMonitor = (teamId: string, monitorId: string, params: RequestParams = {}) =>
    this.request<GetMonitorResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/monitors/${monitorId}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Updates a monitor.
   *
   * @tags Monitors
   * @name UpdateMonitor
   * @summary Update Monitor
   * @request POST:/api/v1/teams/{team_id}/monitors/{monitor_id}
   * @secure
   */
  updateMonitor = (teamId: string, monitorId: string, data?: CreateMonitorRequest, params: RequestParams = {}) =>
    this.request<UpdateMonitorResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/monitors/${monitorId}`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Gets a notification by id.
   *
   * @tags Teams
   * @name GetNotification
   * @summary Get notification
   * @request GET:/api/v1/teams/{team_id}/notifications/{notification_id}
   * @secure
   */
  getNotification = (teamId: string, notificationId: string, params: RequestParams = {}) =>
    this.request<GetNotificationResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/teams/${teamId}/notifications/${notificationId}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Creates a new user via email and password-based authentication. This endpoint is only registered if the environment variable `BASIC_AUTH_ENABLED` is set.
   *
   * @tags Users
   * @name CreateUser
   * @summary Create a new user
   * @request POST:/api/v1/users
   * @secure
   */
  createUser = (data?: CreateUserRequest, params: RequestParams = {}) =>
    this.request<
      CreateUserResponse,
      APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample
    >({
      path: `/api/v1/users`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Deletes the current user.
   *
   * @tags Users
   * @name DeleteCurrentUser
   * @summary Delete the current user.
   * @request DELETE:/api/v1/users/current
   * @secure
   */
  deleteCurrentUser = (params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorForbiddenExample>({
      path: `/api/v1/users/current`,
      method: "DELETE",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Retrieves the current user object based on the data passed in auth.
   *
   * @tags Users
   * @name GetCurrentUser
   * @summary Retrieve the current user.
   * @request GET:/api/v1/users/current
   * @secure
   */
  getCurrentUser = (params: RequestParams = {}) =>
    this.request<GetUserResponse, APIErrorForbiddenExample>({
      path: `/api/v1/users/current`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Updates the current user.
   *
   * @tags Users
   * @name UpdateCurrentUser
   * @summary Update the current user.
   * @request POST:/api/v1/users/current
   * @secure
   */
  updateCurrentUser = (data?: UpdateOrganizationRequest, params: RequestParams = {}) =>
    this.request<UpdateUserResponse, APIErrorForbiddenExample>({
      path: `/api/v1/users/current`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Lists the github app installations for the currently authenticated user.
   *
   * @tags Users
   * @name ListGithubAppInstallations
   * @summary List Github App installations
   * @request GET:/api/v1/users/current/github_app/installations
   * @secure
   */
  listGithubAppInstallations = (
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      page?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<
      ListGithubAppInstallationsResponse,
      APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample
    >({
      path: `/api/v1/users/current/github_app/installations`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Lists organizations for a user.
   *
   * @tags Users
   * @name ListUserOrganizations
   * @summary List user organizations
   * @request GET:/api/v1/users/current/organizations
   * @secure
   */
  listUserOrganizations = (
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      page?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<ListUserOrgsResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/users/current/organizations`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Resets a password for a user using the old password as validation.
   *
   * @tags Users
   * @name ResetPasswordManual
   * @summary Reset password (manual)
   * @request POST:/api/v1/users/current/reset_password_manual
   * @secure
   */
  resetPasswordManual = (data?: ResetPasswordManualRequest, params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/users/current/reset_password_manual`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Lists personal access token for a user.
   *
   * @tags Users
   * @name ListPersonalAccessTokens
   * @summary List personal access tokens.
   * @request GET:/api/v1/users/current/settings/pats
   * @secure
   */
  listPersonalAccessTokens = (
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      page?: number;
    },
    params: RequestParams = {},
  ) =>
    this.request<ListPATsResponse, APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>({
      path: `/api/v1/users/current/settings/pats`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Creates a new personal access token for a user.
   *
   * @tags Users
   * @name CreatePersonalAccessToken
   * @summary Create a new personal access token
   * @request POST:/api/v1/users/current/settings/pats
   * @secure
   */
  createPersonalAccessToken = (data?: CreatePersonalAccessTokenRequest, params: RequestParams = {}) =>
    this.request<CreatePATResponse, APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>(
      {
        path: `/api/v1/users/current/settings/pats`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      },
    );
  /**
   * @description Deletes the personal access token for the user
   *
   * @tags Users
   * @name DeletePersonalAccessToken
   * @summary Delete the personal access token.
   * @request DELETE:/api/v1/users/current/settings/pats/{pat_id}
   * @secure
   */
  deletePersonalAccessToken = (patId: string, params: RequestParams = {}) =>
    this.request<DeletePATResponse, APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>(
      {
        path: `/api/v1/users/current/settings/pats/${patId}`,
        method: "DELETE",
        secure: true,
        format: "json",
        ...params,
      },
    );
  /**
   * @description Gets a personal access token for a user, specified by the path param `pat_id`.
   *
   * @tags Users
   * @name GetPersonalAccessToken
   * @summary Get a personal access token
   * @request GET:/api/v1/users/current/settings/pats/{pat_id}
   * @secure
   */
  getPersonalAccessToken = (patId: string, params: RequestParams = {}) =>
    this.request<GetPATResponse, APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>({
      path: `/api/v1/users/current/settings/pats/${patId}`,
      method: "GET",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Revokes the personal access token for the user
   *
   * @tags Users
   * @name RevokePersonalAccessToken
   * @summary Revoke the personal access token.
   * @request POST:/api/v1/users/current/settings/pats/{pat_id}/revoke
   * @secure
   */
  revokePersonalAccessToken = (patId: string, params: RequestParams = {}) =>
    this.request<
      RevokePATResponseExample,
      APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample
    >({
      path: `/api/v1/users/current/settings/pats/${patId}/revoke`,
      method: "POST",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Lists teams for a user, optionally filtered by organization id.
   *
   * @tags Users
   * @name ListUserTeams
   * @summary List user teams
   * @request GET:/api/v1/users/current/teams
   * @secure
   */
  listUserTeams = (
    query?: {
      /**
       * The page to query for
       * @format int64
       */
      page?: number;
      /** the id of the organization to filter by (optional) */
      organization_id?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<ListUserTeamsResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/users/current/teams`,
      method: "GET",
      query: query,
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Verifies a user's email via a token-based mechanism.
   *
   * @tags Users
   * @name VerifyEmail
   * @summary Verify email
   * @request POST:/api/v1/users/current/verify_email/finalize
   * @secure
   */
  verifyEmail = (data?: VerifyEmailRequest, params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/users/current/verify_email/finalize`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Resends a verification email for the user.
   *
   * @tags Users
   * @name ResendVerificationEmail
   * @summary Resend verification email.
   * @request POST:/api/v1/users/current/verify_email/resend
   * @secure
   */
  resendVerificationEmail = (params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/users/current/verify_email/resend`,
      method: "POST",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Logs a user in via email and password-based authentication. This endpoint is only registered if the environment variable `BASIC_AUTH_ENABLED` is set.
   *
   * @tags Users
   * @name LoginUser
   * @summary Login user
   * @request POST:/api/v1/users/login
   * @secure
   */
  loginUser = (data?: LoginUserRequest, params: RequestParams = {}) =>
    this.request<LoginUserResponse, APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>(
      {
        path: `/api/v1/users/login`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      },
    );
  /**
   * @description Logs a user out. This endpoint only performs an action if it's called with cookie-based authentication.
   *
   * @tags Users
   * @name LogoutUser
   * @summary Logout user
   * @request POST:/api/v1/users/logout
   * @secure
   */
  logoutUser = (params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/users/logout`,
      method: "POST",
      secure: true,
      format: "json",
      ...params,
    });
  /**
   * @description Resets a password for a user by sending them a verification email.
   *
   * @tags Users
   * @name ResetPasswordEmail
   * @summary Reset password (email)
   * @request POST:/api/v1/users/reset_password_email
   * @secure
   */
  resetPasswordEmail = (data?: ResetPasswordEmailRequest, params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/users/reset_password_email`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Resets a user's password given a token-based reset password mechanism.
   *
   * @tags Users
   * @name ResetPasswordEmailFinalize
   * @summary Reset password
   * @request POST:/api/v1/users/reset_password_email/finalize
   * @secure
   */
  resetPasswordEmailFinalize = (data?: ResetPasswordEmailFinalizeRequest, params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/users/reset_password_email/finalize`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Verifies that the token id and token are valid for a given reset password request.
   *
   * @tags Users
   * @name ResetPasswordEmailVerify
   * @summary Verify password reset data
   * @request POST:/api/v1/users/reset_password_email/verify
   * @secure
   */
  resetPasswordEmailVerify = (data?: ResetPasswordEmailVerifyRequest, params: RequestParams = {}) =>
    this.request<EmptyResponse, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/users/reset_password_email/verify`,
      method: "POST",
      body: data,
      secure: true,
      type: ContentType.Json,
      format: "json",
      ...params,
    });
  /**
   * @description Implements a Github App webhook.
   *
   * @tags Github Apps
   * @name GithubAppWebhook
   * @summary Github App Webhook
   * @request POST:/api/v1/webhooks/github_app
   * @secure
   */
  githubAppWebhook = (params: RequestParams = {}) =>
    this.request<void, APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>({
      path: `/api/v1/webhooks/github_app`,
      method: "POST",
      secure: true,
      ...params,
    });
}
