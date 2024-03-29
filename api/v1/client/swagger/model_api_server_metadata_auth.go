/*
 * API v1
 *
 * # Introduction Welcome to the documentation for Hatchet's API.  
 *
 * API version: 1.0.0
 * Contact: support@hatchet.run
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type ApiServerMetadataAuth struct {
	// whether email verification is required in order to use the api/dashboard
	RequireEmailVerification bool `json:"require_email_verification,omitempty"`
}
