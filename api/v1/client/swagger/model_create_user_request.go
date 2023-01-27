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

type CreateUserRequest struct {
	// the display name for this user
	DisplayName string `json:"display_name"`
	// the email address for this user
	Email string `json:"email"`
	// the password for this user
	Password string `json:"password"`
}
