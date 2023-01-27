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

type LockTerraformStateRequest struct {
	Created string `json:"Created,omitempty"`
	ID string `json:"ID,omitempty"`
	Info string `json:"Info,omitempty"`
	Operation string `json:"Operation,omitempty"`
	Path string `json:"Path,omitempty"`
	Version string `json:"Version,omitempty"`
	Who string `json:"Who,omitempty"`
}
