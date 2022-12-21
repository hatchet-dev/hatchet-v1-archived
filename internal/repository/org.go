package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// OrgRepository represents the set of queries on the Organization model
type OrgRepository interface {
	// --- Org queries ---
	//
	// CreateOrg creates a new organization in the database
	CreateOrg(org *models.Organization) (*models.Organization, RepositoryError)

	// ReadOrgByID reads the organization by it's unique UUID
	ReadOrgByID(id string) (*models.Organization, RepositoryError)

	// ListOrgsByUserID lists all organizations for a user
	ListOrgsByUserID(userID string, opts ...QueryOption) ([]*models.Organization, *PaginatedResult, RepositoryError)

	// UpdateOrg updates any modified values for an organization
	UpdateOrg(org *models.Organization) (*models.Organization, RepositoryError)

	// DeleteOrg soft-deletes an organization
	DeleteOrg(org *models.Organization) (*models.Organization, RepositoryError)

	// --- Org member queries ---
	//
	// CreateOrgMember creates a new organization member for that organization
	CreateOrgMember(org *models.Organization, orgMember *models.OrganizationMember) (*models.OrganizationMember, RepositoryError)

	// ReadOrgMemberByUserOrInviteeEmail finds an org member by email address. The email address
	// can be set either in the OrganizationMember.User.Email field, or the OrganizationMember.InviteLink.InviteeEmail
	// field. Some users won't have InviteLink set (i.e. owners of the organization), and not all org members
	// will have user fields set (those with open invitations).
	ReadOrgMemberByUserOrInviteeEmail(orgID, email string) (*models.OrganizationMember, RepositoryError)

	// ReadOrgMemberByUserID finds an org member by their user UUID. Not to be confused with ReadOrgMemberByID,
	// which finds an org member by the **org member ID**.
	ReadOrgMemberByUserID(orgID, userID string) (*models.OrganizationMember, RepositoryError)

	// ReadOrgMemberByID finds an org member by their unique org member UUID.
	ReadOrgMemberByID(orgID, memberID string) (*models.OrganizationMember, RepositoryError)

	// ListOrgMembersByOrgID lists org members that are part of that organization.
	ListOrgMembersByOrgID(orgID string, opts ...QueryOption) ([]*models.OrganizationMember, *PaginatedResult, RepositoryError)

	// UpdateOrgMember updates org members. This MAY have the side effect of updating dependent models,
	// depending on the implementation. Gorm is inconsistent about this so make sure any update methods
	// are tested.
	UpdateOrgMember(orgMember *models.OrganizationMember) (*models.OrganizationMember, RepositoryError)

	// DeleteOrgMember deletes an org member.
	DeleteOrgMember(orgMember *models.OrganizationMember) (*models.OrganizationMember, RepositoryError)

	// AppendOrgPolicyToOrgMember adds an org policy to that member.
	AppendOrgPolicyToOrgMember(orgMember *models.OrganizationMember, orgPolicy *models.OrganizationPolicy) (*models.OrganizationMember, RepositoryError)

	// ReplaceOrgPoliciesForOrgMember replaces all org policies for that org member.
	ReplaceOrgPoliciesForOrgMember(orgMember *models.OrganizationMember, policies []*models.OrganizationPolicy) (*models.OrganizationMember, RepositoryError)

	// RemoveOrgPolicyFromOrgMember removes an org policy for that member.
	RemoveOrgPolicyFromOrgMember(orgMember *models.OrganizationMember, orgPolicy *models.OrganizationPolicy) (*models.OrganizationMember, RepositoryError)

	// --- Invite link queries ---
	//
	// ReadOrgInviteByID finds an org invite link by its UUID.
	ReadOrgInviteByID(inviteID string) (*models.OrganizationInviteLink, RepositoryError)

	// UpdateOrgInvite updates an org invite link. It is beneficial to use this method instead of
	// UpdateOrgMember when you don't want to make additional DB queries and you already have access
	// to the OrganizationInviteLink model.
	UpdateOrgInvite(orgInvite *models.OrganizationInviteLink) (*models.OrganizationInviteLink, RepositoryError)

	// --- Org policy queries ---
	//
	// CreateOrgPolicy creates a new organization policy
	CreateOrgPolicy(org *models.Organization, orgPolicy *models.OrganizationPolicy) (*models.OrganizationPolicy, RepositoryError)

	// ReadPresetPolicyByName finds a preset policy for an organization by its PolicyName
	ReadPresetPolicyByName(orgID string, presetName models.PresetPolicyName) (*models.OrganizationPolicy, RepositoryError)

	// ReadPolicyByID finds a policy by its ID
	ReadPolicyByID(orgID, policyID string) (*models.OrganizationPolicy, RepositoryError)

	// ListOrgPoliciesByOrgID lists policies for the organization
	ListOrgPoliciesByOrgID(orgID string, opts ...QueryOption) ([]*models.OrganizationPolicy, *PaginatedResult, RepositoryError)
}
