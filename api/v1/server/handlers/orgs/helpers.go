package orgs

import (
	"net/http"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/v1/types"
	"github.com/hatchet-dev/hatchet/internal/models"
	"github.com/hatchet-dev/hatchet/internal/repository"
)

func getPoliciesFromRequest(orgRepo repository.OrgRepository, orgID string, refs []types.OrganizationPolicyReference) []*models.OrganizationPolicy {
	policies := make([]*models.OrganizationPolicy, 0)

	// load any preset policies
	for _, policyReq := range refs {
		var policy *models.OrganizationPolicy
		var err error

		if policyReq.Name != "" {
			policy, err = orgRepo.ReadPresetPolicyByName(orgID, models.PresetPolicyName(policyReq.Name))
		} else if policyReq.ID != "" {
			policy, err = orgRepo.ReadPolicyByID(orgID, policyReq.ID)
		}

		if err == nil {
			policies = append(policies, policy)
		}
	}

	return policies
}

func verifyNotOwner(orgMember *models.OrganizationMember) apierrors.RequestError {
	// ensure that the org member is not an owner
	isOwner := false

	for _, policy := range orgMember.OrgPolicies {
		if !policy.IsCustom && policy.PolicyName == string(models.PresetPolicyNameOwner) {
			isOwner = true
			break
		}
	}

	if isOwner {
		return apierrors.NewErrPassThroughToClient(
			types.APIError{
				Code:        types.ErrCodeBadRequest,
				Description: "Cannot perform this operation on the organization owner. Please change the owner of the organization before attempting this operation.",
			},
			http.StatusBadRequest,
		)
	}

	return nil
}
