import {
  ErrorBar,
  HorizontalSpacer,
  MaterialIcon,
  StandardButton,
} from "hatchet-components";
import React from "react";
import {
  OrganizationMemberSanitized,
  TeamMember,
} from "shared/api/generated/data-contracts";
import { capitalize } from "shared/utils";
import { lookupOrgMember } from "../teamwithmembers";
import {
  MemberListContainer,
  MemberContainer,
  MemberNameOrEmail,
  PolicyName,
} from "./styles";

export type Props = {
  members: TeamMember[];
  org_members: OrganizationMemberSanitized[];
  add_member?: boolean;
  remove_member?: (member: TeamMember) => void;
  err?: string;
};

const TeamMemberList: React.FC<Props> = ({
  members,
  org_members,
  remove_member,
  add_member = false,
  err,
}) => {
  return (
    <MemberListContainer>
      {members.map((member, i) => {
        const orgMember = lookupOrgMember(org_members, member.org_member.id);

        return (
          <MemberContainer key={member.id}>
            <MemberNameOrEmail>
              {orgMember?.user?.email || orgMember?.invite?.invitee_email}
            </MemberNameOrEmail>
            <PolicyName>
              <MaterialIcon className="material-icons">person</MaterialIcon>
              <div>{capitalize(member.team_policies[0]?.name)}</div>
            </PolicyName>
            {remove_member && (
              <StandardButton
                label="Remove"
                style_kind="muted"
                on_click={() => remove_member(member)}
              />
            )}
          </MemberContainer>
        );
      })}
      {err && <HorizontalSpacer spacepixels={20} />}
      {err && <ErrorBar text={err} />}
    </MemberListContainer>
  );
};

export default TeamMemberList;
