import {
  StandardButton,
  FlexRowRight,
  HorizontalSpacer,
  H4,
  Placeholder,
  Spinner,
  SectionArea,
  SmallSpan,
  FlexCol,
  ErrorBar,
} from "@hatchet-dev/hatchet-components";
import { useMutation, useQuery } from "@tanstack/react-query";
import EnvVars, { getInternalEnvVars, newEnvVarAtom } from "components/envvars";
import ExpandableSettings from "components/expandablesettings";
import SelectGitSource from "components/module/selectgitpath";
import SetModuleValues from "components/module/setmodulevalues";
import React, { useMemo, useState } from "react";
import api from "shared/api";
import {
  CreateModuleRequestGithub,
  CreateModuleValuesRequestGithub,
  Module,
  UpdateModuleRequest,
  UpdateMonitorRequest,
  ModuleMonitor,
  ModuleMonitorKind,
} from "shared/api/generated/data-contracts";
import useIsModified from "shared/hooks/useismodified";
import { useAtom } from "jotai";
import CodeBlock from "components/codeblock";
import MonitorSettingsContainer from "../monitorsettings/components/MonitorSettingsContainer";

type Props = {
  team_id: string;
  monitor: ModuleMonitor;
};

const PolicyTab: React.FC<Props> = ({ team_id, monitor }) => {
  const monitor_id = monitor.id;
  const [policyBytes, setPolicyBytes] = useState(monitor.policy_bytes);
  const [err, setErr] = useState("");

  const mutation = useMutation({
    mutationKey: ["update_monitor", team_id, monitor_id],
    mutationFn: async (request: UpdateMonitorRequest) => {
      const res = await api.updateMonitor(team_id, monitor_id, request);
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

  return (
    <MonitorSettingsContainer>
      <FlexCol>
        <CodeBlock
          value={policyBytes}
          readOnly={monitor.is_default}
          onChange={(v) => setPolicyBytes(v)}
        />
        {err && <HorizontalSpacer spacepixels={20} />}
        {err && <ErrorBar text={err} />}
        {!monitor.is_default && (
          <FlexRowRight>
            <StandardButton
              label="Update"
              material_icon="chevron_right"
              icon_side="right"
              on_click={() => {
                mutation.mutate({
                  policy_bytes: policyBytes,
                });
              }}
              margin={"0"}
              is_loading={mutation.isLoading}
            />
          </FlexRowRight>
        )}
      </FlexCol>
    </MonitorSettingsContainer>
  );
};

export default PolicyTab;
