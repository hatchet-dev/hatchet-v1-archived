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

// Public data about the user that other members of the org and team can access
type UserOrgPublishedData struct {
	// the display name for this user
	DisplayName string `json:"display_name,omitempty"`
	// the email address for this user
	Email string `json:"email,omitempty"`
}
