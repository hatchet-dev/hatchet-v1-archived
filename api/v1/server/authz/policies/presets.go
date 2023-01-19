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
	orgOwnerPolicy, err := opa.LoadQueryFromBytes("hatchet_org_presets.org_owner", PresetOrgOwnerPolicy)

	if err != nil {
		log.Fatal(err)
	}

	orgAdminPolicy, err := opa.LoadQueryFromBytes("hatchet_org_presets.org_admin", PresetOrgAdminPolicy)

	if err != nil {
		log.Fatal(err)
	}

	orgMemberPolicy, err := opa.LoadQueryFromBytes("hatchet_org_presets.org_member", PresetOrgMemberPolicy)

	if err != nil {
		log.Fatal(err)
	}

	PresetOrgPolicies.OrgOwnerPolicy = &Policy{orgOwnerPolicy}
	PresetOrgPolicies.OrgAdminPolicy = &Policy{orgAdminPolicy}
	PresetOrgPolicies.OrgMemberPolicy = &Policy{orgMemberPolicy}

	teamAdminPolicy, err := opa.LoadQueryFromBytes("hatchet_org_presets.team_admin", PresetTeamAdminPolicy)

	if err != nil {
		log.Fatal(err)
	}

	teamMemberPolicy, err := opa.LoadQueryFromBytes("hatchet_org_presets.team_member", PresetTeamMemberPolicy)

	if err != nil {
		log.Fatal(err)
	}

	PresetTeamPolicies.TeamAdminPolicy = &Policy{teamAdminPolicy}
	PresetTeamPolicies.TeamMemberPolicy = &Policy{teamMemberPolicy}
}
