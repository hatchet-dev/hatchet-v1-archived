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

type ResetPasswordEmailVerifyRequest struct {
	// the email address for this user
	Email string `json:"email"`
	// the token
	Token string `json:"token"`
	// the token id
	TokenId string `json:"token_id"`
}
