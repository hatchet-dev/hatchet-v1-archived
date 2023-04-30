import {
  FlexColScroll,
  FlexRowRight,
  MaterialIcon,
  Paginator,
  Placeholder,
  Spinner,
  Table,
} from "hatchet-components";
import { useQuery } from "@tanstack/react-query";
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
import ExpandedResult from "../expandedresult";

export type Props = {
  team_id: string;
  module_id?: string;
  module_monitor?: ModuleMonitor;
};

const ResultsList: React.FC<Props> = ({
  team_id,
  module_id,
  module_monitor,
}) => {
  const [selectedResult, setSelectedResult] = useState<ModuleMonitorResult>();

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
      module_monitor?.id,
    ],
    queryFn: async () => {
      const res = await api.listMonitorResults(team_id, {
        page: currentPage,
        module_id: module_id,
        module_monitor_id: module_monitor?.id,
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
      width: 60,
      Cell: ({ row }: any) => {
        return (
          <Status
            kind="color"
            color={
              row.original.status == "succeeded"
                ? theme.text.default
                : "#ff385d"
            }
            status_text={capitalize(row.original.status)}
            material_icon="check"
          />
        );
      },
    },
    {
      Header: "Module",
      accessor: "module_name",
      width: 200,
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
      id: result.id,
      status: result.status,
      title: result.title,
      message: result.message,
      created_at: relativeDate(result.created_at),
      module_id: result.module_id,
      module_name: result.module_name,
    };
  });

  if (selectedResult) {
    return (
      <ExpandedResult
        team_id={team_id}
        module_monitor={module_monitor}
        module_monitor_result={selectedResult}
        back={() => setSelectedResult(null)}
      />
    );
  }

  return (
    <>
      <FlexColScroll height="calc(100% - 250px)">
        <Table
          rowHeight={"3.5em"}
          columns={columns}
          data={tableData}
          onRowClick={(row: any) => {
            const matchedResult = listResultsQuery?.data?.data?.rows?.filter(
              (result) => {
                return result.id == row.original.id;
              }
            );

            if (matchedResult.length == 1) {
              setSelectedResult(matchedResult[0]);
            }
          }}
        />
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
