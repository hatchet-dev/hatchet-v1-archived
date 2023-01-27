package hatchet_org_presets.team_admin

import future.keywords.contains
import future.keywords.every
import future.keywords.if
import future.keywords.in

read_verbs = {"get", "list"}

write_verbs = {"create", "update", "delete"}

all_verbs = {"get", "list", "create", "update", "delete"}

# TODO(abelanger5): don't allow admin members to invite other members to the team

allow if {
	not has_module_service_account_scope
}

has_module_service_account_scope if {
	some i
	some j
	resource := input.endpoint.resources[i]
	resource.scope == "module_service_account_scope"
	resource.verb == all_verbs[j]
}
