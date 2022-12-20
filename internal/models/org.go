package models

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/encryption"
	"golang.org/x/crypto/bcrypt"
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

	InviteLink     OrganizationInviteLink
	InviteAccepted bool

	// The referenced user
	UserID string
	User   User `gorm:"foreignKey:UserID"`

	// The attached roles for this user.
	OrgPolicies []OrganizationPolicy `gorm:"many2many:organization_member_policies;"`
}

func (o *OrganizationMember) ToAPIType(key *[32]byte) *types.OrganizationMember {
	res := &types.OrganizationMember{
		APIResourceMeta: o.ToAPITypeMetadata(),
		InviteAccepted:  o.InviteAccepted,
	}

	invite := &o.InviteLink

	if invite != nil {
		res.Invite = *invite.ToAPIType(key)
	}

	user := &o.User

	if user != nil {
		res.User = *user.ToOrgUserPublishedData()
	}

	policies := make([]types.OrganizationPolicyMeta, 0)

	for _, modelPolicy := range o.OrgPolicies {
		policies = append(policies, *modelPolicy.ToAPITypeMeta())
	}

	res.OrgPolicies = policies

	return res
}

func (o *OrganizationMember) AfterFind(tx *gorm.DB) (err error) {
	// this ensures that AfterFind is called on the invite link even if called with a
	// Joins method, instead of just Preload
	return o.InviteLink.AfterFind(tx)
}

type OrganizationInviteLink struct {
	Base
	HasEncryptedFields

	OrganizationID       string
	OrganizationMemberID string
	InviteeEmail         string

	Expires *time.Time

	// Whether the invite link has been used
	Used bool

	// Hashed before write
	Token []byte

	// Encrypted before write
	InviteLinkURL []byte
}

func NewOrganizationInviteLink(serverURL, orgID string) (*OrganizationInviteLink, error) {
	res := &OrganizationInviteLink{
		OrganizationID: orgID,
	}

	res.Base.ID = uuid.New().String()

	// generate the token
	tok, err := encryption.GenerateRandomBytes(32)

	if err != nil {
		return nil, err
	}

	// generate an ID for this token
	link := fmt.Sprintf("%s/api/v1/invites/%s/%s", serverURL, res.Base.ID, tok)

	// ensure valid url
	if _, err := url.Parse(link); err != nil {
		return nil, fmt.Errorf("not a valid url: %v", err)
	}

	res.Token = []byte(tok)
	res.InviteLinkURL = []byte(link)

	return res, nil
}

func (o *OrganizationInviteLink) ToAPIType(key *[32]byte) *types.OrganizationInvite {
	o.Decrypt(key)

	return &types.OrganizationInvite{
		InviteLinkURL: string(o.InviteLinkURL),
		InviteeEmail:  o.InviteeEmail,
		Expires:       o.Expires,
	}
}

func (o *OrganizationInviteLink) BeforeCreate(tx *gorm.DB) error {
	// return an error if encrypt has not been called
	if !o.FieldsAreEncrypted {
		return fmt.Errorf("fields should be encrypted before create")
	}

	err := o.Base.BeforeCreate(tx)

	if err != nil {
		return err
	}

	// hash the token before write
	hashedTok, err := bcrypt.GenerateFromPassword([]byte(o.Token), 8)

	if err != nil {
		return err
	}

	o.Token = hashedTok

	expiryTime := time.Now().Add(24 * time.Hour)

	o.Expires = &expiryTime

	return nil
}

func (o *OrganizationInviteLink) IsExpired() bool {
	timeLeft := o.Expires.Sub(time.Now())
	return timeLeft < 0
}

func (o *OrganizationInviteLink) VerifyToken(tok []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(o.Token, tok)

	return err == nil, err
}

func (o *OrganizationInviteLink) Encrypt(key *[32]byte) error {
	if !o.HasEncryptedFields.FieldsAreEncrypted {
		ciphertext, err := encryption.Encrypt(o.InviteLinkURL, key)

		if err != nil {
			return err
		}

		o.InviteLinkURL = ciphertext
		o.HasEncryptedFields.FieldsAreEncrypted = true
	}

	return nil
}

func (o *OrganizationInviteLink) Decrypt(key *[32]byte) error {
	if o.HasEncryptedFields.FieldsAreEncrypted {
		plaintext, err := encryption.Decrypt(o.InviteLinkURL, key)

		if err != nil {
			return err
		}

		o.InviteLinkURL = plaintext
		o.HasEncryptedFields.FieldsAreEncrypted = false
	}

	return nil
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

func (o *OrganizationPolicy) ToAPITypeMeta() *types.OrganizationPolicyMeta {
	return &types.OrganizationPolicyMeta{
		APIResourceMeta: o.ToAPITypeMetadata(),
		Name:            o.PolicyName,
	}
}
