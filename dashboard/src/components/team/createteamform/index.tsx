import {
  FlexCol,
  HorizontalSpacer,
  MaterialIcon,
  ErrorBar,
  TextInput,
} from "hatchet-components";
import React, { useEffect, useState } from "react";
import { CreateTeamRequest } from "shared/api/generated/data-contracts";
import { CreateTeamContainer, TeamAddButton } from "./styles";

export type Props = {
  submit: (team: CreateTeamRequest, cb: () => void) => void;
  err?: string;
};

const CreateTeamForm: React.FC<Props> = ({ submit, err }) => {
  const [name, setName] = useState("");
  const [reset, setReset] = useState(0);
  const [submitted, setSubmitted] = useState(false);

  useEffect(() => {
    if (submitted && !err) {
      setName("");
      setSubmitted(false);
      setReset(reset + 1);
    }
  }, [submitted, err]);

  const onSubmit = () => {
    if (name != "") {
      submit(
        {
          display_name: name,
        },
        () => {
          setSubmitted(true);
        }
      );
    }
  };

  return (
    <FlexCol>
      <CreateTeamContainer>
        <TextInput
          placeholder="ex. Team 1"
          type="text"
          width="100%"
          on_change={(val) => {
            setName(val);
          }}
          reset={reset}
        />
        <TeamAddButton disabled={name == ""} onClick={onSubmit}>
          <MaterialIcon className="material-icons">add</MaterialIcon>
        </TeamAddButton>
      </CreateTeamContainer>
      {err && <HorizontalSpacer spacepixels={20} />}
      {err && <ErrorBar text={err} />}
    </FlexCol>
  );
};

export default CreateTeamForm;
