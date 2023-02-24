import {
  FlexColScroll,
  FlexRowRight,
  MaterialIcon,
  Paginator,
  Placeholder,
  Spinner,
  Table,
} from "@hatchet-dev/hatchet-components";
import { useQuery } from "@tanstack/react-query";
import React, { useState } from "react";
import api from "shared/api";
import { ModuleRun } from "shared/api/generated/data-contracts";
import usePagination from "shared/hooks/usepagination";
import { relativeDate } from "shared/utils";

export type Props = {
  team_id: string;
  module_id?: string;
  module_monitor_id?: string;
};

const ResultsList: React.FC<Props> = ({
  team_id,
  module_id,
  module_monitor_id,
}) => {
  const {
    currentPage,
    maxPage,
    cursor_forward,
    cursor_backward,
    set_data,
  } = usePagination();

  const listResultsQuery = useQuery({
    queryKey: [
      "monitor_results",
      team_id,
      currentPage,
      module_id,
      module_monitor_id,
    ],
    queryFn: async () => {
      const res = await api.listMonitorResults(team_id, {
        page: currentPage,
        module_id: module_id,
        module_monitor_id: module_monitor_id,
      });

      return res;
    },
    retry: false,
    onSuccess: (data) => {
      set_data(data?.data?.pagination);
    },
  });

  const columns = [
    {
      Header: "Status",
      accessor: "status",
      width: 100,
    },
    {
      Header: "Message",
      accessor: "message",
      width: 300,
    },
    {
      Header: "Run at",
      accessor: "created_at",
      width: 100,
    },
  ];

  if (listResultsQuery.isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  const tableData = listResultsQuery?.data?.data?.rows?.map((result) => {
    return {
      status: result.status,
      title: result.title,
      message: result.message,
      created_at: relativeDate(result.created_at),
    };
  });

  return (
    <>
      <FlexColScroll height="calc(100% - 250px)">
        <Table rowHeight={"3.5em"} columns={columns} data={tableData} />
      </FlexColScroll>
      <FlexRowRight>
        <Paginator
          curr_page={currentPage}
          last_page={maxPage}
          cursor_backward={cursor_backward}
          cursor_forward={cursor_forward}
        />
      </FlexRowRight>
    </>
  );
};

export default ResultsList;
