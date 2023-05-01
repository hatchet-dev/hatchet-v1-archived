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
} from "hatchet-components";
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
  setMonitorName: (name: string) => void;
  setMonitorDescription: (description: string) => void;
};

const UpdateMonitorMeta: React.FC<Props> = ({
  monitor,
  setMonitorName,
  setMonitorDescription,
}) => {
  return (
    <>
      <TextInput
        placeholder="ex. my-monitor"
        initial_value={monitor.name}
        label="Monitor name"
        type="text"
        width="400px"
        on_change={(val) => {
          setMonitorName(val);
        }}
      />
      <HorizontalSpacer spacepixels={16} />
      <TextInput
        initial_value={monitor.description}
        label="Monitor description"
        width="600px"
        placeholder="ex. Detects drift"
        type="text"
        on_change={(val) => {
          setMonitorDescription(val);
        }}
      />
    </>
  );
};

export default UpdateMonitorMeta;
