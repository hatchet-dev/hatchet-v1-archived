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

type TeamAddMemberRequest struct {
	// the organization member id of the new team member
	OrgMemberId string `json:"org_member_id,omitempty"`
	// the set of policies for this user
	Policies []TeamPolicyReference `json:"policies"`
}
