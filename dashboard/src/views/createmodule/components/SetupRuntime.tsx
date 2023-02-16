import {
  H2,
  HorizontalSpacer,
  P,
  SectionArea,
  FlexRowRight,
  StandardButton,
  H1,
  Breadcrumbs,
  ErrorBar,
} from "@hatchet-dev/hatchet-components";
import React, { useMemo, useState } from "react";
import { css } from "styled-components";
import theme from "shared/theme";
import {
  CreateModuleRequest,
  CreateModuleValuesRequestGithub,
} from "shared/api/generated/data-contracts";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";
import EnvVars, { newEnvVarAtom } from "components/envvars";
import SetModuleValues from "components/module/setmodulevalues";

type Props = {
  req: CreateModuleRequest;
  submit: (req: CreateModuleRequest) => void;
  err?: string;
};

const SetupRuntime: React.FC<Props> = ({ req, submit, err }) => {
  const [jsonValues, setJSONValues] = useState<Record<string, object>>();
  const [
    githubValues,
    setGithubValues,
  ] = useState<CreateModuleValuesRequestGithub>({
    github_app_installation_id: req?.github.github_app_installation_id,
    github_repository_branch: req?.github.github_repository_branch,
    github_repository_name: req?.github.github_repository_name,
    github_repository_owner: req?.github.github_repository_owner,
    path: "",
  });
  const [valuesSource, setValuesSource] = useState<string>();
  const [currTeam] = useAtom(currTeamAtom);
  const envVarAtom = useMemo(() => {
    return newEnvVarAtom([]);
  }, []);

  const [envVars, setEnvVars] = useAtom(envVarAtom);

  const breadcrumbs = [
    {
      label: "Modules",
      link: `/team/${currTeam.id}/modules`,
    },
    {
      label: "Step 1: Choose Git Source",
      link: `/team/${currTeam.id}/modules/create/step_1`,
    },
    {
      label: "Step 2: Configure Runtime",
      link: `/team/${currTeam.id}/modules/create/step_2`,
    },
  ];

  const onSubmit = () => {
    if (valuesSource == "github") {
      req.values_github = githubValues;
    } else {
      req.values_raw = jsonValues;
    }

    let mappedEnvVars: Record<string, string> = {};

    envVars.forEach((envVar) => {
      const strArr = envVar.value.split("~~=~~");
      if (strArr.length == 2) {
        mappedEnvVars[strArr[0]] = strArr[1];
      }
    });

    req.env_vars = mappedEnvVars;

    submit(req);
  };

  return (
    <>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
      <HorizontalSpacer spacepixels={12} />
      <H1>Create a new module</H1>
      <HorizontalSpacer spacepixels={20} />
      <SectionArea>
        <H2>Step 2: Configure Runtime Environment</H2>
        <HorizontalSpacer
          spacepixels={14}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={16} />
        <SetModuleValues
          set_github_values={setGithubValues}
          current_github_params={githubValues}
          set_raw_values={setJSONValues}
          set_values_source={setValuesSource}
        />
        <P>
          Add any additional environment variables (for example, credentials or
          Terraform variables) below.
        </P>
        <HorizontalSpacer spacepixels={24} />
        <EnvVars envVarAtom={envVarAtom} />
      </SectionArea>
      <HorizontalSpacer spacepixels={20} />
      {err && <ErrorBar text={err} />}
      {err && <HorizontalSpacer spacepixels={20} />}
      <FlexRowRight>
        <StandardButton
          label="Submit"
          material_icon="chevron_right"
          icon_side="right"
          on_click={onSubmit}
          disabled={!valuesSource}
        />
      </FlexRowRight>
    </>
  );
};

export default SetupRuntime;
