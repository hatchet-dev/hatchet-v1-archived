package types

type APIErrorCode uint

const (
	APIErrorCodeUnknown APIErrorCode = 0
)

type APIErrors struct {
	Errors []APIError `json:"errors"`
}

type APIError struct {
	Code        APIErrorCode `json:"code"`
	Description string       `json:"description"`
	DocsLink    string       `json:"docs_link"`
}
