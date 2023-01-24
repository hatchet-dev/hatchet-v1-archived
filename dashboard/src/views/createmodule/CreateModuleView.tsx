import {
  FlexRowRight,
  H1,
  HorizontalSpacer,
  Breadcrumbs,
  StandardButton,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useHistory, useParams } from "react-router-dom";
import ChooseGitSource from "./components/ChooseGitSource";
import LinkVariables from "./components/LinkVariables";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";
import { CreateModuleRequest } from "shared/api/generated/data-contracts";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";

const CreateModuleView: React.FunctionComponent = () => {
  const history = useHistory();
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);
  const [request, setRequest] = useState<CreateModuleRequest>(null);
  const [err, setErr] = useState("");

  const { step } = useParams<{ step: string }>();

  const mutation = useMutation({
    mutationKey: ["create_organization_invite", currTeam?.id],
    mutationFn: (req: CreateModuleRequest) => {
      return api.createModule(currTeam.id, req);
    },
    onSuccess: (data) => {
      setErr("");
      history.push(`/team/${currTeam.id}/modules`);
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  const submitStepOne = (req: CreateModuleRequest) => {
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
      return <LinkVariables req={request} submit={submitStepTwo} />;
  }
};

export default CreateModuleView;
