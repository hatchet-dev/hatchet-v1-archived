package handlerutils

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/erroralerter"
	"github.com/hatchet-dev/hatchet/internal/logger"
)

type RequestDecoderValidator interface {
	DecodeAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool
	DecodeAndValidateQueryOnly(w http.ResponseWriter, r *http.Request, v interface{}) bool
	DecodeAndValidateNoWrite(r *http.Request, v interface{}) error
}

type DefaultRequestDecoderValidator struct {
	logger    logger.Logger
	alerter   erroralerter.Alerter
	validator Validator
	decoder   Decoder
}

func NewDefaultRequestDecoderValidator(
	logger logger.Logger,
	alerter erroralerter.Alerter,
) RequestDecoderValidator {
	validator := NewDefaultValidator()
	decoder := NewDefaultDecoder()

	return &DefaultRequestDecoderValidator{logger, alerter, validator, decoder}
}

func (j *DefaultRequestDecoderValidator) DecodeAndValidate(
	w http.ResponseWriter,
	r *http.Request,
	v interface{},
) (ok bool) {
	var requestErr apierrors.RequestError

	// decode the request parameters (body and query)
	if requestErr = j.decoder.Decode(v, r); requestErr != nil {
		apierrors.HandleAPIError(j.logger, j.alerter, w, r, requestErr, true)
		return false
	}

	// validate the request object
	if requestErr = j.validator.Validate(v); requestErr != nil {
		apierrors.HandleAPIError(j.logger, j.alerter, w, r, requestErr, true)
		return false
	}

	return true
}

func (j *DefaultRequestDecoderValidator) DecodeAndValidateQueryOnly(
	w http.ResponseWriter,
	r *http.Request,
	v interface{},
) (ok bool) {
	var requestErr apierrors.RequestError

	// decode the request parameters (body and query)
	if requestErr = j.decoder.DecodeQueryOnly(v, r); requestErr != nil {
		apierrors.HandleAPIError(j.logger, j.alerter, w, r, requestErr, true)
		return false
	}

	// validate the request object
	if requestErr = j.validator.Validate(v); requestErr != nil {
		apierrors.HandleAPIError(j.logger, j.alerter, w, r, requestErr, true)
		return false
	}

	return true
}

func (j *DefaultRequestDecoderValidator) DecodeAndValidateNoWrite(
	r *http.Request,
	v interface{},
) error {
	var requestErr apierrors.RequestError

	// decode the request parameters (body and query)
	if requestErr = j.decoder.Decode(v, r); requestErr != nil {
		return fmt.Errorf(requestErr.InternalError())
	}

	// validate the request object
	if requestErr = j.validator.Validate(v); requestErr != nil {
		return fmt.Errorf(requestErr.InternalError())
	}

	return nil
}
