import { MaterialIcon } from "components/globals";
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
};

const MemberList: React.FC<Props> = ({ members }) => {
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
          </MemberContainer>
        );
      })}
    </MemberListContainer>
  );
};

export default MemberList;
