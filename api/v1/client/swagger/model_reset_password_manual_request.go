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

type ResetPasswordManualRequest struct {
	// the new password for this user
	NewPassword string `json:"new_password"`
	// the old password for this user
	OldPassword string `json:"old_password"`
}