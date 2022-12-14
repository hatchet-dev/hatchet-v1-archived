package types

const (
	APIErrorCodeUnknown        uint = 0
	ErrCodeBadRequest          uint = 1400
	ErrCodeUnauthorized        uint = 1401
	ErrCodeForbidden           uint = 1403
	ErrCodeNotFound            uint = 1404
	ErrCodeUnavailable         uint = 1405
	ErrCodeInternalServerError uint = 1500
)

const GenericResourceNotFound string = "the requested resource was not found"

// swagger:model
type APIErrors struct {
	Errors []APIError `json:"errors"`
}

// swagger:model
type APIError struct {
	// a custom Hatchet error code
	// example: 1400
	Code uint `json:"code"`

	// a description for this error
	// example: A descriptive error message
	Description string `json:"description"`

	// a link to the documentation for this error, if it exists
	// example: github.com/hatchet-dev/hatchet
	DocsLink string `json:"docs_link"`
}

// swagger:model
type APIErrorBadRequestExample struct {
	*APIError

	// a custom Hatchet error code
	// example: 1400
	Code uint `json:"code"`

	// a description for this error
	// example: Bad request (detailed error)
	Description string `json:"description"`
}

// swagger:model
type APIErrorForbiddenExample struct {
	*APIError

	// a custom Hatchet error code
	// example: 1403
	Code uint `json:"code"`

	// a description for this error
	// example: Forbidden
	Description string `json:"description"`
}

// swagger:model
type APIErrorNotSupportedExample struct {
	*APIError

	// a custom Hatchet error code
	// example: 1405
	Code uint `json:"code"`

	// a description for this error
	// example: This endpoint is not supported on this Hatchet instance.
	Description string `json:"description"`
}
