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
  FlexColScroll,
  ErrorBar,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { css } from "styled-components";
import theme from "shared/theme";
import { CreateModuleRequest } from "shared/api/generated/data-contracts";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";
import CodeBlock from "components/codeblock";

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
];

type Props = {
  req: CreateModuleRequest;
  submit: (req: CreateModuleRequest) => void;
  err?: string;
};

const SetupRuntime: React.FC<Props> = ({ req, submit, err }) => {
  const [varOption, setSelectedVarOption] = useState<string>();
  const [filePath, setFilePath] = useState("");
  const [jsonValues, setJSONValues] = useState("{\n  \n}");
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);
  const [jsonParseErr, setJSONParseErr] = useState("");

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
      <TextInput
        placeholder="./env1.tfvars"
        on_change={(p) => setFilePath(p)}
      />,
    ];
  };

  const renderManualImportOptions = () => {
    return [
      <HorizontalSpacer spacepixels={24} />,
      <P>Upload your JSON variables here.</P>,
      <HorizontalSpacer spacepixels={24} />,
      <FlexColScroll height="200px" width="100%">
        <CodeBlock
          value={jsonValues}
          height="200px"
          onChange={(e) => setJSONValues(e)}
        />
      </FlexColScroll>,
      jsonParseErr && <HorizontalSpacer spacepixels={20} />,
      jsonParseErr && <ErrorBar text={jsonParseErr} />,
    ];
  };

  const renderAdditionalFormOptions = () => {
    switch (varOption) {
      case "filesystem":
        return renderFilesystemOptions();
      case "manual":
        return renderManualImportOptions();
      default:
        return [];
    }
  };

  const onSubmit = () => {
    switch (varOption) {
      case "filesystem":
        req.values_github = {
          path: filePath,
          github_app_installation_id: req.github.github_app_installation_id,
          github_repository_branch: req.github.github_repository_branch,
          github_repository_name: req.github.github_repository_name,
          github_repository_owner: req.github.github_repository_owner,
        };

        submit(req);
        break;
      case "manual":
        try {
          const values = JSON.parse(jsonValues);

          req.values_raw = values;

          submit(req);
        } catch (e) {
          setJSONParseErr("Could not parse JSON");
        }

        break;
      default:
    }
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
        <HorizontalSpacer spacepixels={24} />
        <P>
          Add any additional environment variables (for example, credentials or
          Terraform variables) below.
        </P>
        <HorizontalSpacer spacepixels={24} />
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
          disabled={!varOption}
        />
      </FlexRowRight>
    </>
  );
};

export default SetupRuntime;
