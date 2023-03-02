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
  CreateMonitorRequest,
} from "shared/api/generated/data-contracts";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";
import EnvVars, { newEnvVarAtom } from "components/envvars";
import SetModuleValues from "components/module/setmodulevalues";
import CodeBlock from "components/codeblock";

type Props = {
  req: CreateMonitorRequest;
  submit: (req: CreateMonitorRequest) => void;
  err?: string;
};

const SetupPolicy: React.FC<Props> = ({ req, submit, err }) => {
  const [policyBytes, setPolicyBytes] = useState("");
  const [currTeam] = useAtom(currTeamAtom);

  const breadcrumbs = [
    {
      label: "Monitors",
      link: `/teams/${currTeam.id}/monitors`,
    },
    {
      label: "Step 1: Monitor Metadata",
      link: `/teams/${currTeam.id}/monitors/create/step_1`,
    },
    {
      label: "Step 2: Configure Policy",
      link: `/teams/${currTeam.id}/monitors/create/step_2`,
    },
  ];

  const onSubmit = () => {
    req.policy_bytes = policyBytes;

    submit(req);
  };

  return (
    <>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
      <HorizontalSpacer spacepixels={12} />
      <H1>Create a new monitor</H1>
      <HorizontalSpacer spacepixels={20} />
      <SectionArea>
        <H2>Step 2: Configure Policy</H2>
        <HorizontalSpacer
          spacepixels={14}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={16} />
        <CodeBlock value={policyBytes} onChange={(v) => setPolicyBytes(v)} />
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
          disabled={!policyBytes}
        />
      </FlexRowRight>
    </>
  );
};

export default SetupPolicy;
