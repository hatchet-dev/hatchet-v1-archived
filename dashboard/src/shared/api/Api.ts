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
  APIErrorBadRequestExample,
  APIErrorForbiddenExample,
  APIErrorNotSupportedExample,
  CreateOrganizationRequest,
  CreateOrganizationResponse,
  CreateOrgMemberInviteResponse,
  CreatePATResponse,
  CreatePersonalAccessTokenRequest,
  CreateUserRequest,
  CreateUserResponse,
  DeleteOrganizationResponse,
  DeletePATResponse,
  GetOrganizationResponse,
  GetPATResponse,
  GetUserResponse,
  ListOrgMembersResponse,
  ListPATsResponse,
  ListUserOrgsResponse,
  RevokePATResponseExample,
  UpdateOrganizationRequest,
  UpdateOrgMemberPoliciesResponse,
  UpdateOrgResponse,
} from "./data-contracts";
import { ContentType, HttpClient, RequestParams } from "./http-client";

export class Api<SecurityDataType = unknown> extends HttpClient<SecurityDataType> {
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
    this.request<void, APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample>({
      path: `/api/v1/invites/${orgMemberInviteId}/${orgMemberInviteTok}`,
      method: "POST",
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
    this.request<void, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/change_owner`,
      method: "POST",
      secure: true,
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
  createOrgMemberInvite = (orgId: string, params: RequestParams = {}) =>
    this.request<
      CreateOrgMemberInviteResponse,
      APIErrorBadRequestExample | APIErrorForbiddenExample | APIErrorNotSupportedExample
    >({
      path: `/api/v1/organizations/${orgId}/members`,
      method: "POST",
      secure: true,
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
    this.request<void, APIErrorBadRequestExample | APIErrorForbiddenExample>({
      path: `/api/v1/organizations/${orgId}/members/${orgMemberId}`,
      method: "DELETE",
      secure: true,
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
    this.request<void, APIErrorForbiddenExample>({
      path: `/api/v1/users/current`,
      method: "DELETE",
      secure: true,
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
}
