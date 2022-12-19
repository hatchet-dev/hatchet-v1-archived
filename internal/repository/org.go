package repository

import "github.com/hatchet-dev/hatchet/internal/models"

// OrgRepository represents the set of queries on the Organization model
type OrgRepository interface {
	CreateOrg(org *models.Organization) (*models.Organization, RepositoryError)
	ReadOrgByID(id string) (*models.Organization, RepositoryError)
	ListOrgsByUserID(userID string, opts ...QueryOption) ([]*models.Organization, *PaginatedResult, RepositoryError)
	UpdateOrg(org *models.Organization) (*models.Organization, RepositoryError)
	DeleteOrg(org *models.Organization) (*models.Organization, RepositoryError)

	CreateOrgMember(org *models.Organization, orgMember *models.OrganizationMember) (*models.OrganizationMember, RepositoryError)
	ReadOrgMemberByUserID(orgID, userID string) (*models.OrganizationMember, RepositoryError)
	ReadOrgMemberByID(orgID, memberID string) (*models.OrganizationMember, RepositoryError)
	ListOrgMembersByOrgID(orgID string, opts ...QueryOption) ([]*models.OrganizationMember, *PaginatedResult, RepositoryError)
	UpdateOrgMember(orgMember *models.OrganizationMember) (*models.OrganizationMember, RepositoryError)
	DeleteOrgMember(orgMember *models.OrganizationMember) (*models.OrganizationMember, RepositoryError)
	AppendOrgPolicyToOrgMember(orgMember *models.OrganizationMember, orgPolicy *models.OrganizationPolicy) (*models.OrganizationMember, RepositoryError)
	ReplaceOrgPoliciesForOrgMember(orgMember *models.OrganizationMember, policies []*models.OrganizationPolicy) (*models.OrganizationMember, RepositoryError)
	RemoveOrgPolicyFromOrgMember(orgMember *models.OrganizationMember, orgPolicy *models.OrganizationPolicy) (*models.OrganizationMember, RepositoryError)

	ReadOrgInviteByID(inviteID string) (*models.OrganizationInviteLink, RepositoryError)
	UpdateOrgInvite(orgInvite *models.OrganizationInviteLink) (*models.OrganizationInviteLink, RepositoryError)

	CreateOrgPolicy(org *models.Organization, orgPolicy *models.OrganizationPolicy) (*models.OrganizationPolicy, RepositoryError)
	ReadPresetPolicyByName(orgID string, presetName models.PresetPolicyName) (*models.OrganizationPolicy, RepositoryError)
	ReadPolicyByID(orgID, policyID string) (*models.OrganizationPolicy, RepositoryError)
	ListOrgPoliciesByOrgID(orgID string, opts ...QueryOption) ([]*models.OrganizationPolicy, *PaginatedResult, RepositoryError)
}
