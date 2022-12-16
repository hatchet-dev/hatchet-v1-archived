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

func (repo *OrgRepository) CreateOrgMember(org *models.Organization, orgMember *models.OrganizationMember) (*models.OrganizationMember, repository.RepositoryError) {
	if err := repo.db.Model(org).Association("OrgMembers").Append(orgMember); err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return orgMember, nil
}

func (repo *OrgRepository) ReadOrgMemberByUserID(orgID, userID string) (*models.OrganizationMember, repository.RepositoryError) {
	member := &models.OrganizationMember{}

	if err := repo.db.Preload("OrgPolicies").Where("organization_members.organization_id = ? AND user_id = ?", orgID, userID).First(&member).Error; err != nil {
		return nil, toRepoError(repo.db, err)
	}

	return member, nil
}

func (repo *OrgRepository) ListOrgMembersByOrgID(orgID string, opts ...repository.QueryOption) ([]*models.OrganizationMember, *repository.PaginatedResult, repository.RepositoryError) {
	var members []*models.OrganizationMember

	db := repo.db.Model(&models.OrganizationMember{})

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&members).Where("organization_id = ?", orgID).Error; err != nil {
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

func (repo *OrgRepository) ListOrgPoliciesByOrgID(orgID string, opts ...repository.QueryOption) ([]*models.OrganizationPolicy, *repository.PaginatedResult, repository.RepositoryError) {
	var policies []*models.OrganizationPolicy

	db := repo.db.Model(&models.OrganizationPolicy{})

	paginatedResult := &repository.PaginatedResult{}

	db = db.Scopes(queryutils.Paginate(opts, db, paginatedResult))

	if err := db.Find(&policies).Where("organization_id = ?", orgID).Error; err != nil {
		return nil, nil, err
	}

	return policies, paginatedResult, nil
}
