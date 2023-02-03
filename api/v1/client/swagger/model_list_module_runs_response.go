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

type ListModuleRunsResponse struct {
	Pagination *PaginationResponse `json:"pagination,omitempty"`
	Rows []ModuleRun `json:"rows,omitempty"`
}
