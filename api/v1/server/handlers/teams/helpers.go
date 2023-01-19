package teams

import (
	"fmt"
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

func getPoliciesFromRequest(teamRepo repository.TeamRepository, teamID string, refs []types.TeamPolicyReference) ([]*models.TeamPolicy, apierrors.RequestError) {
	policies := make([]*models.TeamPolicy, 0)

	errs := make([]types.APIError, 0)

	// load any preset policies
	for _, policyReq := range refs {
		var policy *models.TeamPolicy
		var err error

		if policyReq.Name != "" {
			policy, err = teamRepo.ReadPresetTeamPolicyByName(teamID, models.PresetTeamPolicyName(policyReq.Name))

			if err != nil {
				errs = append(errs, types.APIError{
					Code:        types.ErrCodeBadRequest,
					Description: fmt.Sprintf("Policy not found with name %s", policyReq.Name),
				})
			}
		} else if policyReq.ID != "" {
			policy, err = teamRepo.ReadPolicyByID(teamID, policyReq.ID)

			if err != nil {
				errs = append(errs, types.APIError{
					Code:        types.ErrCodeBadRequest,
					Description: fmt.Sprintf("Policy not found with id %s", policyReq.ID),
				})
			}
		}

		if err == nil {
			policies = append(policies, policy)
		}
	}

	if len(errs) > 0 {
		return nil, apierrors.NewErrPassThroughToClientMulti(types.APIErrors{
			Errors: errs,
		}, http.StatusBadRequest)
	}

	return policies, nil
}
