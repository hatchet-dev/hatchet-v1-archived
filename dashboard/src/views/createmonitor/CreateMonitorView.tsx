import {
  FlexRowRight,
  H1,
  HorizontalSpacer,
  Breadcrumbs,
  StandardButton,
  TextInput,
  SectionArea,
  H2,
  P,
} from "@hatchet-dev/hatchet-components";
import React, { useEffect, useState } from "react";
import { useHistory, useLocation, useParams } from "react-router-dom";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";
import {
  CreateModuleRequest,
  CreateMonitorRequest,
} from "shared/api/generated/data-contracts";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";
import theme from "shared/theme";
import { css } from "styled-components";

const CreateMonitorView: React.FunctionComponent = () => {
  const history = useHistory();
  const location = useLocation();
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);
  const [request, setRequest] = useState<CreateModuleRequest>(null);
  const [name, setName] = useState("");
  const [schedule, setSchedule] = useState("");
  const [submittedStepOne, setSubmittedStepOne] = useState(false);
  const [err, setErr] = useState("");

  const { step } = useParams<{ step: string }>();

  const mutation = useMutation({
    mutationKey: ["create_monitor", currTeam?.id],
    mutationFn: (req: CreateMonitorRequest) => {
      return api.createMonitor(currTeam.id, req);
    },
    onSuccess: (data) => {
      setErr("");
      history.push(`/team/${currTeam.id}/monitors`);
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  return (
    <>
      <HorizontalSpacer spacepixels={12} />
      <H1>Create a new monitor</H1>
      <HorizontalSpacer spacepixels={20} />
      <SectionArea>
        <H2>Configure Monitor Triggers</H2>
        <HorizontalSpacer
          spacepixels={14}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={16} />
        <P>Give the monitor a name.</P>
        <HorizontalSpacer spacepixels={12} />
        <TextInput
          placeholder="ex. my-monitor"
          on_change={(val) => {
            setName(val);
          }}
        />
        <HorizontalSpacer spacepixels={16} />
        <TextInput
          label="Cron Schedule"
          placeholder="ex. * * * * *"
          on_change={(val) => {
            setSchedule(val);
          }}
        />
      </SectionArea>
      <HorizontalSpacer spacepixels={20} />
      <FlexRowRight>
        <StandardButton
          label="Next"
          material_icon="chevron_right"
          icon_side="right"
          on_click={() => {
            mutation.mutate({
              cron_schedule: schedule,
            });
          }}
        />
      </FlexRowRight>
    </>
  );
};

export default CreateMonitorView;
