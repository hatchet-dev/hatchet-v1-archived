import {
  StandardButton,
  FlexRowRight,
  HorizontalSpacer,
  H4,
  Placeholder,
  Spinner,
  SectionArea,
} from "hatchet-components";
import { useMutation, useQuery } from "@tanstack/react-query";
import EnvVars, { getInternalEnvVars, newEnvVarAtom } from "components/envvars";
import ExpandableSettings from "components/expandablesettings";
import SelectGitSource from "components/module/selectgitpath";
import SetModuleValues from "components/module/setmodulevalues";
import React, { useMemo, useState } from "react";
import api from "shared/api";
import {
  CreateModuleRequestGithub,
  CreateModuleValuesRequestGithub,
  Module,
  UpdateModuleRequest,
} from "shared/api/generated/data-contracts";
import useIsModified from "shared/hooks/useismodified";
import DeleteModuleForm from "./components/DeleteModuleForm";
import ModuleSettingsCard from "./components/ModuleSettingsCard";
import ModuleSettingsContainer from "./components/ModuleSettingsContainer";
import UnlockModuleForm from "./components/UnlockModuleForm";
import UpdateModuleName from "./components/UpdateModuleName";
import { useAtom } from "jotai";

type Props = {
  team_id: string;
  module: Module;
};

const ModuleSettings: React.FC<Props> = ({ team_id, module }) => {
  const module_id = module.id;
  const [name, setName] = useState("");
  const [err, setErr] = useState("");

  const [githubParams, setGithubParams] = useState<CreateModuleRequestGithub>();
  const [
    githubValueParams,
    setGithubValueParams,
  ] = useState<CreateModuleValuesRequestGithub>();
  const [rawValues, setRawValues] = useState<Record<string, object>>();
  const [valuesSource, setValuesSource] = useState<string>();

  const envVarAtom = useMemo(() => {
    return newEnvVarAtom([]);
  }, []);

  const [envVars, setEnvVars] = useAtom(envVarAtom);

  const { refetch } = useQuery({
    queryKey: ["module", team_id, module.id],
    queryFn: async () => {
      const res = await api.getModule(team_id, module.id);
      return res;
    },
    retry: false,
  });

  const envVarsQuery = useQuery({
    queryKey: [
      "module_env_vars",
      team_id,
      module_id,
      module.current_env_vars_version_id,
    ],
    queryFn: async () => {
      const res = await api.getModuleEnvVars(
        team_id,
        module_id,
        module.current_env_vars_version_id
      );

      return res;
    },
    onSuccess: (res) => {
      const newVars = getInternalEnvVars(
        res.data.env_vars.map((envVar) => {
          return `${envVar.key}~~=~~${envVar.val}`;
        })
      );

      setEnvVars(newVars);
      envVarsModified && envVarsModified.reset(newVars);
    },
    retry: false,
    refetchOnWindowFocus: false,
  });

  const valuesQuery = useQuery({
    queryKey: [
      "module_values",
      team_id,
      module_id,
      module.current_values_version_id,
    ],
    queryFn: async () => {
      const res = await api.getModuleValues(
        team_id,
        module_id,
        module.current_values_version_id
      );

      return res;
    },
    retry: false,
    refetchOnWindowFocus: false,
  });

  const valuesGH = valuesQuery.data?.data?.github;

  const currentGithubParams = {
    github_app_installation_id: module.deployment.github_app_installation_id,
    github_repository_owner: module.deployment.git_repo_owner,
    github_repository_name: module.deployment.git_repo_name,
    github_repository_branch: module.deployment.git_repo_branch,
    path: module.deployment.path,
  };

  const currentValuesGithubParams = valuesQuery.data && {
    github_app_installation_id:
      valuesGH?.github_app_installation_id ||
      module.deployment.github_app_installation_id,
    github_repository_branch:
      valuesGH?.github_repo_branch || module.deployment.git_repo_branch,
    github_repository_owner:
      valuesGH?.github_repo_owner || module.deployment.git_repo_owner,
    github_repository_name:
      valuesGH?.github_repo_name || module.deployment.git_repo_name,
    path: valuesGH?.path,
  };

  const currentValuesSource = valuesQuery.data && (valuesGH ? "github" : "raw");

  const githubParamsModified = useIsModified(githubParams, currentGithubParams);
  const githubValueParamsModified = useIsModified(
    githubValueParams,
    currentValuesGithubParams
  );
  const valuesSourceModified = useIsModified(valuesSource, currentValuesSource);
  const rawValuesParamsModified = useIsModified(rawValues);
  const envVarsModified = useIsModified(envVars);

  const request = useMemo<UpdateModuleRequest>(() => {
    let req: UpdateModuleRequest = {
      name: name,
    };

    if (githubParamsModified.isModified && githubParams) {
      req.github = githubParams;
    }

    if (
      (valuesSourceModified.isModified ||
        githubValueParamsModified.isModified) &&
      valuesSource == "github"
    ) {
      req.values_github = githubValueParams;
    }

    if (
      (valuesSourceModified.isModified || rawValuesParamsModified.isModified) &&
      valuesSource == "raw"
    ) {
      req.values_raw = rawValues;
    }

    if (envVarsModified.isModified) {
      let mappedEnvVars: Record<string, string> = {};

      envVars.forEach((envVar) => {
        const strArr = envVar.value.split("~~=~~");
        if (strArr.length == 2) {
          mappedEnvVars[strArr[0]] = strArr[1];
        }
      });

      req.env_vars = mappedEnvVars;
    }

    return req;
  }, [
    name,
    githubParams,
    githubParamsModified,
    envVars,
    envVarsModified,
    valuesSource,
    valuesSourceModified,
    githubValueParams,
    githubValueParamsModified,
    rawValues,
    rawValuesParamsModified,
  ]);

  const mutation = useMutation({
    mutationKey: ["update_module", team_id, module_id],
    mutationFn: async (request: UpdateModuleRequest) => {
      const res = await api.updateModule(team_id, module_id, request);
      return res;
    },
    onSuccess: (data) => {
      setErr("");

      // TODO: compute refetchable
      refetch();
      envVarsQuery.refetch();
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  if (envVarsQuery.isLoading || valuesQuery.isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  return (
    <ModuleSettingsContainer>
      <SectionArea>
        <H4>Configuration</H4>
        <HorizontalSpacer spacepixels={20} />
        <UpdateModuleName module={module} setModuleName={setName} />
        <HorizontalSpacer spacepixels={24} />
        {module?.deployment_mechanism == "github" && (
          <ExpandableSettings text="Github settings">
            <SelectGitSource
              set_request={setGithubParams}
              current_params={currentGithubParams}
            />
          </ExpandableSettings>
        )}
        {module?.deployment_mechanism == "github" && (
          <HorizontalSpacer spacepixels={8} />
        )}
        <ExpandableSettings text="Values configuration">
          <SetModuleValues
            set_github_values={setGithubValueParams}
            current_github_params={currentValuesGithubParams}
            set_raw_values={setRawValues}
            current_raw_values={valuesQuery.data?.data.raw_values}
            set_values_source={setValuesSource}
            current_values_source={currentValuesSource}
          />
        </ExpandableSettings>
        <HorizontalSpacer spacepixels={8} />
        <ExpandableSettings text="Environment variables">
          <EnvVars envVarAtom={envVarAtom} />
        </ExpandableSettings>
        <HorizontalSpacer spacepixels={24} />
        <FlexRowRight>
          <StandardButton
            label="Update"
            material_icon="chevron_right"
            icon_side="right"
            on_click={() => {
              mutation.mutate(request);
            }}
            margin={"0"}
            is_loading={mutation.isLoading}
          />
        </FlexRowRight>
      </SectionArea>
      <HorizontalSpacer spacepixels={20} />
      <SectionArea>
        <H4>Locks</H4>
        <HorizontalSpacer spacepixels={20} />
        <UnlockModuleForm team_id={team_id} module={module} />
      </SectionArea>
      <HorizontalSpacer spacepixels={20} />
      <SectionArea>
        <H4>Delete Module</H4>
        <HorizontalSpacer spacepixels={20} />
        <DeleteModuleForm team_id={team_id} module={module} />
      </SectionArea>
    </ModuleSettingsContainer>
  );
};

export default ModuleSettings;
