package types

import (
	"time"
)

type APIVerb string

const (
	APIVerbGet    APIVerb = "get"
	APIVerbCreate APIVerb = "create"
	APIVerbList   APIVerb = "list"
	APIVerbUpdate APIVerb = "update"
	APIVerbDelete APIVerb = "delete"
)

type APIVerbGroup []APIVerb

func ReadVerbGroup() APIVerbGroup {
	return []APIVerb{APIVerbGet, APIVerbList}
}

func ReadWriteVerbGroup() APIVerbGroup {
	return []APIVerb{APIVerbGet, APIVerbList, APIVerbCreate, APIVerbUpdate, APIVerbDelete}
}

type URLParam string

type HTTPVerb string

const (
	HTTPVerbGet    HTTPVerb = "GET"
	HTTPVerbPost   HTTPVerb = "POST"
	HTTPVerbPut    HTTPVerb = "PUT"
	HTTPVerbPatch  HTTPVerb = "PATCH"
	HTTPVerbDelete HTTPVerb = "DELETE"
	HTTPVerbLock   HTTPVerb = "LOCK"
	HTTPVerbUnlock HTTPVerb = "UNLOCK"
)

type PermissionScope string

const (
	NoUserScope                PermissionScope = "no_user_scope"
	UserScope                  PermissionScope = "user_scope"
	BasicAuthUserScope         PermissionScope = "basic_auth_user_scope"
	OrgScope                   PermissionScope = "org_scope"
	OrgMemberScope             PermissionScope = "org_member_scope"
	OrgOwnerScope              PermissionScope = "org_owner_scope"
	TeamScope                  PermissionScope = "team_scope"
	TeamMemberScope            PermissionScope = "team_member_scope"
	GithubAppInstallationScope PermissionScope = "github_app_installation_scope"
	ModuleScope                PermissionScope = "module_scope"
	ModuleRunScope             PermissionScope = "module_run_scope"

	// ModuleServiceAccountScope restricts the scope to service account types only. This is enforced
	// by OPA policies
	ModuleServiceAccountScope PermissionScope = "module_service_account_scope"
)

const OrgMemberLookupKey string = "org_member"
const TeamMemberLookupKey string = "team_member"

// swagger:model
type APIResourceMeta struct {
	// the time that this resource was created
	// example: 2022-12-13T15:06:48.888358-05:00
	CreatedAt *time.Time `json:"created_at"`

	// the time that this resource was last updated
	// example: 2022-12-13T15:06:48.888358-05:00
	UpdatedAt *time.Time `json:"updated_at"`

	// the id of this resource, in UUID format
	// example: bb214807-246e-43a5-a25d-41761d1cff9e
	ID string `json:"id"`
}

// swagger:model
type EmptyResponse struct{}

type UsageMetric string

// swagger:model
type APIServerMetadata struct {
	// auth metadata options
	Auth *APIServerMetadataAuth `json:"auth"`
}

// swagger:model
type APIServerMetadataAuth struct {
	// whether email verification is required in order to use the api/dashboard
	RequireEmailVerification bool `json:"require_email_verification"`
}
