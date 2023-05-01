import {
  H4,
  HorizontalSpacer,
  BackText,
  FlexRow,
  SmallSpan,
  MaterialIcon,
  Placeholder,
  Spinner,
  FlexCol,
  FlexRowRight,
  FlexRowLeft,
} from "hatchet-components";
import { useQuery } from "@tanstack/react-query";
import React from "react";
import api from "shared/api";
import { parseTerraformPlanSummary, relativeDate } from "shared/utils";
import { Module, ModuleRun } from "shared/api/generated/data-contracts";
import Logs from "components/logs";
import Status from "components/status";
import GithubRef from "components/githubref";
import { StatusText } from "components/status/styles";
import { RunSectionCard } from "../expandedrun/styles";
import { TriggerPRContainer } from "./styles";
import CodeBlock from "components/codeblock";

type Props = {
  team_id: string;
  module: Module;
  module_run: ModuleRun;
};

const ExpandedRunOverview: React.FC<Props> = ({
  team_id,
  module,
  module_run,
}) => {
  const module_id = module.id;

  const status = module_run?.status;
  const kind = module_run?.kind;
  const triggerKind = module_run?.config?.trigger_kind;

  const planSummaryEnabled =
    (kind == "plan" && status == "completed") ||
    (kind == "apply" && triggerKind == "github");

  const planSummaryQuery = useQuery({
    queryKey: ["module_run_plan_summary", team_id, module_id, module_run.id],
    queryFn: async () => {
      const res = await api.getModuleRunPlanSummary(
        team_id,
        module_id,
        module_run.id
      );
      return res;
    },
    retry: false,
    enabled: planSummaryEnabled,
  });

  const selectPR = () => {
    const pr = module_run?.github_pull_request;
    window.open(
      `https://github.com/${pr.github_repository_owner}/${pr.github_repository_name}/pull/${pr.github_pull_request_number}`
    );
  };

  const getPRCommitLink = () => {
    const pr = module_run?.github_pull_request;
    const sha = module_run?.config.git_commit_sha;

    return `https://github.com/${pr.github_repository_owner}/${pr.github_repository_name}/pull/${pr.github_pull_request_number}/commits/${sha}`;
  };

  const getCommitLink = () => {
    const gh = module.deployment;
    const sha = module_run?.config.git_commit_sha;

    return `https://github.com/${gh.git_repo_owner}/${gh.git_repo_name}/commit/${sha}`;
  };

  if (planSummaryEnabled && planSummaryQuery.isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  const renderStatusContainer = () => {
    let materialIcon = "";
    let text = "";

    switch (module_run?.status) {
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

    return <Status material_icon={materialIcon} status_text={text} />;
  };

  const renderTimeContainer = () => {
    return (
      <Status
        status_text={relativeDate(module_run?.updated_at)}
        material_icon="schedule"
      />
    );
  };

  const renderPlannedChanges = () => {
    if (!planSummaryQuery.data?.data) {
      return null;
    }

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

  const renderGithubTrigger = () => {
    return (
      <FlexRowLeft gap="4px">
        <SmallSpan>Triggered by </SmallSpan>
        <TriggerPRContainer onClick={selectPR}>
          <MaterialIcon className="fa-solid fa-code-pull-request" />
          <StatusText>
            Pull Request #
            {module_run?.github_pull_request.github_pull_request_number}
          </StatusText>
        </TriggerPRContainer>
        <SmallSpan>into</SmallSpan>
        <CodeBlock
          value={
            module_run?.github_pull_request.github_pull_request_base_branch
          }
        />
        <SmallSpan>from</SmallSpan>
        <CodeBlock
          value={
            module_run?.github_pull_request.github_pull_request_head_branch
          }
        />
      </FlexRowLeft>
    );
  };

  const renderManualTrigger = () => {
    return (
      <FlexRowLeft gap="4px">
        <SmallSpan>Triggered by a manual run.</SmallSpan>
      </FlexRowLeft>
    );
  };

  const renderPlanOverview = () => {
    return (
      <FlexCol>
        <SmallSpan>{module_run?.status_description} </SmallSpan>
        <HorizontalSpacer spacepixels={8} />
        {triggerKind == "github"
          ? renderGithubTrigger()
          : renderManualTrigger()}
        <HorizontalSpacer spacepixels={8} />
        {status == "completed" && renderPlannedChanges()}
      </FlexCol>
    );
  };

  const renderApplyOverview = () => {
    return (
      <FlexCol>
        <SmallSpan>{module_run?.status_description} </SmallSpan>
        <HorizontalSpacer spacepixels={8} />
        {renderPlannedChanges()}
      </FlexCol>
    );
  };

  const renderInitOverview = () => {
    return (
      <FlexCol>
        <SmallSpan>{module_run?.status_description} </SmallSpan>
        <HorizontalSpacer spacepixels={8} />
      </FlexCol>
    );
  };

  const renderOverview = () => {
    switch (module_run?.kind) {
      case "plan":
        return renderPlanOverview();
      case "apply":
        return renderApplyOverview();
      case "init":
        return renderInitOverview();
    }
  };

  return (
    <RunSectionCard>
      <FlexRow>
        <H4>Overview</H4>
        <FlexRowRight gap="8px">
          {triggerKind == "github" && (
            <GithubRef
              text={module_run?.config.git_commit_sha.substr(0, 7)}
              link={kind == "plan" ? getPRCommitLink() : getCommitLink()}
            />
          )}
          {renderStatusContainer()}
          {renderTimeContainer()}
        </FlexRowRight>
      </FlexRow>
      <HorizontalSpacer spacepixels={20} />
      {renderOverview()}
    </RunSectionCard>
  );
};

export default ExpandedRunOverview;
