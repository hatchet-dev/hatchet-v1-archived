package hatchet_org_presets.admin

import future.keywords.contains
import future.keywords.every
import future.keywords.if
import future.keywords.in

allow if {
	resource := input.endpoint.resources[_]

	# admins cannot perform any actions on the org owner resource
	resource.scope != "org_owner_scope"
}
