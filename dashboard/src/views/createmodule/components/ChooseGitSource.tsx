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
} from "@hatchet-dev/hatchet-components";
import React, { useMemo, useState } from "react";
import { css } from "styled-components";
import theme from "shared/theme";
import {
  CreateModuleRequest,
  CreateModuleRequestGithub,
} from "shared/api/generated/data-contracts";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";
import SelectGitSource from "components/module/selectgitpath";

type Props = {
  submit: (req: CreateModuleRequest) => void;
};

const ChooseGitSource: React.FunctionComponent<Props> = ({ submit }) => {
  const [name, setName] = useState("");

  const [githubParams, setGithubParams] = useState<CreateModuleRequestGithub>();

  const request = useMemo<CreateModuleRequest>(() => {
    return {
      name: name,
      github: githubParams,
    };
  }, [githubParams, name]);

  const [currTeam] = useAtom(currTeamAtom);

  const breadcrumbs = [
    {
      label: "Modules",
      link: `/team/${currTeam.id}/modules`,
    },
    {
      label: "Step 1: Choose Git Source",
      link: `/team/${currTeam.id}/modules/create/step_1`,
    },
  ];

  const submitEnabled =
    !!request.name &&
    !!request.github?.github_app_installation_id &&
    !!request.github?.github_repository_owner &&
    !!request.github?.github_repository_name &&
    !!request.github?.github_repository_branch &&
    !!request.github?.path;

  return (
    <>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
      <HorizontalSpacer spacepixels={12} />
      <H1>Create a new module</H1>
      <HorizontalSpacer spacepixels={20} />
      <SectionArea>
        <H2>Step 1: Name and Source Configuration</H2>
        <HorizontalSpacer
          spacepixels={14}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={16} />
        <P>Give the module a name.</P>
        <HorizontalSpacer spacepixels={12} />
        <TextInput
          placeholder="ex. my-module"
          on_change={(val) => {
            setName(val);
          }}
        />
        <HorizontalSpacer spacepixels={16} />
        <SelectGitSource set_request={setGithubParams} />
      </SectionArea>
      <HorizontalSpacer spacepixels={20} />
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

export default ChooseGitSource;
