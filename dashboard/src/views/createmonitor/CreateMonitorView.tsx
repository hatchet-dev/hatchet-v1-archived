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
import { CreateMonitorRequest } from "shared/api/generated/data-contracts";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";
import theme from "shared/theme";
import { css } from "styled-components";
import MonitorMeta from "./components/MonitorMeta";
import SetupPolicy from "./components/SetupPolicy";

const CreateMonitorView: React.FunctionComponent = () => {
  const history = useHistory();
  const location = useLocation();
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);
  const [request, setRequest] = useState<CreateMonitorRequest>(null);
  const [name, setName] = useState("");
  const [schedule, setSchedule] = useState("");
  const [submittedStepOne, setSubmittedStepOne] = useState(false);
  const [err, setErr] = useState("");

  const { step } = useParams<{ step: string }>();

  const mutation = useMutation({
    mutationKey: ["create_monitor", currTeam?.id],
    mutationFn: async (req: CreateMonitorRequest) => {
      const res = await api.createMonitor(currTeam.id, req);
      return res;
    },
    onSuccess: (data) => {
      setErr("");
      history.push(`/team/${currTeam.id}/monitors`);
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  const submitStepOne = (req: CreateMonitorRequest) => {
    setSubmittedStepOne(true);
    setRequest(req);
    history.push(`/team/${currTeam.id}/monitors/create/step_2`);
  };

  const submitStepTwo = (req: CreateMonitorRequest) => {
    mutation.mutate(req);
  };

  switch (step) {
    case "step_1":
      return <MonitorMeta submit={submitStepOne} />;
    case "step_2":
      return <SetupPolicy req={request} submit={submitStepTwo} err={err} />;
  }
};

export default CreateMonitorView;
