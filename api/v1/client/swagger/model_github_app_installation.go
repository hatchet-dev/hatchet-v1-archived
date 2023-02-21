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
import (
	"time"
)

type GithubAppInstallation struct {
	AccountAvatarUrl string `json:"account_avatar_url,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	// the time that this resource was created
	CreatedAt time.Time `json:"created_at,omitempty"`
	// the id of this resource, in UUID format
	Id string `json:"id,omitempty"`
	InstallationId int64 `json:"installation_id,omitempty"`
	InstallationSettingsUrl string `json:"installation_settings_url,omitempty"`
	// the time that this resource was last updated
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}