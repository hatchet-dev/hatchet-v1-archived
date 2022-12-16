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
	OrgOwnerPolicy *Policy
}

//go:embed presets/owner_policy.rego
var PresetOwnerPolicy []byte

func init() {
	orgOwnerPolicy, err := opa.LoadQueryFromBytes("hatchet_org_presets.owner", PresetOwnerPolicy)

	if err != nil {
		log.Fatal(err)
	}

	PresetPolicies.OrgOwnerPolicy = &Policy{orgOwnerPolicy}
}
