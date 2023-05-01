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
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";

const DeleteOrganizationForm: React.FunctionComponent = () => {
  const history = useHistory();
  const [currOrg, setCurrOrg] = useAtom(currOrgAtom);
  const [prompt, setPrompt] = useState(false);
  const [err, setErr] = useState("");

  const { mutate, isLoading } = useMutation(api.deleteOrg, {
    mutationKey: ["delete_organization", currOrg.id],
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
          prompt="Type the name of the organization to trigger deletion."
          confirm_text={currOrg.display_name}
          confirm_text_example="My Organization"
          button_label="Delete Organization"
          confirmed={() => {
            mutate(currOrg.id);
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
            label="Delete Organization"
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

export default DeleteOrganizationForm;
