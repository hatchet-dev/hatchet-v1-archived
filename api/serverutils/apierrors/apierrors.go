package apierrors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hatchet-dev/hatchet/api/serverutils/erroralerter"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/logger"
)

type RequestError interface {
	Error() string
	APIError() types.APIErrors
	InternalError() string
	GetStatusCode() int
}

type ErrInternal struct {
	err error
}

func NewErrInternal(err error) RequestError {
	return &ErrInternal{err}
}

func (e *ErrInternal) Error() string {
	return e.err.Error()
}

func (e *ErrInternal) InternalError() string {
	return e.err.Error()
}

func (e *ErrInternal) APIError() types.APIErrors {
	return types.APIErrors{
		Errors: []types.APIError{
			{
				Code:        types.ErrCodeInternalServerError,
				Description: "An internal error occurred.",
			},
		},
	}
}

func (e *ErrInternal) GetStatusCode() int {
	return http.StatusInternalServerError
}

type ErrForbidden struct {
	err error
}

func NewErrForbidden(err error) RequestError {
	return &ErrForbidden{err}
}

func (e *ErrForbidden) Error() string {
	return e.err.Error()
}

func (e *ErrForbidden) InternalError() string {
	return e.err.Error()
}

func (e *ErrForbidden) APIError() types.APIErrors {
	return types.APIErrors{
		Errors: []types.APIError{
			{
				Code:        types.ErrCodeForbidden,
				Description: "Forbidden",
			},
		},
	}
}

func (e *ErrForbidden) GetStatusCode() int {
	return http.StatusForbidden
}

// errors that should be passed directly, with no filter
type ErrPassThroughToClient struct {
	apiErrors  types.APIErrors
	statusCode int
	errDetails []string
}

func NewErrPassThroughToClient(apiError types.APIError, statusCode int, details ...string) RequestError {
	return &ErrPassThroughToClient{types.APIErrors{
		Errors: []types.APIError{apiError},
	}, statusCode, details}
}

func NewErrPassThroughToClientMulti(apiErrors types.APIErrors, statusCode int, details ...string) RequestError {
	return &ErrPassThroughToClient{apiErrors, statusCode, details}
}

func (e *ErrPassThroughToClient) Error() string {
	errs := make([]string, 0)

	for _, apiError := range e.apiErrors.Errors {
		errs = append(errs, apiError.Description)
	}

	return strings.Join(errs, ", ")
}

func (e *ErrPassThroughToClient) InternalError() string {
	if len(e.errDetails) > 0 {
		return fmt.Sprintf("%v: %s", e.Error(), strings.Join(e.errDetails, ", "))
	}

	return e.Error()
}

func (e *ErrPassThroughToClient) APIError() types.APIErrors {
	return e.apiErrors
}

func (e *ErrPassThroughToClient) GetStatusCode() int {
	return e.statusCode
}

// errors that denote that a resource was not found
type ErrNotFound struct {
	err error
}

func NewErrNotFound(err error) RequestError {
	return &ErrNotFound{err}
}

func (e *ErrNotFound) Error() string {
	return e.err.Error()
}

func (e *ErrNotFound) InternalError() string {
	return e.err.Error()
}

func (e *ErrNotFound) APIError() types.APIErrors {
	return types.APIErrors{
		Errors: []types.APIError{
			{
				Code:        types.ErrCodeNotFound,
				Description: "Resource not found",
			},
		},
	}
}

func (e *ErrNotFound) GetStatusCode() int {
	return http.StatusNotFound
}

type ErrorOpts struct {
	Code uint
}

func HandleAPIError(
	l logger.Logger,
	alerter erroralerter.Alerter,
	w http.ResponseWriter,
	r *http.Request,
	err RequestError,
	writeErr bool,
	opts ...ErrorOpts,
) {
	// log the internal error
	event := l.Warn().
		Str("internal_error", err.InternalError())

	data := make(map[string]interface{})
	// data := logger.AddLoggingContextScopes(r.Context(), event)
	// logger.AddLoggingRequestMeta(r, event)

	event.Send()

	// if the status code is internal server error, use alerter
	if err.GetStatusCode() == http.StatusInternalServerError && alerter != nil {
		data["method"] = r.Method
		data["url"] = r.URL.String()

		alerter.SendAlert(r.Context(), err, data)
	}

	if writeErr {
		// send the external error
		resp := err.APIError()

		// write the status code
		w.WriteHeader(err.GetStatusCode())

		writerErr := json.NewEncoder(w).Encode(resp)

		if writerErr != nil {
			event := l.Error().
				Err(writerErr)

			// logger.AddLoggingContextScopes(r.Context(), event)
			// logger.AddLoggingRequestMeta(r, event)

			event.Send()
		}
	}

	return
}
