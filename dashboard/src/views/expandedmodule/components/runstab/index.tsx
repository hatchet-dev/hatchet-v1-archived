import {
  StandardButton,
  FlexRowRight,
  FlexCol,
  HorizontalSpacer,
} from "@hatchet-dev/hatchet-components";
import { useMutation } from "@tanstack/react-query";
import React, { useState } from "react";
import api from "shared/api";
import { Module, ModuleRun } from "shared/api/generated/data-contracts";
import ExpandedRun from "../../../../components/run/expandedrun";
import ModuleRunsList from "../modulerunslist";

type Props = {
  team_id: string;
  module: Module;
};

const RunsTab: React.FC<Props> = ({ team_id, module }) => {
  const module_id = module.id;
  const [selectedRun, setSelectedRun] = useState<ModuleRun>(null);
  const [err, setErr] = useState("");

  const mutation = useMutation({
    mutationKey: ["create_module_run", team_id, module_id],
    mutationFn: async () => {
      const res = await api.createModuleRun(team_id, module_id);
      return res;
    },
    onSuccess: (data) => {
      setErr("");
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  if (selectedRun) {
    return (
      <ExpandedRun
        back={() => setSelectedRun(null)}
        team_id={team_id}
        module={module}
        module_run_id={selectedRun.id}
      />
    );
  }

  return (
    <FlexCol height="100%">
      <HorizontalSpacer spacepixels={20} />
      <FlexRowRight>
        <StandardButton
          label="New Run"
          material_icon="cached"
          on_click={() => {
            mutation.mutate();
          }}
        />
      </FlexRowRight>
      <ModuleRunsList
        team_id={team_id}
        module_id={module_id}
        select_run={setSelectedRun}
      />
    </FlexCol>
  );
};

export default RunsTab;
