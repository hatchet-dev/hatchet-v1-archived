import {
  H2,
  HorizontalSpacer,
  P,
  TextInput,
  Selector,
  Selection,
  SectionArea,
  FlexRowRight,
  StandardButton,
  H1,
  Breadcrumbs,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { css } from "styled-components";
import theme from "shared/theme";
import { CreateModuleRequest } from "shared/api/generated/data-contracts";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";

const variableOptions = [
  {
    label: "Filesystem",
    value: "filesystem",
    material_icon: "folder_open",
  },
  {
    label: "Manual import",
    value: "manual",
    material_icon: "data_object",
  },
  {
    label: "Environment variables",
    value: "env",
    material_icon: "input",
  },
];

type Props = {
  req: CreateModuleRequest;
  submit: (req: CreateModuleRequest) => void;
};

const LinkVariables: React.FC<Props> = ({ req, submit }) => {
  const [varOption, setSelectedVarOption] = useState<string>();
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);

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
      label: "Step 2: Choose Variable Source",
      link: `/team/${currTeam.id}/modules/create/step_2`,
    },
  ];

  const selectVariableOption = (option: Selection) => {
    setSelectedVarOption(option.value);
  };

  const renderFilesystemOptions = () => {
    return [
      <HorizontalSpacer spacepixels={24} />,
      <P>
        Add the path to your tfvars file, relative to the root folder of the
        Github repository.
      </P>,
      <HorizontalSpacer spacepixels={24} />,
      <TextInput placeholder="./env1.tfvars" />,
    ];
  };

  const renderAdditionalFormOptions = () => {
    switch (varOption) {
      case "filesystem":
        return renderFilesystemOptions();
      default:
        return [];
    }
  };

  return (
    <>
      <Breadcrumbs breadcrumbs={breadcrumbs} />
      <HorizontalSpacer spacepixels={12} />
      <H1>Create a new module</H1>
      <HorizontalSpacer spacepixels={20} />
      <SectionArea>
        <H2>Step 2: Link Variables</H2>
        <HorizontalSpacer
          spacepixels={14}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={16} />
        <P>
          Choose how you would like to link your Terraform variables. You can
          link variables from your Git repository, via a variable file, or set
          them in your CI pipeline as environment variables.
        </P>
        <HorizontalSpacer spacepixels={24} />
        <Selector
          placeholder="Variable Source"
          placeholder_material_icon="edit_note"
          options={variableOptions}
          select={selectVariableOption}
        />
        {renderAdditionalFormOptions()}
      </SectionArea>
      <HorizontalSpacer spacepixels={20} />
      <FlexRowRight>
        <StandardButton
          label="Submit"
          material_icon="chevron_right"
          icon_side="right"
          on_click={() => {
            submit(req);
          }}
        />
      </FlexRowRight>
    </>
  );
};

export default LinkVariables;
