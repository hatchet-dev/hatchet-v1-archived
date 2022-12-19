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

export interface APIError {
  /**
   * a custom Hatchet error code
   * @format uint64
   * @example 1400
   */
  code?: number;
  /**
   * a description for this error
   * @example "A descriptive error message"
   */
  description?: string;
  /**
   * a link to the documentation for this error, if it exists
   * @example "github.com/hatchet-dev/hatchet"
   */
  docs_link?: string;
}

export interface APIErrorBadRequestExample {
  /**
   * a custom Hatchet error code
   * @format uint64
   * @example 1400
   */
  code?: number;
  /**
   * a description for this error
   * @example "Bad request (detailed error)"
   */
  description?: string;
  /**
   * a link to the documentation for this error, if it exists
   * @example "github.com/hatchet-dev/hatchet"
   */
  docs_link?: string;
}

export interface APIErrorForbiddenExample {
  /**
   * a custom Hatchet error code
   * @format uint64
   * @example 1403
   */
  code?: number;
  /**
   * a description for this error
   * @example "Forbidden"
   */
  description?: string;
  /**
   * a link to the documentation for this error, if it exists
   * @example "github.com/hatchet-dev/hatchet"
   */
  docs_link?: string;
}

export interface APIErrorNotSupportedExample {
  /**
   * a custom Hatchet error code
   * @format uint64
   * @example 1405
   */
  code?: number;
  /**
   * a description for this error
   * @example "This endpoint is not supported on this Hatchet instance."
   */
  description?: string;
  /**
   * a link to the documentation for this error, if it exists
   * @example "github.com/hatchet-dev/hatchet"
   */
  docs_link?: string;
}

export interface APIErrors {
  errors?: APIError[];
}

export interface APIResourceMeta {
  /**
   * the time that this resource was created
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  created_at?: string;
  /**
   * the id of this resource, in UUID format
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  id?: string;
  /**
   * the time that this resource was last updated
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  updated_at?: string;
}

export interface CreateOrgMemberInviteRequest {
  /**
   * the email address to use for the invite
   * @example "user1@gmail.com"
   */
  invitee_email: string;
  /** the set of policies for this user */
  invitee_policies: OrganizationPolicyReference[];
}

export type CreateOrgMemberInviteResponse = OrganizationMember;

/** @example {"display_name":"Organization 1"} */
export interface CreateOrganizationRequest {
  /**
   * the display name for this user
   * @example "Organization 1"
   */
  display_name: string;
}

export type CreateOrganizationResponse = Organization;

export interface CreatePATRequest {
  /**
   * the display name for the personal access token
   * @example "cli-token-1234"
   */
  display_name: string;
}

/** @example {"pat":{"expires":"2023-01-12T22:09:28.350Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"cli-token-1234","revoked":false},"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."} */
export interface CreatePATResponse {
  pat?: PersonalAccessToken;
  /**
   * the raw JWT token. see API documentation for details
   * @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."
   */
  token?: string;
}

/** @example {"password":"Securepassword123","display_name":"User 1","email":"user1@gmail.com"} */
export interface CreateUserRequest {
  /**
   * the display name for this user
   * @example "User 1"
   */
  display_name: string;
  /**
   * the email address for this user
   * @example "user1@gmail.com"
   */
  email: string;
  /**
   * the password for this user
   * @example "Securepassword123"
   */
  password: string;
}

export type CreateUserResponse = User;

export type DeleteOrganizationResponse = Organization;

export type DeletePATResponse = PersonalAccessToken;

export type GetOrganizationResponse = Organization;

export type GetPATResponse = PersonalAccessToken;

export type GetUserResponse = User;

/** @example {"pagination":{"next_page":3,"num_pages":10,"current_page":2},"rows":[{"updated_at":"2022-12-13T20:06:48.888Z","invite_accepted":true,"created_at":"2022-12-13T20:06:48.888Z","organization_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite":{"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite_link_url":"invite_link_url"},"user":{"display_name":"User 1","email":"user1@gmail.com"}},{"updated_at":"2022-12-13T20:06:48.888Z","invite_accepted":true,"created_at":"2022-12-13T20:06:48.888Z","organization_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite":{"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite_link_url":"invite_link_url"},"user":{"display_name":"User 1","email":"user1@gmail.com"}}]} */
export interface ListOrgMembersResponse {
  pagination?: PaginationResponse;
  rows?: OrganizationMember[];
}

/** @example {"pagination":{"next_page":3,"num_pages":10,"current_page":2},"rows":[{"expires":"2023-01-12T22:09:28.350Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"cli-token-1234","revoked":false},{"expires":"2023-01-12T22:09:28.350Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"cli-token-1234","revoked":false}]} */
export interface ListPATsResponse {
  pagination?: PaginationResponse;
  rows?: PersonalAccessToken[];
}

/** @example {"pagination":{"next_page":3,"num_pages":10,"current_page":2},"rows":[{"owner":{"display_name":"User 1","email":"user1@gmail.com"},"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Organization 1"},{"owner":{"display_name":"User 1","email":"user1@gmail.com"},"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Organization 1"}]} */
export interface ListUserOrgsResponse {
  pagination?: PaginationResponse;
  rows?: Organization[];
}

/** @example {"owner":{"display_name":"User 1","email":"user1@gmail.com"},"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Organization 1"} */
export interface Organization {
  /**
   * the time that this resource was created
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  created_at?: string;
  /**
   * the display name for the personal access token
   * @example "Organization 1"
   */
  display_name?: string;
  /**
   * the id of this resource, in UUID format
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  id?: string;
  /**
   * Public data about the user that other members of the org and team
   * can access
   */
  owner?: UserOrgPublishedData;
  /**
   * the time that this resource was last updated
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  updated_at?: string;
}

/** @example {"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite_link_url":"invite_link_url"} */
export interface OrganizationInvite {
  /**
   * the time that this resource was created
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  created_at?: string;
  /** @format date-time */
  expires?: string;
  /**
   * the id of this resource, in UUID format
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  id?: string;
  invite_link_url?: string;
  invitee_email?: string;
  /**
   * the time that this resource was last updated
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  updated_at?: string;
}

/** @example {"updated_at":"2022-12-13T20:06:48.888Z","invite_accepted":true,"created_at":"2022-12-13T20:06:48.888Z","organization_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite":{"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite_link_url":"invite_link_url"},"user":{"display_name":"User 1","email":"user1@gmail.com"}} */
export interface OrganizationMember {
  /**
   * the time that this resource was created
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  created_at?: string;
  /**
   * the id of this resource, in UUID format
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  id?: string;
  invite?: OrganizationInvite;
  invite_accepted?: boolean;
  organization_policies?: OrganizationPolicyMeta[];
  /**
   * the time that this resource was last updated
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  updated_at?: string;
  /**
   * Public data about the user that other members of the org and team
   * can access
   */
  user?: UserOrgPublishedData;
}

/** @example {"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"} */
export interface OrganizationPolicyMeta {
  /**
   * the time that this resource was created
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  created_at?: string;
  /**
   * the id of this resource, in UUID format
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  id?: string;
  name?: string;
  /**
   * the time that this resource was last updated
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  updated_at?: string;
}

export interface OrganizationPolicyReference {
  ID?: string;
  Name?: string;
}

/** @example {"next_page":3,"num_pages":10,"current_page":2} */
export interface PaginationResponse {
  /**
   * the current page
   * @format int64
   * @example 2
   */
  current_page?: number;
  /**
   * the next page
   * @format int64
   * @example 3
   */
  next_page?: number;
  /**
   * the total number of pages for listing
   * @format int64
   * @example 10
   */
  num_pages?: number;
}

/** @example {"expires":"2023-01-12T22:09:28.350Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"cli-token-1234","revoked":false} */
export interface PersonalAccessToken {
  /**
   * the time that this resource was created
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  created_at?: string;
  /**
   * the display name for the personal access token
   * @example "cli-token-1234"
   */
  display_name?: string;
  /**
   * when the token expires
   * @format date-time
   * @example "2023-01-12T22:09:28.350Z"
   */
  expires?: string;
  /**
   * the id of this resource, in UUID format
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  id?: string;
  /**
   * whether the token has been revoked
   * @example false
   */
  revoked?: boolean;
  /**
   * the time that this resource was last updated
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  updated_at?: string;
}

export type RevokePATResponse = PersonalAccessToken;

/** @example {"expires":"2023-01-12T22:09:28.350Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"cli-token-1234","revoked":true} */
export interface RevokePATResponseExample {
  /**
   * the time that this resource was created
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  created_at?: string;
  /**
   * the display name for the personal access token
   * @example "cli-token-1234"
   */
  display_name?: string;
  /**
   * when the token expires
   * @format date-time
   * @example "2023-01-12T22:09:28.350Z"
   */
  expires?: string;
  /**
   * the id of this resource, in UUID format
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  id?: string;
  /**
   * whether the token is revoked
   * @example true
   */
  revoked?: boolean;
  /**
   * the time that this resource was last updated
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  updated_at?: string;
}

export interface UpdateOrgMemberPoliciesRequest {
  /** the set of policies for this user */
  policies: OrganizationPolicyReference[];
}

export type UpdateOrgMemberPoliciesResponse = OrganizationMember;

export interface UpdateOrgOwnerRequest {
  /**
   * the member id of the new owner
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  new_owner_member_id?: string;
}

export interface UpdateOrgRequest {
  /**
   * the display name for this user
   * @example "User 1"
   */
  display_name: string;
}

export type UpdateOrgResponse = Organization;

/** @example {"email_verified":false,"updated_at":"2022-12-13T20:06:48.888Z","icon":"https://avatars.githubusercontent.com/u/25448214?v=4","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"User 1","email":"user1@gmail.com"} */
export interface User {
  /**
   * the time that this resource was created
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  created_at?: string;
  /**
   * the display name for this user
   * @example "User 1"
   */
  display_name?: string;
  /**
   * the email address for this user
   * @example "user1@gmail.com"
   */
  email?: string;
  /**
   * whether this user's email address has been verified
   * @example false
   */
  email_verified?: boolean;
  /**
   * a URI for the user icon
   * @example "https://avatars.githubusercontent.com/u/25448214?v=4"
   */
  icon?: string;
  /**
   * the id of this resource, in UUID format
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  id?: string;
  /**
   * the time that this resource was last updated
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  updated_at?: string;
}

/**
 * Public data about the user that other members of the org and team
 * can access
 * @example {"display_name":"User 1","email":"user1@gmail.com"}
 */
export interface UserOrgPublishedData {
  /**
   * the display name for this user
   * @example "User 1"
   */
  display_name?: string;
  /**
   * the email address for this user
   * @example "user1@gmail.com"
   */
  email?: string;
}

export interface CreateOrganizationRequest {
  /**
   * the display name for this user
   * @example "Organization 1"
   */
  display_name: string;
}

export interface UpdateOrganizationRequest {
  /**
   * the display name for this user
   * @example "User 1"
   */
  display_name: string;
}

export interface CreateUserRequest {
  /**
   * the display name for this user
   * @example "User 1"
   */
  display_name: string;
  /**
   * the email address for this user
   * @example "user1@gmail.com"
   */
  email: string;
  /**
   * the password for this user
   * @example "Securepassword123"
   */
  password: string;
}

export interface CreatePersonalAccessTokenRequest {
  /**
   * the display name for the personal access token
   * @example "cli-token-1234"
   */
  display_name: string;
}
