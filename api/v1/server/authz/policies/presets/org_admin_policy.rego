package hatchet_org_presets.org_admin

import future.keywords.contains
import future.keywords.every
import future.keywords.if
import future.keywords.in

allow if {
	not has_org_owner_scope
}

has_org_owner_scope if {
	some i
	resource := input.endpoint.resources[i]
	resource.scope == "org_owner_scope"
}
