import {
  FlexColScroll,
  FlexRowLeft,
  TextInput,
  MaterialIcon,
  StandardButton,
  HorizontalSpacer,
} from "@hatchet-dev/hatchet-components";
import React from "react";
import { EnvVarRow, EnvVarRemoveButton } from "./styles";

type Props = {
  envVars: string[];
  setEnvVars: (vars: string[]) => void;
};

const EnvVars: React.FC<Props> = ({ envVars, setEnvVars }) => {
  const renderEnvVars = () => {
    return envVars.map((envVar, i) => {
      const splArr = envVar.split("~~=~~");
      if (splArr.length == 2) {
        const rowVal = (
          <EnvVarRow key={i}>
            <TextInput
              placeholder="KEY"
              initial_value={splArr[0]}
              on_change={(newKey) => updateEnvVarKey(i, newKey)}
            />
            <TextInput
              placeholder="VALUE"
              initial_value={splArr[1]}
              on_change={(newVal) => updateEnvVarVal(i, newVal)}
            />
            <EnvVarRemoveButton onClick={() => removeEnvVar(i)}>
              <MaterialIcon className="material-icons">close</MaterialIcon>
            </EnvVarRemoveButton>
          </EnvVarRow>
        );

        if (i == envVars.length - 1) {
          return rowVal;
        }

        return (
          <>
            {rowVal}
            <HorizontalSpacer spacepixels={10} />
          </>
        );
      }
    });
  };

  const addNewEnvVar = () => {
    const newEnvVars = [...envVars, "~~=~~"];

    setEnvVars(newEnvVars);
  };

  const updateEnvVarKey = (index: number, key: string) => {
    const newEnvVars = [...envVars];
    const strArr = envVars[index].split("~~=~~");
    const oldEnvVarVal = strArr[1];

    newEnvVars[index] = `${key}~~=~~${oldEnvVarVal}`;

    setEnvVars(newEnvVars);
  };

  const updateEnvVarVal = (index: number, val: string) => {
    const newEnvVars = [...envVars];
    const strArr = envVars[index].split("~~=~~");
    const oldEnvVarKey = strArr[0];

    newEnvVars[index] = `${oldEnvVarKey}~~=~~${val}`;

    setEnvVars(newEnvVars);
  };

  const removeEnvVar = (index: number) => {
    let newEnvVars = [...envVars];

    newEnvVars.splice(index, 1);

    setEnvVars(newEnvVars);
  };

  return (
    <FlexColScroll>
      {renderEnvVars()}
      <HorizontalSpacer spacepixels={10} />
      <FlexRowLeft>
        <StandardButton
          material_icon="add"
          style_kind="muted"
          label="Add environment variable"
          margin="0"
          on_click={addNewEnvVar}
        />
      </FlexRowLeft>
    </FlexColScroll>
  );
};

export default EnvVars;
