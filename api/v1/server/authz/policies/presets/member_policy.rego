package hatchet_org_presets.member

import future.keywords.contains
import future.keywords.every
import future.keywords.if
import future.keywords.in

read_verbs = {"get", "list"}

write_verbs = {"create", "update", "delete"}

allow if {
	not has_org_write_scope
	not has_org_owner_scope
}

has_org_owner_scope if {
	some i
	resource := input.endpoint.resources[i]
	resource.scope == "org_owner_scope"
}

has_org_write_scope if {
	some i
	some j
	resource := input.endpoint.resources[i]
	resource.scope == "org_scope"
	resource.verb == write_verbs[j]
}
