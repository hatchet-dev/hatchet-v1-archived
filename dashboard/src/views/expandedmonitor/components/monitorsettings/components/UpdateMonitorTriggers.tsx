import {
  StandardButton,
  FlexRowRight,
  FlexCol,
  HorizontalSpacer,
  SmallSpan,
  TextInput,
  Selector,
  Selection,
  P,
} from "@hatchet-dev/hatchet-components";
import { useMutation } from "@tanstack/react-query";
import React, { useState } from "react";
import api from "shared/api";
import {
  ModuleMonitor,
  ModuleMonitorKind,
} from "shared/api/generated/data-contracts";
import { MonitorKindOptions } from "views/createmonitor/components/MonitorMeta";

type Props = {
  monitor: ModuleMonitor;
  setMonitorKind: (kind: ModuleMonitorKind) => void;
  setMonitorSchedule: (schedule: string) => void;
};

const UpdateMonitorTriggers: React.FC<Props> = ({
  monitor,
  setMonitorKind,
  setMonitorSchedule,
}) => {
  const [kind, setKind] = useState(monitor.kind);

  const selectKind = (option: Selection) => {
    setMonitorKind(option.value);
    setKind(option.value);
  };

  const isCronMonitor = kind == "plan" || kind == "state";

  return (
    <>
      <P>
        Choose when this policy check should be run. You can configure this to
        run periodically against a Terraform plan or the Terraform state, or you
        can run checks before/after Terraform operations.
      </P>
      <HorizontalSpacer spacepixels={12} />
      <Selector
        placeholder="Scheduled Plan Check"
        placeholder_material_icon="schedule"
        options={MonitorKindOptions}
        select={selectKind}
      />
      {isCronMonitor && (
        <>
          <HorizontalSpacer spacepixels={20} />
          <P>Provide a cron schedule to run these policy checks.</P>
          <HorizontalSpacer spacepixels={12} />
          <TextInput
            placeholder="ex. * * * * *"
            initial_value={monitor.cron_schedule}
            width="400px"
            on_change={(val) => {
              setMonitorSchedule(val);
            }}
          />
        </>
      )}
    </>
  );
};

export default UpdateMonitorTriggers;
