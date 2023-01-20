package hatchet_org_presets.team_member

import future.keywords.contains
import future.keywords.every
import future.keywords.if
import future.keywords.in

read_verbs = {"get", "list"}

write_verbs = {"create", "update", "delete"}

allow if {
	# don't allow team members to call write operations on the team itself, which includes
	# inviting other team members
	not has_team_write_scope
	not has_team_member_write_scope
}

has_team_write_scope if {
	some i
	some j
	resource := input.endpoint.resources[i]
	resource.scope == "team_scope"
	resource.verb == write_verbs[j]
}

has_team_member_write_scope if {
	some i
	some j
	resource := input.endpoint.resources[i]
	resource.scope == "team_member_scope"
	resource.verb == write_verbs[j]
}
