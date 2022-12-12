package endpoint

import (
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/hatchet-dev/hatchet/internal/config/server"
)

// APIEndpoint contains the implementation of a Hatchet API endpoint, including which
// decoder, validator, and writer to use.
type APIEndpoint struct {
	Metadata         *EndpointMetadata
	DecoderValidator handlerutils.RequestDecoderValidator
	Writer           handlerutils.ResultWriter
}

// APIEndpointFactory contains common methods for generating an API endpoint from a shared
// set of metadata. All server endpoints should likely use the DefaultAPIEndpointFactory, this
// abstraction is mostly for testing purposes.
type APIEndpointFactory interface {
	NewAPIEndpoint(metadata *EndpointMetadata) *APIEndpoint
	GetDecoderValidator() handlerutils.RequestDecoderValidator
	GetResultWriter() handlerutils.ResultWriter
}

type APIObjectEndpointFactory struct {
	DecoderValidator handlerutils.RequestDecoderValidator
	ResultWriter     handlerutils.ResultWriter
}

func NewAPIObjectEndpointFactory(conf *server.Config) APIEndpointFactory {
	decoderValidator := handlerutils.NewDefaultRequestDecoderValidator(conf.Logger, conf.ErrorAlerter)
	resultWriter := handlerutils.NewDefaultResultWriter(conf.Logger, conf.ErrorAlerter)

	return &APIObjectEndpointFactory{
		DecoderValidator: decoderValidator,
		ResultWriter:     resultWriter,
	}
}

func (factory *APIObjectEndpointFactory) NewAPIEndpoint(
	metadata *EndpointMetadata,
) *APIEndpoint {
	return &APIEndpoint{
		Metadata:         metadata,
		DecoderValidator: factory.DecoderValidator,
		Writer:           factory.ResultWriter,
	}
}

func (factory *APIObjectEndpointFactory) GetDecoderValidator() handlerutils.RequestDecoderValidator {
	return factory.DecoderValidator
}

func (factory *APIObjectEndpointFactory) GetResultWriter() handlerutils.ResultWriter {
	return factory.ResultWriter
}
