import {
  StandardButton,
  FlexRowRight,
  HorizontalSpacer,
  H4,
  Placeholder,
  Spinner,
} from "@hatchet-dev/hatchet-components";
import { useMutation, useQuery } from "@tanstack/react-query";
import EnvVars from "components/envvars";
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
import ModuleSettingsCard from "./components/ModuleSettingsCard";
import ModuleSettingsContainer from "./components/ModuleSettingsContainer";
import UpdateModuleName from "./components/UpdateModuleName";

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
  const [envVars, setEnvVars] = useState<string[]>([]);

  const request = useMemo<UpdateModuleRequest>(() => {
    return {
      name: name,
      github: githubParams,
    };
  }, [githubParams, name]);

  const { refetch } = useQuery({
    queryKey: ["module", team_id, module.id],
    queryFn: async () => {
      const res = await api.getModule(team_id, module.id);
      return res;
    },
    retry: false,
  });

  const envVarsQuery = useQuery({
    queryKey: ["module_env_vars", team_id, module_id],
    queryFn: async () => {
      const res = await api.getModuleEnvVars(
        team_id,
        module_id,
        module.current_env_vars_version_id
      );

      setEnvVars(
        res.data.env_vars.map((envVar) => {
          return `${envVar.key}~~=~~${envVar.val}`;
        })
      );

      return res;
    },
    retry: false,
  });

  const valuesQuery = useQuery({
    queryKey: ["module_values", team_id, module_id],
    queryFn: async () => {
      const res = await api.getModuleValues(
        team_id,
        module_id,
        module.current_values_version_id
      );

      return res;
    },
    retry: false,
  });

  const mutation = useMutation({
    mutationKey: ["update_module", team_id, module_id],
    mutationFn: (request: UpdateModuleRequest) => {
      return api.updateModule(team_id, module_id, request);
    },
    onSuccess: (data) => {
      setErr("");
      refetch();
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
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

  const gh = valuesQuery.data?.data?.github;

  return (
    <ModuleSettingsContainer>
      <ModuleSettingsCard>
        <H4>Configuration</H4>
        <HorizontalSpacer spacepixels={20} />
        <UpdateModuleName module={module} setModuleName={setName} />
        <HorizontalSpacer spacepixels={60} />
        <SelectGitSource
          set_request={setGithubParams}
          current_params={{
            github_app_installation_id:
              module.deployment.github_app_installation_id,
            github_repository_owner: module.deployment.github_repo_owner,
            github_repository_name: module.deployment.github_repo_name,
            github_repository_branch: module.deployment.github_repo_branch,
            path: module.deployment.path,
          }}
        />
        <HorizontalSpacer spacepixels={60} />
        <SetModuleValues
          set_github_values={setGithubValueParams}
          current_github_params={{
            github_app_installation_id: gh?.github_app_installation_id,
            github_repository_branch: gh?.github_repo_branch,
            github_repository_owner: gh?.github_repo_owner,
            github_repository_name: gh?.github_repo_name,
            path: gh?.path,
          }}
          set_raw_values={setRawValues}
          current_raw_values={valuesQuery.data?.data.raw_values}
          set_values_source={setValuesSource}
          current_values_source={gh ? "github" : "raw"}
        />
        <HorizontalSpacer spacepixels={24} />
        <EnvVars envVars={envVars} setEnvVars={setEnvVars} />
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
      </ModuleSettingsCard>
      <HorizontalSpacer spacepixels={20} />
      <ModuleSettingsCard>
        <H4>Locks</H4>
        <HorizontalSpacer spacepixels={20} />
      </ModuleSettingsCard>
      <HorizontalSpacer spacepixels={20} />
      <ModuleSettingsCard>
        <H4>Delete Module</H4>
        <HorizontalSpacer spacepixels={20} />
      </ModuleSettingsCard>
    </ModuleSettingsContainer>
  );
};

export default ModuleSettings;
