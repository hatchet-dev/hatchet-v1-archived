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

type ModuleRunConfig struct {
	EnvVarVersionId string `json:"env_var_version_id,omitempty"`
	GithubCommitSha string `json:"github_commit_sha,omitempty"`
	TriggerKind string `json:"trigger_kind,omitempty"`
	ValuesVersionId string `json:"values_version_id,omitempty"`
}
