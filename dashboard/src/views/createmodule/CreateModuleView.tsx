import {
  FlexRowRight,
  H1,
  HorizontalSpacer,
  Breadcrumbs,
  StandardButton,
} from "@hatchet-dev/hatchet-components";
import React, { useEffect, useState } from "react";
import { useHistory, useLocation, useParams } from "react-router-dom";
import ChooseGitSource from "./components/ChooseGitSource";
import SetupRuntime from "./components/SetupRuntime";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";
import { CreateModuleRequest } from "shared/api/generated/data-contracts";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";

const CreateModuleView: React.FunctionComponent = () => {
  const history = useHistory();
  const location = useLocation();
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);
  const [request, setRequest] = useState<CreateModuleRequest>(null);
  const [submittedStepOne, setSubmittedStepOne] = useState(false);
  const [err, setErr] = useState("");

  const { step } = useParams<{ step: string }>();

  const mutation = useMutation({
    mutationKey: ["create_module", currTeam?.id],
    mutationFn: async (req: CreateModuleRequest) => {
      const res = await api.createModule(currTeam.id, req);
      return res;
    },
    onSuccess: (data) => {
      setErr("");
      history.push(`/team/${currTeam.id}/modules`);
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  useEffect(() => {
    if (location.pathname.includes("step_2") && !submittedStepOne) {
      history.push(`/team/${currTeam.id}/modules/create/step_1`);
    }
  }, [submittedStepOne, location.pathname]);

  const submitStepOne = (req: CreateModuleRequest) => {
    setSubmittedStepOne(true);
    setRequest(req);
    history.push(`/team/${currTeam.id}/modules/create/step_2`);
  };

  const submitStepTwo = (req: CreateModuleRequest) => {
    mutation.mutate(req);
  };

  switch (step) {
    case "step_1":
      return <ChooseGitSource submit={submitStepOne} />;
    case "step_2":
      return <SetupRuntime req={request} submit={submitStepTwo} err={err} />;
  }
};

export default CreateModuleView;
