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
} from "@hatchet-dev/hatchet-components";
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
import ExpandedResultOverview from "../expandedresultoverview";
import { ExpandedResultContainer, ResultSectionCard } from "./styles";

export type Props = {
  team_id: string;
  module_monitor: ModuleMonitor;
  module_monitor_result: ModuleMonitorResult;
  back: () => void;
};

const ExpandedResult: React.FC<Props> = ({
  team_id,
  module_monitor,
  module_monitor_result,
  back,
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

  const getTriggerDescription = () => {
    switch (module_monitor.kind) {
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
    if (moduleQuery.isLoading) {
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
    <ExpandedResultContainer>
      <HorizontalSpacer spacepixels={24} />
      <BackText text="All Results" back={back} />
      <HorizontalSpacer spacepixels={24} />
      <ExpandedResultOverview
        team_id={team_id}
        module_monitor_result={module_monitor_result}
      />
      <ConfigurationSection
        team_id={team_id}
        module_id={module_monitor_result.module_id}
        module_run={moduleRunQuery.data.data}
      />
    </ExpandedResultContainer>
  );
};

export default ExpandedResult;
