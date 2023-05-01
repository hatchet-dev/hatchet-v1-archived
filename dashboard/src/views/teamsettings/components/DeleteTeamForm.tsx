import {
  FlexRowRight,
  HorizontalSpacer,
  SmallSpan,
  SectionArea,
  StandardButton,
  ErrorBar,
  Confirmation,
} from "hatchet-components";
import React, { useState } from "react";
import { useMutation } from "@tanstack/react-query";

import { useHistory } from "react-router-dom";
import api from "shared/api";
import { currOrgAtom, currTeamAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";

const DeleteTeamForm: React.FunctionComponent = () => {
  const history = useHistory();
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);
  const [prompt, setPrompt] = useState(false);
  const [err, setErr] = useState("");

  const { mutate, isLoading } = useMutation(api.deleteTeam, {
    mutationKey: ["delete_team", currTeam.id],
    onSuccess: (data) => {
      history.push("/");
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  const renderContents = () => {
    if (prompt) {
      return (
        <Confirmation
          prompt="Type the name of the team to trigger deletion."
          confirm_text={currTeam.display_name}
          confirm_text_example={currTeam.display_name}
          button_label="Delete Team"
          confirmed={() => {
            mutate(currTeam.id);
          }}
        />
      );
    }

    return (
      <>
        <SmallSpan>This operation cannot be undone.</SmallSpan>
        {err && <HorizontalSpacer spacepixels={12} />}
        {err && <ErrorBar text={err} />}
        <HorizontalSpacer spacepixels={12} />
        <FlexRowRight>
          <StandardButton
            label="Delete Team"
            material_icon="delete"
            icon_side="right"
            on_click={() => {
              setPrompt(true);
            }}
            margin={"0"}
            is_loading={isLoading}
          />
        </FlexRowRight>
      </>
    );
  };

  return <SectionArea>{renderContents()}</SectionArea>;
};

export default DeleteTeamForm;
