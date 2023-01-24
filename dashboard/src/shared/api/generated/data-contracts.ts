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

/** @example {"auth":{"require_email_verification":true}} */
export interface APIServerMetadata {
  auth?: APIServerMetadataAuth;
}

/** @example {"require_email_verification":true} */
export interface APIServerMetadataAuth {
  /** whether email verification is required in order to use the api/dashboard */
  require_email_verification?: boolean;
}

/** @example {"invitee_email":"user1@gmail.com","invitee_policies":[{"name":"name","id":"id"},{"name":"name","id":"id"}]} */
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
   * the display name for the organization
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

/** @example {"display_name":"Team 1"} */
export interface CreateTeamRequest {
  /**
   * the display name for the team
   * @example "Team 1"
   */
  display_name: string;
}

export type CreateTeamResponse = Team;

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

export type DeleteTeamResponse = Team;

export type EmptyResponse = object;

export type GetOrgMemberResponse = OrganizationMember;

export type GetOrganizationResponse = Organization;

export type GetPATResponse = PersonalAccessToken;

export type GetUserResponse = User;

/** @example {"installation_settings_url":"installation_settings_url","updated_at":"2022-12-13T20:06:48.888Z","account_name":"account_name","created_at":"2022-12-13T20:06:48.888Z","installation_id":0,"id":"bb214807-246e-43a5-a25d-41761d1cff9e","account_avatar_url":"account_avatar_url"} */
export interface GithubAppInstallation {
  account_avatar_url?: string;
  account_name?: string;
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
  /** @format int64 */
  installation_id?: number;
  installation_settings_url?: string;
  /**
   * the time that this resource was last updated
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  updated_at?: string;
}

/** @example {"pagination":{"next_page":3,"num_pages":10,"current_page":2},"rows":[{"installation_settings_url":"installation_settings_url","updated_at":"2022-12-13T20:06:48.888Z","account_name":"account_name","created_at":"2022-12-13T20:06:48.888Z","installation_id":0,"id":"bb214807-246e-43a5-a25d-41761d1cff9e","account_avatar_url":"account_avatar_url"},{"installation_settings_url":"installation_settings_url","updated_at":"2022-12-13T20:06:48.888Z","account_name":"account_name","created_at":"2022-12-13T20:06:48.888Z","installation_id":0,"id":"bb214807-246e-43a5-a25d-41761d1cff9e","account_avatar_url":"account_avatar_url"}]} */
export interface ListGithubAppInstallationsResponse {
  pagination?: PaginationResponse;
  rows?: GithubAppInstallation[];
}

/** @example {"pagination":{"next_page":3,"num_pages":10,"current_page":2},"rows":[{"updated_at":"2022-12-13T20:06:48.888Z","invite_accepted":true,"created_at":"2022-12-13T20:06:48.888Z","organization_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite":{"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},"user":{"display_name":"User 1","email":"user1@gmail.com"}},{"updated_at":"2022-12-13T20:06:48.888Z","invite_accepted":true,"created_at":"2022-12-13T20:06:48.888Z","organization_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite":{"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},"user":{"display_name":"User 1","email":"user1@gmail.com"}}]} */
export interface ListOrgMembersResponse {
  pagination?: PaginationResponse;
  rows?: OrganizationMemberSanitized[];
}

/** @example {"pagination":{"next_page":3,"num_pages":10,"current_page":2},"rows":[{"expires":"2023-01-12T22:09:28.350Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"cli-token-1234","revoked":false},{"expires":"2023-01-12T22:09:28.350Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"cli-token-1234","revoked":false}]} */
export interface ListPATsResponse {
  pagination?: PaginationResponse;
  rows?: PersonalAccessToken[];
}

/** @example {"pagination":{"next_page":3,"num_pages":10,"current_page":2},"rows":[{"team_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"updated_at":"2022-12-13T20:06:48.888Z","org_member":{"updated_at":"2022-12-13T20:06:48.888Z","invite_accepted":true,"created_at":"2022-12-13T20:06:48.888Z","organization_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite":{"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},"user":{"display_name":"User 1","email":"user1@gmail.com"}},"created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"team_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"updated_at":"2022-12-13T20:06:48.888Z","org_member":{"updated_at":"2022-12-13T20:06:48.888Z","invite_accepted":true,"created_at":"2022-12-13T20:06:48.888Z","organization_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite":{"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},"user":{"display_name":"User 1","email":"user1@gmail.com"}},"created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}]} */
export interface ListTeamMembersResponse {
  pagination?: PaginationResponse;
  rows?: TeamMember[];
}

/** @example {"pagination":{"next_page":3,"num_pages":10,"current_page":2},"rows":[{"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Team 1"},{"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Team 1"}]} */
export interface ListTeamsResponse {
  pagination?: PaginationResponse;
  rows?: Team[];
}

/** @example {"pagination":{"next_page":3,"num_pages":10,"current_page":2},"rows":[{"owner":{"display_name":"User 1","email":"user1@gmail.com"},"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Organization 1"},{"owner":{"display_name":"User 1","email":"user1@gmail.com"},"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Organization 1"}]} */
export interface ListUserOrgsResponse {
  pagination?: PaginationResponse;
  rows?: Organization[];
}

/** @example {"pagination":{"next_page":3,"num_pages":10,"current_page":2},"rows":[{"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Team 1"},{"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Team 1"}]} */
export interface ListUserTeamsResponse {
  pagination?: PaginationResponse;
  rows?: Team[];
}

/** @example {"password":"Securepassword123","email":"user1@gmail.com"} */
export interface LoginUserRequest {
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

export type LoginUserResponse = User;

/** @example {"owner":{"display_name":"User 1","email":"user1@gmail.com"},"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Organization 1"} */
export interface Organization {
  /**
   * the time that this resource was created
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  created_at?: string;
  /**
   * the display name for the team
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

/** @example {"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"} */
export interface OrganizationInviteSanitized {
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

/**
 * OrganizationMemberSanitized represents an organization member without a sensitive invite
 * link exposed.
 * @example {"updated_at":"2022-12-13T20:06:48.888Z","invite_accepted":true,"created_at":"2022-12-13T20:06:48.888Z","organization_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite":{"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},"user":{"display_name":"User 1","email":"user1@gmail.com"}}
 */
export interface OrganizationMemberSanitized {
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
  invite?: OrganizationInviteSanitized;
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

/** @example {"name":"name","id":"id"} */
export interface OrganizationPolicyReference {
  id?: string;
  name?: string;
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

/** @example {"token_id":"bb214807-246e-43a5-a25d-41761d1cff9e","new_password":"Newpassword123","email":"user1@gmail.com","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."} */
export interface ResetPasswordEmailFinalizeRequest {
  /**
   * the email address for this user
   * @example "user1@gmail.com"
   */
  email: string;
  /**
   * the new password for this user
   * @example "Newpassword123"
   */
  new_password: string;
  /**
   * the token
   * @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."
   */
  token: string;
  /**
   * the token id
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  token_id: string;
}

/** @example {"email":"user1@gmail.com"} */
export interface ResetPasswordEmailRequest {
  /**
   * the email address for this user
   * @example "user1@gmail.com"
   */
  email: string;
}

export interface ResetPasswordEmailVerifyTokenRequest {
  /**
   * the email address for this user
   * @example "user1@gmail.com"
   */
  email: string;
  /**
   * the token
   * @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."
   */
  token: string;
  /**
   * the token id
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  token_id: string;
}

/** @example {"old_password":"Securepassword123","new_password":"Newpassword123"} */
export interface ResetPasswordManualRequest {
  /**
   * the new password for this user
   * @example "Newpassword123"
   */
  new_password: string;
  /**
   * the old password for this user
   * @example "Securepassword123"
   */
  old_password: string;
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

/** @example {"updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e","display_name":"Team 1"} */
export interface Team {
  /**
   * the time that this resource was created
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  created_at?: string;
  /**
   * the display name for the team
   * @example "Team 1"
   */
  display_name?: string;
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

export interface TeamAddMemberRequest {
  /** the organization member id of the new team member */
  org_member_id?: string;
  /** the set of policies for this user */
  policies: TeamPolicyReference[];
}

export type TeamAddMemberResponse = TeamMember;

/** @example {"team_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"updated_at":"2022-12-13T20:06:48.888Z","org_member":{"updated_at":"2022-12-13T20:06:48.888Z","invite_accepted":true,"created_at":"2022-12-13T20:06:48.888Z","organization_policies":[{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},{"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"}],"id":"bb214807-246e-43a5-a25d-41761d1cff9e","invite":{"invitee_email":"invitee_email","expires":"2000-01-23T04:56:07.000Z","updated_at":"2022-12-13T20:06:48.888Z","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"},"user":{"display_name":"User 1","email":"user1@gmail.com"}},"created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"} */
export interface TeamMember {
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
   * OrganizationMemberSanitized represents an organization member without a sensitive invite
   * link exposed.
   */
  org_member?: OrganizationMemberSanitized;
  team_policies?: TeamPolicyMeta[];
  /**
   * the time that this resource was last updated
   * @format date-time
   * @example "2022-12-13T20:06:48.888Z"
   */
  updated_at?: string;
}

/** @example {"updated_at":"2022-12-13T20:06:48.888Z","name":"name","created_at":"2022-12-13T20:06:48.888Z","id":"bb214807-246e-43a5-a25d-41761d1cff9e"} */
export interface TeamPolicyMeta {
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

export interface TeamPolicyReference {
  id?: string;
  name?: string;
}

export interface TeamUpdateRequest {
  /**
   * the display name for the team
   * @example "Team 1"
   */
  display_name: string;
}

export type TeamUpdateResponse = Team;

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

export interface UpdateUserRequest {
  /**
   * the display name for this user
   * @example "User 1"
   */
  display_name: string;
}

export type UpdateUserResponse = User;

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

/** @example {"token_id":"bb214807-246e-43a5-a25d-41761d1cff9e","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."} */
export interface VerifyEmailRequest {
  /**
   * the token
   * @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."
   */
  token: string;
  /**
   * the token id
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  token_id: string;
}

export interface CreateOrganizationRequest {
  /**
   * the display name for the organization
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

export interface CreateOrgMemberInviteRequest {
  /**
   * the email address to use for the invite
   * @example "user1@gmail.com"
   */
  invitee_email: string;
  /** the set of policies for this user */
  invitee_policies: OrganizationPolicyReference[];
}

export interface CreateTeamRequest {
  /**
   * the display name for the team
   * @example "Team 1"
   */
  display_name: string;
}

export interface AddTeamMemberRequest {
  /** the organization member id of the new team member */
  org_member_id?: string;
  /** the set of policies for this user */
  policies: TeamPolicyReference[];
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

export interface ResetPasswordManualRequest {
  /**
   * the new password for this user
   * @example "Newpassword123"
   */
  new_password: string;
  /**
   * the old password for this user
   * @example "Securepassword123"
   */
  old_password: string;
}

export interface CreatePersonalAccessTokenRequest {
  /**
   * the display name for the personal access token
   * @example "cli-token-1234"
   */
  display_name: string;
}

export interface VerifyEmailRequest {
  /**
   * the token
   * @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."
   */
  token: string;
  /**
   * the token id
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  token_id: string;
}

export interface LoginUserRequest {
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

export interface ResetPasswordEmailRequest {
  /**
   * the email address for this user
   * @example "user1@gmail.com"
   */
  email: string;
}

export interface ResetPasswordEmailFinalizeRequest {
  /**
   * the email address for this user
   * @example "user1@gmail.com"
   */
  email: string;
  /**
   * the new password for this user
   * @example "Newpassword123"
   */
  new_password: string;
  /**
   * the token
   * @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."
   */
  token: string;
  /**
   * the token id
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  token_id: string;
}

export interface ResetPasswordEmailVerifyRequest {
  /**
   * the email address for this user
   * @example "user1@gmail.com"
   */
  email: string;
  /**
   * the token
   * @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...."
   */
  token: string;
  /**
   * the token id
   * @example "bb214807-246e-43a5-a25d-41761d1cff9e"
   */
  token_id: string;
}
