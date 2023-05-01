import {
  BackText,
  FlexColScroll,
  FlexRowRight,
  HorizontalSpacer,
  MaterialIcon,
  Paginator,
  Placeholder,
  Spinner,
  Table,
  FlexRow,
  H4,
  P,
} from "hatchet-components";
import { useQuery } from "@tanstack/react-query";
import ConfigurationSection from "components/run/expandedrun/components/ConfigurationSection";
import Status from "components/status";
import React, { useState } from "react";
import api from "shared/api";
import {
  ModuleRun,
  ModuleMonitor,
  ModuleMonitorResult,
} from "shared/api/generated/data-contracts";
import usePagination from "shared/hooks/usepagination";
import theme from "shared/theme";
import { capitalize, relativeDate } from "shared/utils";
import { ResultSectionCard } from "../expandedresult/styles";

export type Props = {
  team_id: string;
  module_monitor_result: ModuleMonitorResult;
};

const ExpandedResultOverview: React.FC<Props> = ({
  team_id,
  module_monitor_result,
}) => {
  const moduleRunQuery = useQuery({
    queryKey: [
      "module_run",
      team_id,
      module_monitor_result.module_id,
      module_monitor_result.module_run_id,
    ],
    queryFn: async () => {
      const res = await api.getModuleRun(
        team_id,
        module_monitor_result.module_id,
        module_monitor_result.module_run_id
      );
      return res;
    },
    retry: false,
  });

  const moduleQuery = useQuery({
    queryKey: ["module", team_id, module_monitor_result.module_id],
    queryFn: async () => {
      const res = await api.getModule(team_id, module_monitor_result.module_id);
      return res;
    },
    retry: false,
  });

  const moduleMonitorQuery = useQuery({
    queryKey: ["monitor", team_id, module_monitor_result.module_monitor_id],
    queryFn: async () => {
      const res = await api.getMonitor(
        team_id,
        module_monitor_result.module_monitor_id
      );
      return res;
    },
    retry: false,
  });

  const getTriggerDescription = () => {
    switch (moduleMonitorQuery.data.data.kind) {
      case "plan":
        return "scheduled plan check";
      case "state":
        return "scheduled state check";
      case "before_plan":
        return "before plan";
      case "after_plan":
        return "after plan";
      case "before_apply":
        return "before apply";
      case "after_apply":
        return "after apply";
    }
  };

  const renderTrigger = () => {
    if (moduleQuery.isLoading || moduleMonitorQuery.isLoading) {
      return (
        <Placeholder>
          <Spinner />
        </Placeholder>
      );
    }

    return (
      <P>
        Triggered by {getTriggerDescription()} for module{" "}
        {moduleQuery?.data?.data?.name}
      </P>
    );
  };

  if (moduleRunQuery.isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  return (
    <ResultSectionCard>
      <FlexRow>
        <H4>Overview</H4>
        <FlexRowRight gap="8px">
          <Status
            status_text={relativeDate(module_monitor_result.created_at)}
            material_icon="schedule"
          />
          <Status
            kind="color"
            color={
              module_monitor_result.status == "succeeded"
                ? theme.text.default
                : "#ff385d"
            }
            status_text={capitalize(module_monitor_result.status)}
          />
        </FlexRowRight>
      </FlexRow>
      <HorizontalSpacer spacepixels={12} />
      <P>{module_monitor_result.message}</P>
      <HorizontalSpacer spacepixels={20} />
      {renderTrigger()}
      <HorizontalSpacer spacepixels={20} />
    </ResultSectionCard>
  );
};

export default ExpandedResultOverview;
