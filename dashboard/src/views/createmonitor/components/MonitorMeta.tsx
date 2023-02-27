import {
  H2,
  HorizontalSpacer,
  P,
  TextInput,
  SectionArea,
  FlexRowRight,
  StandardButton,
  Breadcrumbs,
  H1,
  Selector,
  Selection,
} from "@hatchet-dev/hatchet-components";
import React, { useMemo, useState } from "react";
import { css } from "styled-components";
import theme from "shared/theme";
import { CreateMonitorRequest } from "shared/api/generated/data-contracts";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";
import SelectGitSource from "components/module/selectgitpath";
import { useQuery } from "@tanstack/react-query";
import api from "shared/api";

type Props = {
  submit: (req: CreateMonitorRequest) => void;
};

export const MonitorKindOptions = [
  {
    label: "Scheduled Plan Check",
    value: "plan",
    material_icon: "schedule",
  },
  {
    label: "Scheduled State Check",
    value: "state",
    material_icon: "schedule",
  },
];

const MonitorMeta: React.FunctionComponent<Props> = ({ submit }) => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [kind, setKind] = useState("plan");
  const [schedule, setSchedule] = useState("");
  const [moduleOptions, setModuleOptions] = useState<Selection[]>([]);

  const [currTeam] = useAtom(currTeamAtom);

  const request = useMemo<CreateMonitorRequest>(() => {
    return {
      name: name,
      description: description,
      kind: kind,
      cron_schedule: schedule,
      modules: moduleOptions?.map((option) => option.value) || [],
    };
  }, [name, description, kind, schedule, moduleOptions]);

  useQuery({
    queryKey: ["modules", currTeam?.id],
    queryFn: async () => {
      const res = await api.listModules(currTeam?.id);

      return res;
    },
    retry: false,
    onSuccess: (data) => {
      const modOptions: Selection[] = data.data.rows.map((module) => {
        return {
          label: module.name,
          value: module.id,
          checked: false,
        };
      });

      setModuleOptions(modOptions);
    },
    refetchOnWindowFocus: false,
  });

  const breadcrumbs = [
    {
      label: "Monitors",
      link: `/team/${currTeam.id}/monitors`,
    },
    {
      label: "Step 1: Monitor Metadata",
      link: `/team/${currTeam.id}/monitors/create/step_1`,
    },
  ];

  const selectKind = (option: Selection) => {
    setKind(option.value);
  };

  const submitEnabled =
    !!request.name &&
    !!request.description &&
    !!request.kind &&
    !!request.cron_schedule;

  return (
    <>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
      <HorizontalSpacer spacepixels={12} />
      <H1>Create a new monitor</H1>
      <HorizontalSpacer spacepixels={20} />
      <SectionArea>
        <H2>Step 1: Name and Trigger Configuration</H2>
        <HorizontalSpacer
          spacepixels={14}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={20} />
        <P>Give the monitor a name.</P>
        <HorizontalSpacer spacepixels={12} />
        <TextInput
          placeholder="ex. Drift detection"
          on_change={(val) => {
            setName(val);
          }}
        />
        <HorizontalSpacer spacepixels={20} />
        <P>Give the monitor a description.</P>
        <HorizontalSpacer spacepixels={12} />
        <TextInput
          width="600px"
          placeholder="ex. Detects drift"
          on_change={(val) => {
            setDescription(val);
          }}
        />
        <HorizontalSpacer spacepixels={20} />
        <P>
          Choose when this policy check should be run. You can configure this to
          run periodically against a Terraform plan or the Terraform state, or
          you can run checks before/after Terraform operations.
        </P>
        <HorizontalSpacer spacepixels={12} />
        <Selector
          placeholder="Scheduled Plan Check"
          placeholder_material_icon="schedule"
          options={MonitorKindOptions}
          select={selectKind}
        />
        <HorizontalSpacer spacepixels={20} />
        <P>Provide a cron schedule to run these policy checks.</P>
        <HorizontalSpacer spacepixels={12} />
        <TextInput
          placeholder="ex. * * * * *"
          on_change={(val) => {
            setSchedule(val);
          }}
        />
        <HorizontalSpacer spacepixels={20} />
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
          }}
        />
      </SectionArea>
      <HorizontalSpacer spacepixels={30} />
      <FlexRowRight>
        <StandardButton
          label="Next"
          material_icon="chevron_right"
          icon_side="right"
          disabled={!submitEnabled}
          on_click={() => {
            if (!submitEnabled) {
              return;
            }

            submit(request);
          }}
        />
      </FlexRowRight>
    </>
  );
};

export default MonitorMeta;
