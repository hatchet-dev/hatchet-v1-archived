import StandardButton from "components/buttons";
import {
  FlexRowRight,
  MaterialIcon,
  StyledDeprecatedText,
} from "components/globals";
import React from "react";
import { OrganizationMemberSanitized } from "shared/api/generated/data-contracts";
import { capitalize } from "shared/utils";
import {
  MemberListContainer,
  MemberContainer,
  MemberNameOrEmail,
  PolicyName,
} from "./styles";

export type Props = {
  members: OrganizationMemberSanitized[];
  remove_member?: (member: OrganizationMemberSanitized) => void;
};

const MemberList: React.FC<Props> = ({ members, remove_member }) => {
  return (
    <MemberListContainer>
      {members.map((member, i) => {
        return (
          <MemberContainer key={member.id}>
            <MemberNameOrEmail>
              {member.user?.email || member.invite?.invitee_email}
            </MemberNameOrEmail>
            <PolicyName>
              <MaterialIcon className="material-icons">person</MaterialIcon>
              <div>{capitalize(member.organization_policies[0]?.name)}</div>
            </PolicyName>
            {remove_member &&
              member.organization_policies[0]?.name != "owner" && (
                <StandardButton
                  label="Remove"
                  style_kind="muted"
                  on_click={() => remove_member(member)}
                />
              )}
          </MemberContainer>
        );
      })}
    </MemberListContainer>
  );
};

export default MemberList;
