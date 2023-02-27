import {
  FlexRowRight,
  HorizontalSpacer,
  SmallSpan,
  SectionArea,
  StandardButton,
  ErrorBar,
  Confirmation,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation } from "@tanstack/react-query";

import { useHistory } from "react-router-dom";
import api from "shared/api";
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import { Module } from "shared/api/generated/data-contracts";

type Props = {
  team_id: string;
  module: Module;
};

const DeleteModuleForm: React.FC<Props> = ({ team_id, module }) => {
  const history = useHistory();
  const [prompt, setPrompt] = useState(false);
  const [err, setErr] = useState("");

  const { mutate, isLoading } = useMutation({
    mutationKey: ["delete_module", team_id, module.id],
    mutationFn: async () => {
      const res = await api.deleteModule(team_id, module.id);

      return res;
    },
    onSuccess: (data) => {
      history.push(`/team/${team_id}/modules`);
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  if (prompt) {
    return (
      <Confirmation
        prompt="Type the name of the module to trigger deletion."
        confirm_text={module.name}
        confirm_text_example="my-module"
        button_label="Delete Module"
        confirmed={() => {
          mutate();
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
          label="Delete Module"
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

export default DeleteModuleForm;
