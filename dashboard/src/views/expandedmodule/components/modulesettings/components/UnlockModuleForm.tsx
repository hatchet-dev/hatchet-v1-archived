import {
  FlexRowRight,
  HorizontalSpacer,
  SmallSpan,
  SectionArea,
  StandardButton,
  ErrorBar,
  Confirmation,
  Placeholder,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";

import { useHistory } from "react-router-dom";
import api from "shared/api";
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import { Module } from "shared/api/generated/data-contracts";

type Props = {
  team_id: string;
  module: Module;
};

const UnlockModuleForm: React.FC<Props> = ({ team_id, module }) => {
  const history = useHistory();
  const [prompt, setPrompt] = useState(false);
  const [err, setErr] = useState("");

  const { refetch } = useQuery({
    queryKey: ["module", team_id, module.id],
    queryFn: async () => {
      const res = await api.getModule(team_id, module.id);
      return res;
    },
    retry: false,
  });

  const { mutate, isLoading } = useMutation({
    mutationKey: ["force_unlock_module", team_id, module.id],
    mutationFn: async () => {
      const res = await api.forceUnlockModule(team_id, module.id);

      return res;
    },
    onSuccess: (data) => {
      refetch();
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  if (!module.lock_id) {
    return (
      <Placeholder>
        <SmallSpan>Your module is currently unlocked.</SmallSpan>
      </Placeholder>
    );
  }

  if (prompt) {
    return (
      <Confirmation
        prompt={`Type ${module.lock_id} to force unlock the module.`}
        confirm_text={module.lock_id}
        confirm_text_example={module.lock_id}
        button_label="Unlock Module"
        button_material_icon="lock_open"
        confirmed={() => {
          mutate();
        }}
      />
    );
  }

  return (
    <>
      <SmallSpan>
        Your module is currently locked with lock id {module.lock_id}.
      </SmallSpan>
      {err && <HorizontalSpacer spacepixels={12} />}
      {err && <ErrorBar text={err} />}
      <HorizontalSpacer spacepixels={12} />
      <FlexRowRight>
        <StandardButton
          label="Unlock Module"
          material_icon="lock_open"
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

export default UnlockModuleForm;
