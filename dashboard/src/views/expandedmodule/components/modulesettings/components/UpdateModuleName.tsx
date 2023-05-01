import {
  StandardButton,
  FlexRowRight,
  FlexCol,
  HorizontalSpacer,
  SmallSpan,
  TextInput,
} from "hatchet-components";
import { useMutation } from "@tanstack/react-query";
import React, { useState } from "react";
import api from "shared/api";
import {
  Module,
  ModuleRun,
  UpdateModuleRequest,
} from "shared/api/generated/data-contracts";

type Props = {
  module: Module;
  setModuleName: (name: string) => void;
};

const UpdateModuleName: React.FC<Props> = ({ module, setModuleName }) => {
  return (
    <TextInput
      placeholder="ex. my-module"
      initial_value={module.name}
      label="Module name"
      type="text"
      width="400px"
      on_change={(val) => {
        setModuleName(val);
      }}
    />
  );
};

export default UpdateModuleName;
