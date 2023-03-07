import {
  StandardButton,
  FlexRowRight,
  FlexCol,
  HorizontalSpacer,
  Selector,
  SectionArea,
  P,
  BackText,
  H4,
  ErrorBar,
} from "@hatchet-dev/hatchet-components";
import { useMutation } from "@tanstack/react-query";
import React, { useState } from "react";
import api from "shared/api";
import {
  CreateModuleRunRequest,
  Module,
  ModuleRun,
} from "shared/api/generated/data-contracts";
import ExpandedRun from "../../../../components/run/expandedrun";
import ModuleRunsList from "../modulerunslist";

type Props = {
  back: () => void;
  team_id: string;
  module: Module;
};

export const ManualRunKindOptions = [
  {
    label: "Plan",
    value: "plan",
    material_icon: "check",
  },
  {
    label: "Apply",
    value: "apply",
    material_icon: "check",
  },
];

const CreateManualRun: React.FC<Props> = ({ back, team_id, module }) => {
  const module_id = module.id;
  const [err, setErr] = useState("");
  const [kind, setKind] = useState("");

  const mutation = useMutation({
    mutationKey: ["create_module_run", team_id, module_id],
    mutationFn: async (req: CreateModuleRunRequest) => {
      const res = await api.createModuleRun(team_id, module_id, req);
      return res;
    },
    onSuccess: (data) => {
      setErr("");
      back();
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  return (
    <FlexCol height="100%">
      <HorizontalSpacer spacepixels={24} />
      <BackText text="All Runs" back={back} />
      <HorizontalSpacer spacepixels={24} />
      <SectionArea>
        <H4>Module Run Configuration</H4>
        <HorizontalSpacer spacepixels={12} />
        <P>Select the module run kind.</P>
        <HorizontalSpacer spacepixels={12} />
        <Selector
          placeholder="Run Kind"
          placeholder_material_icon="check"
          options={ManualRunKindOptions}
          select={(opt) => setKind(opt.value)}
        />
      </SectionArea>
      {err && <HorizontalSpacer spacepixels={20} />}
      {err && <ErrorBar text={err} />}
      <HorizontalSpacer spacepixels={20} />
      <FlexRowRight>
        <StandardButton
          label="Create run"
          material_icon="chevron_right"
          icon_side="right"
          on_click={() => {
            mutation.mutate({
              kind: kind,
            });
          }}
          disabled={!kind}
        />
      </FlexRowRight>
    </FlexCol>
  );
};

export default CreateManualRun;
