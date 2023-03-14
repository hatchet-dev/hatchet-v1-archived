package types

import "time"

const (
	URLParamOrgID              URLParam = "org_id"
	URLParamOrgMemberID        URLParam = "org_member_id"
	URLParamOrgMemberInviteID  URLParam = "org_member_invite_id"
	URLParamOrgMemberInviteTok URLParam = "org_member_invite_tok"
)

// swagger:model
type Organization struct {
	*APIResourceMeta

	// the display name for the team
	// example: Organization 1
	DisplayName string `json:"display_name"`

	// information about the organization owner
	Owner UserOrgPublishedData `json:"owner"`
}

// OrganizationMemberSanitized represents an organization member without a sensitive invite
// link exposed.
// swagger:model
type OrganizationMemberSanitized struct {
	*APIResourceMeta

	User UserOrgPublishedData `json:"user"`

	InviteAccepted bool                        `json:"invite_accepted"`
	Invite         OrganizationInviteSanitized `json:"invite"`

	OrgPolicies []OrganizationPolicyMeta `json:"organization_policies"`
}

// swagger:model
type OrganizationInviteSanitized struct {
	*APIResourceMeta

	InviteeEmail string `json:"invitee_email"`

	Expires *time.Time `json:"expires"`
}

// swagger:model
type OrganizationMember struct {
	*APIResourceMeta

	User UserOrgPublishedData `json:"user"`

	Invite         OrganizationInvite `json:"invite"`
	InviteAccepted bool               `json:"invite_accepted"`

	OrgPolicies []OrganizationPolicyMeta `json:"organization_policies"`
}

// swagger:model
type OrganizationPolicyMeta struct {
	*APIResourceMeta

	Name string `json:"name"`
}

// swagger:model
type OrganizationInvite struct {
	*APIResourceMeta

	InviteLinkURL       string `json:"invite_link_url"`
	PublicInviteLinkURL string `json:"public_invite_link_url"`
	InviteeEmail        string `json:"invitee_email"`

	Expires *time.Time `json:"expires"`
}

// swagger:model
type CreateOrganizationRequest struct {
	// the display name for the organization
	//
	// required: true
	// example: Organization 1
	DisplayName string `json:"display_name" form:"required,max=255"`
}

// swagger:model
type CreateOrganizationResponse Organization

// swagger:model
type GetOrganizationResponse Organization

// swagger:model
type DeleteOrganizationResponse Organization

// swagger:parameters listUserOrganizations
type ListUserOrgsRequest struct {
	*PaginationRequest
}

// swagger:model
type ListUserOrgsResponse struct {
	Pagination *PaginationResponse `json:"pagination"`
	Rows       []*Organization     `json:"rows"`
}

// swagger:model
type UpdateOrgOwnerRequest struct {
	// the member id of the new owner
	// example: bb214807-246e-43a5-a25d-41761d1cff9e
	NewOwnerMemberID string `json:"new_owner_member_id" form:"required,uuid"`
}

// swagger:model
type CreateOrgMemberInviteRequest struct {
	// the email address to use for the invite
	//
	// required: true
	// example: user1@gmail.com
	InviteeEmail string `json:"invitee_email" form:"required,max=255,email"`

	// the set of policies for this user
	// required: true
	InviteePolicies []OrganizationPolicyReference `json:"invitee_policies" form:"required,min=1,dive"`
}

// swagger:model
type OrganizationPolicyReference struct {
	Name string `json:"name" form:"omitempty,oneof=admin member"`
	ID   string `json:"id" form:"omitempty,uuid"`
}

// swagger:model
type CreateOrgMemberInviteResponse OrganizationMember

// swagger:parameters listOrgMembers
type ListOrgMembersRequest struct {
	*PaginationRequest
}

// swagger:model
type ListOrgMembersResponse struct {
	Pagination *PaginationResponse            `json:"pagination"`
	Rows       []*OrganizationMemberSanitized `json:"rows"`
}

// swagger:model
type GetOrgMemberResponse OrganizationMember

// swagger:model
type UpdateOrgMemberPoliciesRequest struct {
	// the set of policies for this user
	// required: true
	Policies []OrganizationPolicyReference `json:"policies" form:"required,min=1,dive"`
}

// swagger:model
type UpdateOrgMemberPoliciesResponse OrganizationMember

// swagger:model
type UpdateOrgRequest struct {
	// the display name for this user
	//
	// required: true
	// example: User 1
	DisplayName string `json:"display_name" form:"required,max=255"`
}

// swagger:model
type UpdateOrgResponse Organization
