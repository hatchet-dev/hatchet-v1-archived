package types

import (
	"time"
)

const (
	ErrCodeUnavailable uint = 1405
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
	UserScope PermissionScope = "user_scope"
)

// swagger:model
type APIResourceMeta struct {
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	ID        string     `json:"id"`
}

type UsageMetric string
