import {
  StandardButton,
  FlexRowRight,
  FlexCol,
  HorizontalSpacer,
} from "hatchet-components";
import { useMutation } from "@tanstack/react-query";
import React, { useState } from "react";
import api from "shared/api";
import { Module, ModuleRun } from "shared/api/generated/data-contracts";
import ExpandedRun from "../../../../components/run/expandedrun";
import CreateManualRun from "../createmanualrun";
import ModuleRunsList from "../modulerunslist";

type Props = {
  team_id: string;
  module: Module;
};

const RunsTab: React.FC<Props> = ({ team_id, module }) => {
  const module_id = module.id;
  const [selectedRun, setSelectedRun] = useState<ModuleRun>(null);
  const [createRun, setCreateRun] = useState(false);

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

  if (createRun) {
    return (
      <CreateManualRun
        back={() => setCreateRun(false)}
        team_id={team_id}
        module={module}
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
            setCreateRun(true);
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
