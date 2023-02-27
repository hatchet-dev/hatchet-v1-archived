import {
  StandardButton,
  FlexRowRight,
  HorizontalSpacer,
  H4,
  Placeholder,
  Spinner,
  SectionArea,
  SmallSpan,
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
import UpdateMonitorMeta from "./components/UpdateMonitorMeta";
import MonitorSettingsContainer from "./components/MonitorSettingsContainer";
import UpdateMonitorTriggers from "./components/UpdateMonitorTriggers";
import DeleteMonitorForm from "./components/DeleteMonitorForm";
import UpdateMonitorFilters from "./components/UpdateMonitorFilters";

type Props = {
  team_id: string;
  monitor: ModuleMonitor;
};

const MonitorSettings: React.FC<Props> = ({ team_id, monitor }) => {
  const monitor_id = monitor.id;
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [kind, setKind] = useState<ModuleMonitorKind>("");
  const [schedule, setSchedule] = useState("");
  const [modules, setModules] = useState<string[]>();
  const [err, setErr] = useState("");

  const { refetch } = useQuery({
    queryKey: ["monitor", team_id, monitor.id],
    queryFn: async () => {
      const res = await api.getMonitor(team_id, monitor.id);
      return res;
    },
    retry: false,
  });

  const nameModified = useIsModified(name, monitor.name);
  const modulesModified = useIsModified(modules, monitor.modules);
  const kindModified = useIsModified(kind, monitor.kind);
  const scheduleModified = useIsModified(schedule, monitor.cron_schedule);
  const descriptionModified = useIsModified(description, monitor.description);

  const request = useMemo<UpdateMonitorRequest>(() => {
    let req: UpdateMonitorRequest = {};

    if (nameModified.isModified && name != "") {
      req.name = name;
    }

    if (kindModified.isModified && kind != "") {
      req.kind = kind;
    }

    if (descriptionModified.isModified && description != "") {
      req.description = description;
    }

    if (scheduleModified.isModified && schedule != "") {
      req.cron_schedule = schedule;
    }

    if (modulesModified.isModified && modules) {
      req.modules = modules;
    }

    return req;
  }, [
    name,
    nameModified,
    description,
    descriptionModified,
    kind,
    kindModified,
    schedule,
    scheduleModified,
    modules,
    modulesModified,
  ]);

  const mutation = useMutation({
    mutationKey: ["update_monitor", team_id, monitor_id],
    mutationFn: async (request: UpdateMonitorRequest) => {
      const res = await api.updateMonitor(team_id, monitor_id, request);
      return res;
    },
    onSuccess: (data) => {
      setErr("");

      refetch();
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  if (monitor.is_default) {
    return (
      <Placeholder>
        <SmallSpan>
          Cannot update settings for default monitors. Contact your Hatchet
          administrator to update default monitors.
        </SmallSpan>
      </Placeholder>
    );
  }

  return (
    <MonitorSettingsContainer>
      <SectionArea>
        <H4>Configuration</H4>
        <HorizontalSpacer spacepixels={20} />
        <UpdateMonitorMeta
          monitor={monitor}
          setMonitorName={setName}
          setMonitorDescription={setDescription}
        />
        <HorizontalSpacer spacepixels={24} />
        <ExpandableSettings text="Configure triggers">
          <UpdateMonitorTriggers
            monitor={monitor}
            setMonitorKind={setKind}
            setMonitorSchedule={setSchedule}
          />
        </ExpandableSettings>
        <HorizontalSpacer spacepixels={16} />
        <ExpandableSettings text="Configure filters">
          <UpdateMonitorFilters
            team_id={team_id}
            monitor={monitor}
            setMonitorModules={setModules}
          />
        </ExpandableSettings>
        <HorizontalSpacer spacepixels={24} />
        <FlexRowRight>
          <StandardButton
            label="Update"
            material_icon="chevron_right"
            icon_side="right"
            on_click={() => {
              mutation.mutate(request);
            }}
            margin={"0"}
            is_loading={mutation.isLoading}
          />
        </FlexRowRight>
      </SectionArea>
      <HorizontalSpacer spacepixels={20} />
      <SectionArea>
        <H4>Delete Monitor</H4>
        <HorizontalSpacer spacepixels={20} />
        <DeleteMonitorForm team_id={team_id} monitor={monitor} />
      </SectionArea>
    </MonitorSettingsContainer>
  );
};

export default MonitorSettings;
