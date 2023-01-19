import {
  FlexCol,
  HorizontalSpacer,
  MaterialIcon,
  ErrorBar,
  Selector,
  TextInput,
} from "@hatchet-dev/hatchet-components";
import React, { useEffect, useState } from "react";
import {
  AddTeamMemberRequest,
  CreateOrgMemberInviteRequest,
  OrganizationMemberSanitized,
  TeamMember,
} from "shared/api/generated/data-contracts";
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

  const getTeamMemberOptions = () => {
    let current_member_ids = current_team_members.map(
      (team_member) => team_member.org_member.id
    );

    return org_members
      .map((org_member) => {
        return {
          label: org_member.user?.email || org_member.invite?.invitee_email,
          value: org_member.id,
        };
      })
      .filter((org_member) => !current_member_ids.includes(org_member.value));
  };

  return (
    <FlexCol>
      <InviteContainer>
        <Selector
          placeholder={"Org Member"}
          placeholder_material_icon="person"
          options={getTeamMemberOptions()}
          select={(option) => {
            setOrgMemberID(option.value);
          }}
          reset={reset}
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
