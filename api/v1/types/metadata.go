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
)

type PermissionScope string

const (
	NoUserScope    PermissionScope = "no_user_scope"
	UserScope      PermissionScope = "user_scope"
	OrgScope       PermissionScope = "org_scope"
	OrgMemberScope PermissionScope = "org_member_scope"
	OrgOwnerScope  PermissionScope = "org_owner_scope"
)

const OrgMemberLookupKey string = "org_member"

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
