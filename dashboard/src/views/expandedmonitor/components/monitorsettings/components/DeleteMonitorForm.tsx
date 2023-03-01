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
import { ModuleMonitor } from "shared/api/generated/data-contracts";

type Props = {
  team_id: string;
  monitor: ModuleMonitor;
};

const DeleteMonitorForm: React.FC<Props> = ({ team_id, monitor }) => {
  const history = useHistory();
  const [prompt, setPrompt] = useState(false);
  const [err, setErr] = useState("");

  const { mutate, isLoading } = useMutation({
    mutationKey: ["delete_monitor", team_id, monitor.id],
    mutationFn: async () => {
      const res = await api.deleteMonitor(team_id, monitor.id);

      return res;
    },
    onSuccess: (data) => {
      history.push(`/team/${team_id}/monitors`);
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
        prompt="Type the name of the monitor to trigger deletion."
        confirm_text={monitor.name}
        confirm_text_example="my-monitor"
        button_label="Delete Monitor"
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
          label="Delete Monitor"
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

export default DeleteMonitorForm;
