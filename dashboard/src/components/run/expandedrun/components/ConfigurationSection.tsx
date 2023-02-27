import {
  H4,
  HorizontalSpacer,
  SmallSpan,
  FlexColScroll,
  Placeholder,
  Spinner,
} from "@hatchet-dev/hatchet-components";
import { useQuery } from "@tanstack/react-query";
import CodeBlock from "components/codeblock";
import EnvVars, { getInternalEnvVars, newEnvVarAtom } from "components/envvars";
import React, { useMemo } from "react";
import api from "shared/api";
import { Module, ModuleRun } from "shared/api/generated/data-contracts";
import GithubRef from "components/githubref";
import { useAtom } from "jotai";
import { RunSectionCard } from "./RunSectionCard";

type Props = {
  team_id: string;
  module_id: string;
  module_run: ModuleRun;
};

const ConfigurationSection: React.FC<Props> = ({
  team_id,
  module_id,
  module_run,
}) => {
  const envVarAtom = useMemo(() => {
    return newEnvVarAtom([]);
  }, []);

  const [envVars, setEnvVars] = useAtom(envVarAtom);

  const envVarsQuery = useQuery({
    queryKey: [
      "module_run_env_vars",
      team_id,
      module_id,
      module_run?.config?.env_var_version_id,
    ],
    queryFn: async () => {
      const res = await api.getModuleEnvVars(
        team_id,
        module_id,
        module_run?.config?.env_var_version_id
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
    },
    retry: false,
    enabled: !!module_run?.config?.env_var_version_id,
  });

  const valuesQuery = useQuery({
    queryKey: [
      "module_run_values",
      team_id,
      module_id,
      module_run?.config?.values_version_id,
    ],
    queryFn: async () => {
      const res = await api.getModuleValues(
        team_id,
        module_id,
        module_run?.config?.values_version_id
      );

      return res;
    },
    retry: false,
    enabled: !!module_run?.config?.values_version_id,
  });

  const getGithubFileRefLink = () => {
    const gh = valuesQuery.data.data?.github;
    const sha = module_run?.config.github_commit_sha;

    return `https://github.com/${gh.github_repo_owner}/${
      gh.github_repo_name
    }/blob/${sha}${gh.path.replace(/^(\.)/, "")}`;
  };

  const renderValuesSection = () => {
    if (valuesQuery.isLoading) {
      return (
        <Placeholder>
          <Spinner />
        </Placeholder>
      );
    }

    if (valuesQuery.data.data.github) {
      return (
        <GithubRef
          text={valuesQuery.data.data.github.path}
          link={getGithubFileRefLink()}
        />
      );
    }

    return (
      <FlexColScroll height="200px" width="100%">
        <CodeBlock
          value={JSON.stringify(valuesQuery?.data?.data?.raw_values)}
          height="200px"
          readOnly={true}
        />
      </FlexColScroll>
    );
  };

  const renderEnvVarsSection = () => {
    if (envVarsQuery.isLoading) {
      return (
        <Placeholder>
          <Spinner />
        </Placeholder>
      );
    }

    return <EnvVars envVarAtom={envVarAtom} read_only={true} />;
  };

  return (
    <RunSectionCard>
      <H4>Configuration</H4>
      <HorizontalSpacer spacepixels={12} />
      <FlexColScroll>
        <SmallSpan>Values:</SmallSpan>
        <HorizontalSpacer spacepixels={12} />
        {renderValuesSection()}
        <HorizontalSpacer spacepixels={12} />
        <SmallSpan>Environment variables:</SmallSpan>
        <HorizontalSpacer spacepixels={12} />
        {renderEnvVarsSection()}
      </FlexColScroll>
    </RunSectionCard>
  );
};

export default ConfigurationSection;
