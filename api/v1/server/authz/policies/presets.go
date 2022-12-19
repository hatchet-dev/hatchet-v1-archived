package policies

import (
	_ "embed"
	"log"

	"github.com/hatchet-dev/hatchet/internal/opa"
)

type Policy struct {
	Query *opa.OPAQuery
}

var PresetPolicies struct {
	OrgOwnerPolicy  *Policy
	OrgAdminPolicy  *Policy
	OrgMemberPolicy *Policy
}

//go:embed presets/owner_policy.rego
var PresetOwnerPolicy []byte

//go:embed presets/admin_policy.rego
var PresetAdminPolicy []byte

//go:embed presets/member_policy.rego
var PresetMemberPolicy []byte

func init() {
	orgOwnerPolicy, err := opa.LoadQueryFromBytes("hatchet_org_presets.owner", PresetOwnerPolicy)

	if err != nil {
		log.Fatal(err)
	}

	orgAdminPolicy, err := opa.LoadQueryFromBytes("hatchet_org_presets.admin", PresetAdminPolicy)

	if err != nil {
		log.Fatal(err)
	}

	orgMemberPolicy, err := opa.LoadQueryFromBytes("hatchet_org_presets.member", PresetMemberPolicy)

	if err != nil {
		log.Fatal(err)
	}

	PresetPolicies.OrgOwnerPolicy = &Policy{orgOwnerPolicy}
	PresetPolicies.OrgAdminPolicy = &Policy{orgAdminPolicy}
	PresetPolicies.OrgMemberPolicy = &Policy{orgMemberPolicy}
}
