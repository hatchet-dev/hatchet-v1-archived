import {
  H4,
  HorizontalSpacer,
  BackText,
  FlexRow,
  SmallSpan,
  FlexColScroll,
  MaterialIcon,
  Placeholder,
  Spinner,
  CodeLine,
  FlexCol,
} from "@hatchet-dev/hatchet-components";
import { useQuery } from "@tanstack/react-query";
import CodeBlock from "components/codeblock";
import EnvVars from "components/envvars";
import React from "react";
import api from "shared/api";
import { parseTerraformPlanSummary, relativeDate } from "shared/utils";
import {
  ExpandedRunContainer,
  GithubRefContainer,
  GithubImg,
  RunSectionCard,
  StatusContainer,
  StatusText,
  TriggerContainer,
  TriggerPRContainer,
  StatusAndCommitContainer,
} from "./styles";
import github from "assets/github.png";
import { Module } from "shared/api/generated/data-contracts";

type Props = {
  back: () => void;
  team_id: string;
  module: Module;
  module_run_id: string;
};

const ExpandedRun: React.FC<Props> = ({
  back,
  team_id,
  module,
  module_run_id,
}) => {
  const module_id = module.id;

  const moduleRunQuery = useQuery({
    queryKey: ["module_run", team_id, module_id, module_run_id],
    queryFn: async () => {
      const res = await api.getModuleRun(team_id, module_id, module_run_id);
      return res;
    },
    retry: false,
  });

  const status = moduleRunQuery?.data?.data?.status;
  const kind = moduleRunQuery?.data?.data?.kind;
  const triggerKind = moduleRunQuery?.data?.data?.config?.trigger_kind;

  const planSummaryEnabled =
    (kind == "plan" && status == "completed") ||
    (kind == "apply" && triggerKind == "github");

  const planSummaryQuery = useQuery({
    queryKey: ["module_run_plan_summary", team_id, module_id, module_run_id],
    queryFn: async () => {
      const res = await api.getModuleRunPlanSummary(
        team_id,
        module_id,
        module_run_id
      );
      return res;
    },
    retry: false,
    enabled: planSummaryEnabled,
  });

  const envVarsQuery = useQuery({
    queryKey: [
      "module_run_env_vars",
      team_id,
      module_id,
      moduleRunQuery?.data?.data?.config?.env_var_version_id,
    ],
    queryFn: async () => {
      const res = await api.getModuleEnvVars(
        team_id,
        module_id,
        moduleRunQuery?.data?.data?.config?.env_var_version_id
      );

      return res;
    },
    retry: false,
    enabled: !!moduleRunQuery?.data?.data?.config?.env_var_version_id,
  });

  const valuesQuery = useQuery({
    queryKey: [
      "module_run_values",
      team_id,
      module_id,
      moduleRunQuery?.data?.data?.config?.values_version_id,
    ],
    queryFn: async () => {
      const res = await api.getModuleValues(
        team_id,
        module_id,
        moduleRunQuery?.data?.data?.config?.values_version_id
      );

      return res;
    },
    retry: false,
    enabled: !!moduleRunQuery?.data?.data?.config?.values_version_id,
  });

  const selectPR = () => {
    const pr = moduleRunQuery.data.data?.github_pull_request;
    window.open(
      `https://github.com/${pr.github_repository_owner}/${pr.github_repository_name}/pull/${pr.github_pull_request_number}`
    );
  };

  const selectPRCommit = () => {
    const pr = moduleRunQuery.data.data?.github_pull_request;
    const sha = moduleRunQuery.data.data?.config.github_commit_sha;

    window.open(
      `https://github.com/${pr.github_repository_owner}/${pr.github_repository_name}/pull/${pr.github_pull_request_number}/commits/${sha}`
    );
  };

  const selectCommit = () => {
    const gh = module.deployment;
    const sha = moduleRunQuery.data.data?.config.github_commit_sha;

    window.open(
      `https://github.com/${gh.github_repo_owner}/${gh.github_repo_name}/commit/${sha}`
    );
  };

  const selectGithubFileRef = () => {
    const gh = valuesQuery.data.data?.github;
    const sha = moduleRunQuery.data.data?.config.github_commit_sha;

    window.open(
      `https://github.com/${gh.github_repo_owner}/${
        gh.github_repo_name
      }/blob/${sha}${gh.path.replace(/^(\.)/, "")}`
    );
  };

  if (
    moduleRunQuery.isLoading ||
    (planSummaryEnabled && planSummaryQuery.isLoading)
  ) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

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
        <GithubRefContainer onClick={selectGithubFileRef}>
          <GithubImg src={github} />
          <StatusText>./envs/alexander-test/variables.json</StatusText>
        </GithubRefContainer>
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

    return (
      <EnvVars
        envVars={envVarsQuery.data?.data?.env_vars.map((envVar) => {
          return `${envVar.key}~~=~~${envVar.val}`;
        })}
        read_only={true}
      />
    );
  };

  const renderStatusContainer = () => {
    let materialIcon = "";
    let text = "";

    switch (moduleRunQuery?.data?.data?.status) {
      case "completed":
        materialIcon = "check";
        text = "Completed";
        break;
      case "failed":
        materialIcon = "error";
        text = "Failed";
        break;
      case "in_progress":
        materialIcon = "pending";
        text = "In Progress";
        break;
    }

    return (
      <StatusContainer>
        <MaterialIcon className="material-icons">{materialIcon}</MaterialIcon>
        <StatusText>{text}</StatusText>
      </StatusContainer>
    );
  };

  const renderTimeContainer = () => {
    return (
      <StatusContainer>
        <MaterialIcon className="material-icons">schedule</MaterialIcon>
        <StatusText>
          {relativeDate(moduleRunQuery.data.data?.updated_at)}
        </StatusText>
      </StatusContainer>
    );
  };

  const renderPlannedChanges = () => {
    const [numToCreate, numToUpdate, numToDelete] = parseTerraformPlanSummary(
      planSummaryQuery.data.data
    );

    return (
      <SmallSpan>
        Planned changes: {numToCreate} to create, {numToUpdate} to update,{" "}
        {numToDelete} to delete.
      </SmallSpan>
    );
  };

  const renderPlanOverview = () => {
    return (
      <FlexCol>
        <SmallSpan>{moduleRunQuery?.data?.data?.status_description} </SmallSpan>
        <HorizontalSpacer spacepixels={8} />
        <TriggerContainer>
          <SmallSpan>Triggered by </SmallSpan>
          <TriggerPRContainer onClick={selectPR}>
            <MaterialIcon className="fa-solid fa-code-pull-request" />
            <StatusText>
              Pull Request #
              {
                moduleRunQuery.data.data?.github_pull_request
                  .github_pull_request_number
              }
            </StatusText>
          </TriggerPRContainer>
          <SmallSpan>into</SmallSpan>
          <CodeLine padding="6px">
            {
              moduleRunQuery.data.data?.github_pull_request
                .github_pull_request_base_branch
            }
          </CodeLine>
          <SmallSpan>from</SmallSpan>
          <CodeLine padding="6px">
            {
              moduleRunQuery.data.data?.github_pull_request
                .github_pull_request_head_branch
            }
          </CodeLine>
        </TriggerContainer>
        <HorizontalSpacer spacepixels={8} />
        {status == "completed" && renderPlannedChanges()}
      </FlexCol>
    );
  };

  const renderApplyOverview = () => {
    return (
      <FlexCol>
        <SmallSpan>{moduleRunQuery?.data?.data?.status_description} </SmallSpan>
        <HorizontalSpacer spacepixels={8} />
        {renderPlannedChanges()}
      </FlexCol>
    );
  };

  const renderOverview = () => {
    switch (moduleRunQuery?.data?.data?.kind) {
      case "plan":
        return renderPlanOverview();
      case "apply":
        return renderApplyOverview();
    }
  };

  return (
    <ExpandedRunContainer>
      <HorizontalSpacer spacepixels={24} />
      <BackText text="All Runs" back={back} />
      <HorizontalSpacer spacepixels={24} />
      <RunSectionCard>
        <FlexRow>
          <H4>Overview</H4>
          <StatusAndCommitContainer>
            <GithubRefContainer
              onClick={kind == "plan" ? selectPRCommit : selectCommit}
            >
              <GithubImg src={github} />
              <StatusText>
                {moduleRunQuery.data.data?.config.github_commit_sha.substr(
                  0,
                  7
                )}
              </StatusText>
            </GithubRefContainer>
            {renderStatusContainer()}
            {renderTimeContainer()}
          </StatusAndCommitContainer>
        </FlexRow>
        <HorizontalSpacer spacepixels={20} />
        {renderOverview()}
      </RunSectionCard>
      <HorizontalSpacer spacepixels={24} />
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
      <HorizontalSpacer spacepixels={24} />
      <RunSectionCard>
        <H4>Logs</H4>
      </RunSectionCard>
    </ExpandedRunContainer>
  );
};

export default ExpandedRun;
