package terraform

import (
	"encoding/json"

	"github.com/hatchet-dev/hatchet/api/v1/types"
)

func GetPlanSummaryFromBytes(jsonBytes []byte) (*types.ModulePlanSummary, error) {
	plan := &terraformPlanRepresentation{}

	err := json.Unmarshal(jsonBytes, plan)

	if err != nil {
		return nil, err
	}

	res := make(types.ModulePlanSummary, 0)

	for _, resourceChange := range plan.ResourceChanges {
		res = append(res, types.ModulePlannedChangeSummary{
			Address: resourceChange.Address,
			Actions: resourceChange.Change.Actions,
		})
	}

	return &res, nil
}

// ref: https://developer.hashicorp.com/terraform/internals/json-format#plan-representation
// TODO(abelanger5): incomplete
type terraformPlanRepresentation struct {
	ResourceChanges []terraformPlanResourceChanges `json:"resource_changes"`
}

type terraformPlanResourceChanges struct {
	Address string              `json:"address"`
	Change  terraformPlanChange `json:"change"`
}

// ref: https://developer.hashicorp.com/terraform/internals/json-format#change-representation
type terraformPlanChange struct {
	Actions []string `json:"actions"`
}
