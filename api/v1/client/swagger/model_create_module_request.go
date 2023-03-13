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

type CreateModuleRequest struct {
	EnvVars map[string]string `json:"env_vars,omitempty"`
	Github *CreateModuleRequestGithub `json:"github,omitempty"`
	Local *CreateModuleRequestLocal `json:"local,omitempty"`
	Name string `json:"name,omitempty"`
	ValuesGithub *CreateModuleValuesRequestGithub `json:"values_github,omitempty"`
	ValuesRaw map[string]interface{} `json:"values_raw,omitempty"`
}
