import {
  HorizontalSpacer,
  P,
  TextInput,
  Selector,
  Selection,
} from "hatchet-components";
import React, { useState, useEffect } from "react";
import { CreateModuleRequestGithub } from "shared/api/generated/data-contracts";
import { useAtom } from "jotai";
import { currTeamAtom } from "shared/atoms/atoms";
import UploadJSONValues from "components/module/uploadjsonvalues";

const variableOptions = [
  {
    label: "Filesystem",
    value: "github",
    material_icon: "folder_open",
  },
  {
    label: "Manual import",
    value: "raw",
    material_icon: "data_object",
  },
];

type Props = {
  set_github_values: (values: CreateModuleRequestGithub) => void;
  current_github_params?: CreateModuleRequestGithub;
  set_raw_values: (values: Record<string, object>) => void;
  current_raw_values?: Record<string, object>;
  set_values_source: (source: string) => void;
  current_values_source?: string;
};

const SetModuleValues: React.FC<Props> = ({
  set_github_values,
  current_github_params,
  set_raw_values,
  current_raw_values,
  set_values_source,
  current_values_source,
}) => {
  const [varOption, setSelectedVarOption] = useState<string>(
    current_values_source
  );
  const [filePath, setFilePath] = useState("");
  const [jsonValues, setJSONValues] = useState(
    current_raw_values ? JSON.stringify(current_raw_values) : "{\n  \n}"
  );
  const [jsonParseErr, setJSONParseErr] = useState("");

  useEffect(() => {
    let curr = current_github_params;
    curr.path = filePath;
    set_github_values(curr);
  }, [filePath]);

  useEffect(() => {
    set_values_source(varOption);
  }, [varOption]);

  useEffect(() => {
    try {
      const values = JSON.parse(jsonValues);

      setJSONParseErr("");
      set_raw_values(values);
    } catch (e) {
      setJSONParseErr("Errors while parsing JSON");
      return;
    }
  }, [jsonValues]);

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
        initial_value={current_github_params.path}
      />,
    ];
  };

  const renderAdditionalFormOptions = () => {
    switch (varOption) {
      case "github":
        return renderFilesystemOptions();
      case "raw":
        return (
          <UploadJSONValues
            jsonParseErr={jsonParseErr}
            set_values={setJSONValues}
            current_values={jsonValues}
          />
        );
      default:
        return [];
    }
  };

  const getSelectorPlaceholder = () => {
    switch (current_values_source) {
      case "github":
        return "Filesystem";
      case "raw":
        return "Manual import";
      default:
        return "Variable Source";
    }
  };

  const getSelectorPlaceholderIcon = () => {
    switch (current_values_source) {
      case "github":
        return "folder_open";
      case "raw":
        return "data_object";
      default:
        return "edit_note";
    }
  };

  return (
    <>
      <HorizontalSpacer spacepixels={16} />
      <P>
        Choose how you would like to link your Terraform variables. You can link
        variables from your Git repository or via a manual variable file upload.
      </P>
      <HorizontalSpacer spacepixels={24} />
      <Selector
        placeholder={getSelectorPlaceholder()}
        placeholder_material_icon={getSelectorPlaceholderIcon()}
        options={variableOptions}
        select={selectVariableOption}
      />
      {renderAdditionalFormOptions()}
      <HorizontalSpacer spacepixels={24} />
    </>
  );
};

export default SetModuleValues;
