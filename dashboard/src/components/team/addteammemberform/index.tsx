import {
  FlexCol,
  HorizontalSpacer,
  MaterialIcon,
  ErrorBar,
  Selector,
  TextInput,
  SmallSpan,
} from "hatchet-components";
import React, { useEffect, useState } from "react";
import {
  AddTeamMemberRequest,
  CreateOrgMemberInviteRequest,
  OrganizationMemberSanitized,
  TeamMember,
} from "shared/api/generated/data-contracts";
import { css } from "styled-components";
import { InviteAddButton, InviteContainer } from "./styles";

export type Props = {
  submit: (member: AddTeamMemberRequest, cb: () => void) => void;
  org_members: OrganizationMemberSanitized[];
  current_team_members: TeamMember[];
  err?: string;
};

const policyOptions = [
  {
    material_icon: "person",
    label: "Admin",
    value: "admin",
  },
  {
    material_icon: "person",
    label: "Member",
    value: "member",
  },
];

const AddTeamMemberForm: React.FC<Props> = ({
  submit,
  org_members,
  current_team_members,
  err,
}) => {
  const [orgMemberID, setOrgMemberID] = useState("");
  const [policy, setPolicy] = useState("");
  const [reset, setReset] = useState(0);
  const [submitted, setSubmitted] = useState(false);

  useEffect(() => {
    if (submitted && !err) {
      setOrgMemberID("");
      setPolicy("");
      setSubmitted(false);
      setReset(reset + 1);
    }
  }, [submitted, err]);

  const onSubmit = () => {
    if (orgMemberID != "" && policy != "") {
      submit(
        {
          org_member_id: orgMemberID,
          policies: [
            {
              name: policy,
            },
          ],
        },
        () => {
          setSubmitted(true);
        }
      );
    }
  };

  const getTeamMemberLabel = (member: OrganizationMemberSanitized) => {
    if (member.user?.display_name) {
      return `${member.user?.display_name} | ${
        member.user?.email || member.invite?.invitee_email
      }`;
    }

    return member.user?.email || member.invite?.invitee_email;
  };

  const getTeamMemberOptions = () => {
    let current_member_ids = current_team_members.map(
      (team_member) => team_member.org_member.id
    );

    return org_members
      .map((org_member) => {
        return {
          label: getTeamMemberLabel(org_member),
          value: org_member.id,
          material_icon: "person",
        };
      })
      .filter((org_member) => !current_member_ids.includes(org_member.value));
  };

  return (
    <FlexCol>
      <HorizontalSpacer spacepixels={20} />
      <SmallSpan>Add a new team member</SmallSpan>
      <HorizontalSpacer spacepixels={6} />
      <InviteContainer>
        <Selector
          placeholder={"Org Member"}
          placeholder_material_icon="person"
          options={getTeamMemberOptions()}
          select={(option) => {
            setOrgMemberID(option.value);
          }}
          reset={reset}
          selector_overrides={css({
            width: "100%",
          }).toString()}
          selector_wrapper_overrides={css({
            width: "100%",
          }).toString()}
        />
        <Selector
          placeholder={"Role"}
          placeholder_material_icon="person"
          options={policyOptions}
          select={(option) => {
            setPolicy(option.value);
          }}
          reset={reset}
        />
        <InviteAddButton
          disabled={orgMemberID == "" || policy == ""}
          onClick={onSubmit}
        >
          <MaterialIcon className="material-icons">add</MaterialIcon>
        </InviteAddButton>
      </InviteContainer>
      {err && <HorizontalSpacer spacepixels={20} />}
      {err && <ErrorBar text={err} />}
    </FlexCol>
  );
};

export default AddTeamMemberForm;
