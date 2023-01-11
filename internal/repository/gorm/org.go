package gorm

import (
	"fmt"

	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/gorm/queryutils"
	"gorm.io/gorm"
)

// OrgRepository uses gorm.DB for querying the database
type OrgRepository struct {
	db *gorm.DB
}

// NewOrgRepository returns a DefaultOrgRepository which uses
// gorm.DB for querying the database
func NewOrgRepository(db *gorm.DB) repository.OrgRepository {
	return &OrgRepository{db}
}

// CreateOrg adds a new Org row to the Orgs table in the database
func (repo *OrgRepository) CreateOrg(org *models.Organization) (*models.Organization, repository.RepositoryError) {
	if err := repo.db.Create(org).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	// if the organization's Owner field isn't populated by OwnerID, manually query for it
	if org.Owner.ID != org.OwnerID {
		owner := &models.User{}

		if err := repo.db.Where("id = ?", org.OwnerID).First(owner).Error; err != nil {
			return nil, toRepoError(repo.db, err)
		}

		org.Owner = *owner
	}

	return org, nil
}

// ReadOrgByID finds a single org by its unique id
func (repo *OrgRepository) ReadOrgByID(id string) (*models.Organization, repository.RepositoryError) {
	org := &models.Organization{}

	if err := repo.db.Preload("Owner").Where("id = ?", id).First(&org).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return org, nil
}

// UpdateOrg updates an organization in the database
func (repo *OrgRepository) UpdateOrg(org *models.Organization) (*models.Organization, repository.RepositoryError) {
	if err := repo.db.Save(org).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return org, nil
}

// DeleteOrg deletes a single org by its unique id
func (repo *OrgRepository) DeleteOrg(org *models.Organization) (*models.Organization, repository.RepositoryError) {
	del := repo.db.Delete(&org)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return org, nil
}

func (repo *OrgRepository) ListOrgsByUserID(userID string, opts ...repository.QueryOption) ([]*models.Organization, *repository.PaginatedResult, repository.RepositoryError) {
	// get org members first, then list orgs
	var orgMembers []*models.OrganizationMember

	db := repo.db.Model(&models.OrganizationMember{}).Where("user_id = ?", userID)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&orgMembers).Error; err != nil {
		return nil, nil, err
	}

	// populate organizations based on ids
	orgIDs := make([]string, 0)

	for _, orgMember := range orgMembers {
		orgIDs = append(orgIDs, orgMember.OrganizationID)
	}

	var orgs []*models.Organization

	if err := repo.db.Preload("Owner").Where("id IN (?)", orgIDs).Find(&orgs).Error; err != nil {
		return nil, nil, toRepoError(repo.db, err)
	}

	return orgs, paginatedResult, nil
}

func (repo *OrgRepository) CreateOrgMember(org *models.Organization, orgMember *models.OrganizationMember) (*models.OrganizationMember, repository.RepositoryError) {
	if err := repo.db.Model(org).Association("OrgMembers").Append(orgMember); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return orgMember, nil
}

func (repo *OrgRepository) ReadOrgMemberByID(orgID, memberID string) (*models.OrganizationMember, repository.RepositoryError) {
	member := &models.OrganizationMember{}

	if err := repo.db.Preload("OrgPolicies").Joins("User").Joins("InviteLink").Where("organization_members.organization_id = ? AND organization_members.id = ?", orgID, memberID).First(&member).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return member, nil
}

func (repo *OrgRepository) ReadOrgMemberByUserID(orgID, userID string) (*models.OrganizationMember, repository.RepositoryError) {
	member := &models.OrganizationMember{}

	if err := repo.db.Preload("OrgPolicies").Joins("User").Joins("InviteLink").Where("organization_members.organization_id = ? AND organization_members.user_id = ?", orgID, userID).First(&member).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return member, nil
}

func (repo *OrgRepository) ReadOrgMemberByUserOrInviteeEmail(orgID, email string) (*models.OrganizationMember, repository.RepositoryError) {
	member := &models.OrganizationMember{}

	if err := repo.db.Preload("OrgPolicies").Joins("User").Joins("InviteLink").Where("organization_members.organization_id = ? AND (InviteLink.invitee_email = ? OR User.email = ?)", orgID, email, email).First(&member).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return member, nil
}

func (repo *OrgRepository) ListOrgMembersByOrgID(orgID string, opts ...repository.QueryOption) ([]*models.OrganizationMember, *repository.PaginatedResult, repository.RepositoryError) {
	var members []*models.OrganizationMember

	db := repo.db.Model(&models.OrganizationMember{}).Where("organization_id = ?", orgID)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Preload("InviteLink").Preload("User").Preload("OrgPolicies").Find(&members).Error; err != nil {
		return nil, nil, err
	}

	return members, paginatedResult, nil
}

func (repo *OrgRepository) UpdateOrgMember(orgMember *models.OrganizationMember) (*models.OrganizationMember, repository.RepositoryError) {
	if err := repo.db.Save(orgMember).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return orgMember, nil
}

func (repo *OrgRepository) DeleteOrgMember(orgMember *models.OrganizationMember) (*models.OrganizationMember, repository.RepositoryError) {
	del := repo.db.Delete(&orgMember)

	if del.Error != nil {
		return nil, toRepoError(repo.db, del.Error)
	} else if del.RowsAffected == 0 {
		return nil, repository.RepositoryNoRowsAffected
	}

	return orgMember, nil
}

func (repo *OrgRepository) ReadOrgInviteByID(inviteID string) (*models.OrganizationInviteLink, repository.RepositoryError) {
	invite := &models.OrganizationInviteLink{}

	if err := repo.db.Where("id = ?", inviteID).First(&invite).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return invite, nil
}

func (repo *OrgRepository) UpdateOrgInvite(orgInvite *models.OrganizationInviteLink) (*models.OrganizationInviteLink, repository.RepositoryError) {
	if err := repo.db.Save(orgInvite).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return orgInvite, nil
}

func (repo *OrgRepository) AppendOrgPolicyToOrgMember(orgMember *models.OrganizationMember, orgPolicy *models.OrganizationPolicy) (*models.OrganizationMember, repository.RepositoryError) {
	// we add an additional check to verify that the organization ids are the same
	if orgMember.OrganizationID != orgPolicy.OrganizationID {
		return nil, repository.UnknownRepositoryError(
			fmt.Errorf("organization ids are not equal: member has organization %s, policy has organization %s", orgMember.OrganizationID, orgPolicy.OrganizationID),
		)
	}

	if err := repo.db.Model(orgMember).Association("OrgPolicies").Append(orgPolicy); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return orgMember, nil
}

func (repo *OrgRepository) ReplaceOrgPoliciesForOrgMember(orgMember *models.OrganizationMember, policies []*models.OrganizationPolicy) (*models.OrganizationMember, repository.RepositoryError) {
	if err := repo.db.Model(orgMember).Association("OrgPolicies").Clear(); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	if err := repo.db.Model(orgMember).Association("OrgPolicies").Append(policies); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return orgMember, nil
}

func (repo *OrgRepository) RemoveOrgPolicyFromOrgMember(orgMember *models.OrganizationMember, orgPolicy *models.OrganizationPolicy) (*models.OrganizationMember, repository.RepositoryError) {
	// we add an additional check to verify that the organization ids are the same
	if orgMember.OrganizationID != orgPolicy.OrganizationID {
		return nil, repository.UnknownRepositoryError(
			fmt.Errorf("organization ids are not equal: member has organization %s, policy has organization %s", orgMember.OrganizationID, orgPolicy.OrganizationID),
		)
	}

	if err := repo.db.Model(orgMember).Association("OrgPolicies").Delete(orgPolicy); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return orgMember, nil
}

func (repo *OrgRepository) CreateOrgPolicy(org *models.Organization, orgPolicy *models.OrganizationPolicy) (*models.OrganizationPolicy, repository.RepositoryError) {
	if err := repo.db.Model(org).Association("OrgPolicies").Append(orgPolicy); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return orgPolicy, nil
}

func (repo *OrgRepository) ReadPresetPolicyByName(orgID string, presetName models.PresetPolicyName) (*models.OrganizationPolicy, repository.RepositoryError) {
	policy := &models.OrganizationPolicy{}

	if err := repo.db.Where("organization_id = ? AND is_custom = ? AND policy_name = ?", orgID, false, presetName).First(&policy).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return policy, nil
}

func (repo *OrgRepository) ReadPolicyByID(orgID, policyID string) (*models.OrganizationPolicy, repository.RepositoryError) {
	policy := &models.OrganizationPolicy{}

	if err := repo.db.Where("organization_id = ? AND id = ?", orgID, policyID).First(&policy).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return policy, nil
}

func (repo *OrgRepository) ListOrgPoliciesByOrgID(orgID string, opts ...repository.QueryOption) ([]*models.OrganizationPolicy, *repository.PaginatedResult, repository.RepositoryError) {
	var policies []*models.OrganizationPolicy

	db := repo.db.Model(&models.OrganizationPolicy{}).Where("organization_id = ?", orgID)

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&policies).Error; err != nil {
		return nil, nil, err
	}

	return policies, paginatedResult, nil
}
