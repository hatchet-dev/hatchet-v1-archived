import {
  FlexCol,
  HorizontalSpacer,
  MaterialIcon,
  ErrorBar,
  Selector,
  TextInput,
} from "@hatchet-dev/hatchet-components";
import React, { useEffect, useState } from "react";
import { CreateOrgMemberInviteRequest } from "shared/api/generated/data-contracts";
import { InviteAddButton, InviteContainer } from "./styles";

export type Props = {
  submit: (invite: CreateOrgMemberInviteRequest, cb: () => void) => void;
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
  {
    material_icon: "person",
    label: "Viewer",
    value: "viewer",
  },
];

const InviteMemberForm: React.FC<Props> = ({ submit, err }) => {
  const [email, setEmail] = useState("");
  const [policy, setPolicy] = useState("");
  const [reset, setReset] = useState(0);
  const [submitted, setSubmitted] = useState(false);

  useEffect(() => {
    if (submitted && !err) {
      setEmail("");
      setPolicy("");
      setSubmitted(false);
      setReset(reset + 1);
    }
  }, [submitted, err]);

  const onSubmit = () => {
    if (email != "" && policy != "") {
      submit(
        {
          invitee_email: email,
          invitee_policies: [
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

  return (
    <FlexCol>
      <InviteContainer>
        <TextInput
          placeholder="you@example.com"
          type="text"
          width="100%"
          on_change={(val) => {
            setEmail(val);
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
          disabled={email == "" || policy == ""}
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

export default InviteMemberForm;
