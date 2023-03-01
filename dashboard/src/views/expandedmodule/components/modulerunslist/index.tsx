import {
  FlexCol,
  FlexColCenter,
  FlexRowRight,
  H1,
  HorizontalSpacer,
  P,
  StyledDeprecatedText,
  Table,
  StandardButton,
  Spinner,
  Placeholder,
  FlexColScroll,
  Paginator,
  SmallSpan,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { relativeDate } from "shared/utils";
import RunsList from "components/run/runslist";
import usePagination from "shared/hooks/usepagination";
import { ModuleRun } from "shared/api/generated/data-contracts";

type Props = {
  team_id: string;
  module_id: string;
  select_run: (run: ModuleRun) => void;
};

const ModuleRunsList: React.FunctionComponent<Props> = ({
  team_id,
  module_id,
  select_run,
}) => {
  const {
    currentPage,
    maxPage,
    cursor_forward,
    cursor_backward,
    set_data,
  } = usePagination();

  const { data, isLoading, refetch, isFetching } = useQuery({
    queryKey: ["module_runs", module_id, currentPage],
    queryFn: async () => {
      const res = await api.listModuleRuns(team_id, module_id, {
        page: currentPage,
        status: "",
      });
      return res;
    },
    retry: false,
    onSuccess: (data) => {
      set_data(data?.data?.pagination);
    },
  });

  if (isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  if (data.data.rows.length == 0) {
    return (
      <FlexCol height="calc(100% - 200px)">
        <HorizontalSpacer spacepixels={20} />
        <Placeholder>
          <SmallSpan>No runs found.</SmallSpan>
        </Placeholder>
      </FlexCol>
    );
  }

  return (
    <FlexCol height="calc(100% - 250px)">
      <FlexColScroll maxHeight="100%">
        <RunsList runs={data.data.rows} select_run={select_run} />
      </FlexColScroll>
      <FlexRowRight>
        <Paginator
          curr_page={currentPage}
          last_page={maxPage}
          cursor_backward={cursor_backward}
          cursor_forward={cursor_forward}
        />
      </FlexRowRight>
    </FlexCol>
  );
};

export default ModuleRunsList;
