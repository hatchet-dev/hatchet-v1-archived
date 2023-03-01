package policies

import (
	_ "embed"
	"log"

	"github.com/hatchet-dev/hatchet/internal/opa"
)

type Policy struct {
	Query *opa.OPAQuery
}

var PresetOrgPolicies struct {
	OrgOwnerPolicy  *Policy
	OrgAdminPolicy  *Policy
	OrgMemberPolicy *Policy
}

var PresetTeamPolicies struct {
	TeamAdminPolicy  *Policy
	TeamMemberPolicy *Policy
}

//go:embed presets/org_owner_policy.rego
var PresetOrgOwnerPolicy []byte

//go:embed presets/org_admin_policy.rego
var PresetOrgAdminPolicy []byte

//go:embed presets/org_member_policy.rego
var PresetOrgMemberPolicy []byte

//go:embed presets/team_admin_policy.rego
var PresetTeamAdminPolicy []byte

//go:embed presets/team_member_policy.rego
var PresetTeamMemberPolicy []byte

func init() {
	orgOwnerPolicy, err := opa.LoadQueryFromBytes(opa.PACKAGE_HATCHET_ORGANIZATION, PresetOrgOwnerPolicy)

	if err != nil {
		log.Fatal(err)
	}

	orgAdminPolicy, err := opa.LoadQueryFromBytes(opa.PACKAGE_HATCHET_ORGANIZATION, PresetOrgAdminPolicy)

	if err != nil {
		log.Fatal(err)
	}

	orgMemberPolicy, err := opa.LoadQueryFromBytes(opa.PACKAGE_HATCHET_ORGANIZATION, PresetOrgMemberPolicy)

	if err != nil {
		log.Fatal(err)
	}

	PresetOrgPolicies.OrgOwnerPolicy = &Policy{orgOwnerPolicy}
	PresetOrgPolicies.OrgAdminPolicy = &Policy{orgAdminPolicy}
	PresetOrgPolicies.OrgMemberPolicy = &Policy{orgMemberPolicy}

	teamAdminPolicy, err := opa.LoadQueryFromBytes(opa.PACKAGE_HATCHET_TEAM, PresetTeamAdminPolicy)

	if err != nil {
		log.Fatal(err)
	}

	teamMemberPolicy, err := opa.LoadQueryFromBytes(opa.PACKAGE_HATCHET_TEAM, PresetTeamMemberPolicy)

	if err != nil {
		log.Fatal(err)
	}

	PresetTeamPolicies.TeamAdminPolicy = &Policy{teamAdminPolicy}
	PresetTeamPolicies.TeamMemberPolicy = &Policy{teamMemberPolicy}
}
