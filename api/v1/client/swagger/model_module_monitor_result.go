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

type ModuleMonitorResult struct {
	// the time that this resource was created
	CreatedAt time.Time `json:"created_at,omitempty"`
	// the id of this resource, in UUID format
	Id string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	ModuleId string `json:"module_id,omitempty"`
	ModuleMonitorId string `json:"module_monitor_id,omitempty"`
	ModuleName string `json:"module_name,omitempty"`
	ModuleRunId string `json:"module_run_id,omitempty"`
	Severity string `json:"severity,omitempty"`
	Status string `json:"status,omitempty"`
	Title string `json:"title,omitempty"`
	// the time that this resource was last updated
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
