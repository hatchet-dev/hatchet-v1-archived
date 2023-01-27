package hatchet_org_presets.team_member

import future.keywords.contains
import future.keywords.every
import future.keywords.if
import future.keywords.in

read_verbs = {"get", "list"}

write_verbs = {"create", "update", "delete"}

all_verbs = {"get", "list", "create", "update", "delete"}

allow if {
	# don't allow team members to call write operations on the team itself, which includes
	# inviting other team members
	not has_team_write_scope
	not has_team_member_write_scope
	not has_module_service_account_scope
}

# allow if the request has a module or module run write scope, and this is not
# a service account scope. this supersedes `not has_team_write_scope`
allow if {
	has_module_write_scope
	not has_module_service_account_scope
}

# allow if this is a service account user and the endpoint has the module service account scope
allow if {
	has_module_service_account_scope
	input.user.user_account_kind == "serviceaccount"
}

has_team_write_scope if {
	some i
	some j
	resource := input.endpoint.resources[i]
	resource.scope == "team_scope"
	resource.verb == write_verbs[j]
}

has_module_write_scope if {
	some i
	some j
	resource := input.endpoint.resources[i]
	resource.scope == "module_scope"
	resource.verb == write_verbs[j]
}

has_module_service_account_scope if {
	some i
	some j
	resource := input.endpoint.resources[i]
	resource.scope == "module_service_account_scope"
	resource.verb == all_verbs[j]
}

has_team_member_write_scope if {
	some i
	some j
	resource := input.endpoint.resources[i]
	resource.scope == "team_member_scope"
	resource.verb == write_verbs[j]
}
