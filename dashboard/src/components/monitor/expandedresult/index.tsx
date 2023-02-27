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
  ModuleMonitorResult,
} from "shared/api/generated/data-contracts";
import usePagination from "shared/hooks/usepagination";
import theme from "shared/theme";
import { capitalize, relativeDate } from "shared/utils";
import { ExpandedResultContainer, ResultSectionCard } from "./styles";

export type Props = {
  team_id: string;
  module_monitor_result: ModuleMonitorResult;
  back: () => void;
};

const ExpandedResult: React.FC<Props> = ({
  team_id,
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
      <ResultSectionCard>
        <FlexRow>
          <H4>Overview</H4>
          <FlexRowRight gap="8px">
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
      </ResultSectionCard>
      <ConfigurationSection
        team_id={team_id}
        module_id={module_monitor_result.module_id}
        module_run={moduleRunQuery.data.data}
      />
    </ExpandedResultContainer>
  );
};

export default ExpandedResult;
