import { H2, HorizontalSpacer, P } from "components/globals";
import Selector, { Selection } from "components/selector";
import React, { useState } from "react";
import TextInput from "components/textinput";
import SectionArea from "components/sectionarea";
import { css } from "styled-components";
import theme from "shared/theme";

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

const LinkVariables: React.FunctionComponent = () => {
  const [varOption, setSelectedVarOption] = useState<string>();

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
        Choose how you would like to link your Terraform variables. You can link
        variables from your Git repository, via a variable file, or set them in
        your CI pipeline as environment variables.
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
  );
};

export default LinkVariables;
