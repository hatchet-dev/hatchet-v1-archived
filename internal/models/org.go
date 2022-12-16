package models

import (
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"gorm.io/gorm"
)

type PresetPolicyName string

const (
	PresetPolicyNameOwner  PresetPolicyName = "owner"
	PresetPolicyNameAdmin  PresetPolicyName = "admin"
	PresetPolicyNameMember PresetPolicyName = "member"
)

type Organization struct {
	Base

	DisplayName string
	Icon        string

	// an organization has a single owner, which is a user model
	OwnerID string
	Owner   User `gorm:"foreignKey:OwnerID"`

	// The list of members of this organization. This is not typically returned
	// in queries unless you're explicitly calling list methods for org members.
	OrgMembers []OrganizationMember

	// The list of policies for this organization
	OrgPolicies []OrganizationPolicy
}

func (o *Organization) ToAPIType() *types.Organization {
	return &types.Organization{
		APIResourceMeta: o.Base.ToAPITypeMetadata(),
		DisplayName:     o.DisplayName,
		Owner:           *o.Owner.ToOrgUserPublishedData(),
	}
}

func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	err := o.Base.BeforeCreate(tx)

	if err != nil {
		return err
	}

	// create the preset policies
	o.OrgPolicies = []OrganizationPolicy{
		{
			PolicyName: string(PresetPolicyNameOwner),
		},
		{
			PolicyName: string(PresetPolicyNameAdmin),
		},
		{
			PolicyName: string(PresetPolicyNameMember),
		},
	}

	return nil
}

type OrganizationMember struct {
	Base

	// The parent organization ID, for the hasMany relationship
	OrganizationID string

	InviteLink     string
	InviteAccepted bool

	// The referenced user
	UserID string
	User   User `gorm:"foreignKey:UserID"`

	// The attached roles for this user.
	OrgPolicies []OrganizationPolicy `gorm:"many2many:organization_member_policies;"`
}

type OrganizationPolicy struct {
	Base

	// The parent organization ID, for the hasMany relationship
	OrganizationID string

	IsCustom   bool
	PolicyName string

	// Policy bytes MAY be empty if this is a preset policy, in which case they are preloaded
	// into the server binary.
	Policy []byte
}
