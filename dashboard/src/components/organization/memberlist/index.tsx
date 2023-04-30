import {
  CopyCodeline,
  CopyToClipboard,
  FlexCol,
  FlexRow,
  FlexRowRight,
  HorizontalSpacer,
  MaterialIcon,
  StandardButton,
} from "hatchet-components";
import React, { useState } from "react";
import {
  OrganizationMemberSanitized,
  OrganizationMember,
  APIErrorBadRequestExample,
  APIErrorForbiddenExample,
} from "shared/api/generated/data-contracts";
import { HttpResponse } from "shared/api/generated/http-client";
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
  get_invite_link?: (
    member: OrganizationMemberSanitized
  ) => Promise<
    HttpResponse<
      OrganizationMember,
      APIErrorBadRequestExample | APIErrorForbiddenExample
    >
  >;
};

const MemberList: React.FC<Props> = ({
  members,
  remove_member,
  get_invite_link,
}) => {
  const [orgMember, setOrgMember] = useState<OrganizationMember>();

  const renderInviteLinkCopyButton = (member: OrganizationMemberSanitized) => {
    return (
      get_invite_link &&
      member.organization_policies[0]?.name != "owner" && (
        <StandardButton
          margin="0"
          label="Show Invite Link"
          style_kind="muted"
          on_click={async () => {
            const res = await get_invite_link(member);

            setOrgMember(res.data);
          }}
        />
      )
    );
  };

  return (
    <MemberListContainer>
      {members.map((member, i) => {
        return (
          <FlexCol>
            <MemberContainer key={member.id}>
              <MemberNameOrEmail>
                {member.user?.email || member.invite?.invitee_email}
              </MemberNameOrEmail>
              <FlexRowRight gap="4px">
                <PolicyName>
                  <MaterialIcon className="material-icons">person</MaterialIcon>
                  <div>{capitalize(member.organization_policies[0]?.name)}</div>
                </PolicyName>
                {renderInviteLinkCopyButton(member)}
                {remove_member &&
                  member.organization_policies[0]?.name != "owner" && (
                    <StandardButton
                      label="Remove"
                      style_kind="muted"
                      on_click={() => remove_member(member)}
                      margin="0"
                    />
                  )}
              </FlexRowRight>
            </MemberContainer>
            {orgMember && orgMember.id == member.id && (
              <CopyCodeline
                code_line={orgMember.invite.public_invite_link_url}
              />
            )}
            {orgMember && orgMember.id == member.id && (
              <HorizontalSpacer spacepixels={16} />
            )}
          </FlexCol>
        );
      })}
    </MemberListContainer>
  );
};

export default MemberList;
