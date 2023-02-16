import {
  FlexColScroll,
  FlexRowLeft,
  TextInput,
  MaterialIcon,
  StandardButton,
  HorizontalSpacer,
} from "@hatchet-dev/hatchet-components";
import { atom, useAtom, createStore, useStore, PrimitiveAtom } from "jotai";
import React, { useState, useEffect } from "react";
import { makeid } from "shared/utils";
import { EnvVarRow, EnvVarRemoveButton } from "./styles";

export const newEnvVarAtom = (envVars: string[]) => {
  const envVarAtom = atom<InternalEnvVar[]>(getInternalEnvVars(envVars));

  return envVarAtom;
};

export const getInternalEnvVars = (envVars: string[]) => {
  return envVars.map((envVar) => {
    return {
      value: envVar,
      key: makeid(16),
    };
  });
};

type Props = {
  envVarAtom: PrimitiveAtom<InternalEnvVar[]>;
  read_only?: boolean;
};

type InternalEnvVar = {
  value: string;
  key: string;
};

const EnvVars: React.FC<Props> = ({ envVarAtom, read_only = false }) => {
  const [envVars, setEnvVars] = useAtom(envVarAtom);

  const renderEnvVars = () => {
    return envVars.map((envVar, i) => {
      const splArr = envVar.value.split("~~=~~");
      if (splArr.length == 2) {
        return (
          <React.Fragment key={envVar.key}>
            <EnvVarRow>
              <TextInput
                placeholder="KEY"
                initial_value={splArr[0]}
                on_change={(newKey) => updateEnvVarKey(i, newKey)}
                disabled={read_only}
              />
              <TextInput
                placeholder="VALUE"
                initial_value={splArr[1]}
                on_change={(newVal) => updateEnvVarVal(i, newVal)}
                disabled={read_only}
              />
              {!read_only && (
                <EnvVarRemoveButton onClick={() => removeEnvVar(i)}>
                  <MaterialIcon className="material-icons">close</MaterialIcon>
                </EnvVarRemoveButton>
              )}
            </EnvVarRow>
            {i != envVars.length - 1 && <HorizontalSpacer spacepixels={10} />}
          </React.Fragment>
        );
      }
    });
  };

  const addNewEnvVar = () => {
    const newEnvVars = envVars.map((envVar) => envVar);
    newEnvVars.push({
      value: "~~=~~",
      key: makeid(16),
    });

    setEnvVars(newEnvVars);
  };

  const updateEnvVarKey = (index: number, key: string) => {
    const newEnvVars = envVars.map((envVar) => envVar);
    const strArr = envVars[index].value.split("~~=~~");
    const oldEnvVarVal = strArr[1];

    newEnvVars[index].value = `${key}~~=~~${oldEnvVarVal}`;

    setEnvVars(newEnvVars);
  };

  const updateEnvVarVal = (index: number, val: string) => {
    const newEnvVars = envVars.map((envVar) => envVar);
    const strArr = envVars[index].value.split("~~=~~");
    const oldEnvVarKey = strArr[0];

    newEnvVars[index].value = `${oldEnvVarKey}~~=~~${val}`;

    setEnvVars(newEnvVars);
  };

  const removeEnvVar = (index: number) => {
    const newEnvVars = envVars.map((envVar) => envVar);

    newEnvVars.splice(index, 1);

    setEnvVars(newEnvVars);
  };

  return (
    <FlexColScroll>
      {renderEnvVars()}
      {!read_only && <HorizontalSpacer spacepixels={10} />}
      {!read_only && (
        <FlexRowLeft>
          <StandardButton
            material_icon="add"
            style_kind="muted"
            label="Add environment variable"
            margin="0"
            on_click={addNewEnvVar}
          />
        </FlexRowLeft>
      )}
    </FlexColScroll>
  );
};

export default EnvVars;
