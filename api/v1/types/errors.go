package types

type APIErrorCode uint

const (
	APIErrorCodeUnknown        APIErrorCode = 0
	ErrCodeBadRequest          APIErrorCode = 1400
	ErrCodeUnauthorized        APIErrorCode = 1401
	ErrCodeForbidden           APIErrorCode = 1403
	ErrCodeNotFound            APIErrorCode = 1404
	ErrCodeUnavailable         APIErrorCode = 1405
	ErrCodeInternalServerError APIErrorCode = 1500
)

// swagger:model
type APIErrors struct {
	Errors []APIError `json:"errors"`
}

// swagger:model
type APIError struct {
	Code        APIErrorCode `json:"code"`
	Description string       `json:"description"`
	DocsLink    string       `json:"docs_link"`
}
