package models

import (
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"gorm.io/gorm"
)

type PresetTeamPolicyName string

const (
	PresetTeamPolicyNameAdmin  PresetTeamPolicyName = "admin"
	PresetTeamPolicyNameMember PresetTeamPolicyName = "member"
)

type Team struct {
	Base

	DisplayName string

	// The parent organization id
	OrganizationID string

	// The list of members of this team
	TeamMembers []TeamMember

	// The list of policies for this team
	TeamPolicies []TeamPolicy
}

func (t *Team) ToAPIType() *types.Team {
	return &types.Team{
		APIResourceMeta: t.Base.ToAPITypeMetadata(),
		DisplayName:     t.DisplayName,
	}
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	err := t.Base.BeforeCreate(tx)

	if err != nil {
		return err
	}

	// create the preset policies
	t.TeamPolicies = []TeamPolicy{
		{
			PolicyName: string(PresetTeamPolicyNameAdmin),
		},
		{
			PolicyName: string(PresetTeamPolicyNameMember),
		},
	}

	return nil
}

type TeamMember struct {
	Base

	// The parent team ID, for the hasMany relationship
	TeamID string

	// The referenced org member
	OrgMemberID string
	OrgMember   OrganizationMember `gorm:"foreignKey:OrgMemberID"`

	// The attached roles for this user.
	TeamPolicies []TeamPolicy `gorm:"many2many:team_member_policies;"`
}

func (o *TeamMember) ToAPIType() *types.TeamMember {
	res := &types.TeamMember{
		APIResourceMeta: o.ToAPITypeMetadata(),
		OrgMember:       *o.OrgMember.ToAPITypeSanitized(),
	}

	policies := make([]types.TeamPolicyMeta, 0)

	for _, modelPolicy := range o.TeamPolicies {
		policies = append(policies, *modelPolicy.ToAPITypeMeta())
	}

	res.TeamPolicies = policies

	return res
}

type TeamPolicy struct {
	Base

	// The team organization ID, for the hasMany relationship
	TeamID string

	IsCustom   bool
	PolicyName string

	// Policy bytes MAY be empty if this is a preset policy, in which case they are preloaded
	// into the server binary.
	Policy []byte
}

func (o *TeamPolicy) ToAPITypeMeta() *types.TeamPolicyMeta {
	return &types.TeamPolicyMeta{
		APIResourceMeta: o.ToAPITypeMetadata(),
		Name:            o.PolicyName,
	}
}
