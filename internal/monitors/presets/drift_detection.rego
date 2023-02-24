package hatchet.module

import future.keywords

POLICY_SEVERITY := "high"

POLICY_TITLE := "drift_detection"

POLICY_SUCCESS_MESSAGE := "No drift detected"

POLICY_ALLOW if {
	not input.plan.resource_drift
}

POLICY_ALLOW if {
	count(input.plan.resource_drift) == 0
}

POLICY_FAILURE_MESSAGE contains msg if {
	not POLICY_ALLOW
	msg := "Failed: drift detected in at least one resource"
}

POLICY_STATUS := "succeeded" if POLICY_ALLOW

else := "failed" {
	true
}
