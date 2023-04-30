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
import { ExpandedRunContainer, RunSectionCard } from "./styles";
import { Module } from "shared/api/generated/data-contracts";
import Logs from "components/logs";
import Status from "components/status";
import GithubRef from "components/githubref";
import { StatusText } from "components/status/styles";
import ConfigurationSection from "./components/ConfigurationSection";
import ExpandedRunOverview from "../expandedrunoverview";

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

  const logsQuery = useQuery({
    queryKey: ["module_run_logs", team_id, module_id, module_run_id],
    queryFn: async () => {
      const res = await api.getModuleRunLogs(team_id, module_id, module_run_id);
      return res;
    },
    retry: false,
  });

  const renderLogsSection = () => {
    if (logsQuery.isLoading) {
      return (
        <Placeholder>
          <Spinner />
        </Placeholder>
      );
    }

    if (logsQuery.isError) {
      return (
        <Placeholder>
          <SmallSpan>Could not load logs: an error occurred.</SmallSpan>
        </Placeholder>
      );
    }

    return <Logs logs={logsQuery.data?.data.logs} />;
  };

  if (moduleRunQuery.isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  return (
    <ExpandedRunContainer>
      <HorizontalSpacer spacepixels={24} />
      <BackText text="All Runs" back={back} />
      <HorizontalSpacer spacepixels={24} />
      <ExpandedRunOverview
        team_id={team_id}
        module={module}
        module_run={moduleRunQuery?.data?.data}
      />
      <HorizontalSpacer spacepixels={24} />
      <ConfigurationSection
        team_id={team_id}
        module_id={module.id}
        module_run={moduleRunQuery.data.data}
      />
      <HorizontalSpacer spacepixels={24} />
      <RunSectionCard>
        <H4>Logs</H4>
        {renderLogsSection()}
      </RunSectionCard>
    </ExpandedRunContainer>
  );
};

export default ExpandedRun;
