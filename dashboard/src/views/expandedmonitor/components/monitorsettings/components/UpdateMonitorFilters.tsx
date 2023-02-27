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
import { useMutation, useQuery } from "@tanstack/react-query";
import React, { useState } from "react";
import api from "shared/api";
import {
  ModuleMonitor,
  ModuleMonitorKind,
} from "shared/api/generated/data-contracts";
import { MonitorKindOptions } from "views/createmonitor/components/MonitorMeta";

type Props = {
  team_id: string;
  monitor: ModuleMonitor;
  setMonitorModules: (modules: string[]) => void;
};

const UpdateMonitorFilters: React.FC<Props> = ({
  team_id,
  monitor,
  setMonitorModules,
}) => {
  const [moduleOptions, setModuleOptions] = useState<Selection[]>([]);

  useQuery({
    queryKey: ["modules", team_id],
    queryFn: async () => {
      const res = await api.listModules(team_id);

      return res;
    },
    retry: false,
    onSuccess: (data) => {
      const modOptions: Selection[] = data.data.rows.map((module) => {
        const checked = monitor.modules.includes(module.id);

        return {
          label: module.name,
          value: module.id,
          checked: checked,
        };
      });

      setModuleOptions(modOptions);
    },
    refetchOnWindowFocus: false,
  });

  return (
    <>
      <P>Provide a list of modules that this monitor should trigger on.</P>
      <HorizontalSpacer spacepixels={12} />
      <Selector
        placeholder="All Modules"
        placeholder_material_icon="indeterminate_check_box"
        options={moduleOptions}
        selector_kind="multi"
        selector_multi_object_kind="Modules Selected"
        select={(option) => {
          const newOptions = moduleOptions.map((moduleOption) => {
            if (moduleOption.value == option.value) {
              moduleOption.checked = !moduleOption.checked;
            }

            return moduleOption;
          });

          setModuleOptions(newOptions);
          setMonitorModules(
            newOptions
              .filter((option) => option.checked)
              .map((option) => option.value)
          );
        }}
      />
    </>
  );
};

export default UpdateMonitorFilters;
