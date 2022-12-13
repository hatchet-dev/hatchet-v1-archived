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
	// a custom Hatchet error code
	// example: 1400
	Code APIErrorCode `json:"code"`

	// a description for this error
	// example: A descriptive error message
	Description string `json:"description"`

	// a link to the documentation for this error, if it exists
	// example: github.com/hatchet-dev/hatchet
	DocsLink string `json:"docs_link"`
}
