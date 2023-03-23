import {
  FlexRowRight,
  H1,
  HorizontalSpacer,
  P,
  Filter,
  Paginator,
  Table,
  StandardButton,
  Placeholder,
  Spinner,
} from "@hatchet-dev/hatchet-components";
import { useQuery } from "@tanstack/react-query";
import { useAtom } from "jotai";
import React from "react";
import { useHistory } from "react-router-dom";
import api from "shared/api";
import { currTeamAtom } from "shared/atoms/atoms";
import usePagination from "shared/hooks/usepagination";
import { capitalize, relativeDate } from "shared/utils";

const MonitoringView: React.FunctionComponent = () => {
  let history = useHistory();
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);

  const {
    currentPage,
    maxPage,
    cursor_forward,
    cursor_backward,
    set_data,
  } = usePagination();

  const listMonitorsQuery = useQuery({
    queryKey: ["monitors", currTeam?.id, currentPage],
    queryFn: async () => {
      const res = await api.listMonitors(currTeam?.id, {
        page: currentPage,
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
      Header: "Name",
      accessor: "name",
      width: 100,
    },
    {
      Header: "Description",
      accessor: "description",
      width: 300,
    },
    {
      Header: "Kind",
      accessor: "kind",
      width: 100,
      Cell: ({ row }: any) => {
        switch (row.original.kind) {
          case "plan":
            return <div>Plan</div>;
          case "state":
            return <div>State</div>;
          case "before_plan":
            return <div>Before Plan</div>;
          case "after_plan":
            return <div>After Plan</div>;
          case "before_apply":
            return <div>Before Apply</div>;
          case "after_apply":
            return <div>After Apply</div>;
          case "before_destroy":
            return <div>Before Destroy</div>;
          case "after_destroy":
            return <div>After Destroy</div>;
        }
      },
    },
    {
      Header: "Created",
      accessor: "created_at",
      width: 100,
    },
    {
      Header: "Triggers/Cron",
      accessor: "updated_at",
      width: 100,
      Cell: ({ row }: any) => {
        if (row.original.cron_schedule) {
          return <div>{row.original.cron_schedule}</div>;
        }

        switch (row.original.kind) {
          case "before_plan":
            return <div>Runs before plan</div>;
          case "after_plan":
            return <div>Runs after plan</div>;
          case "before_apply":
            return <div>Runs before apply</div>;
          case "after_apply":
            return <div>Runs after apply</div>;
          case "before_destroy":
            return <div>Runs before destroy</div>;
          case "after_destroy":
            return <div>Runs after destroy</div>;
        }
      },
    },
  ];

  const handleResourceClick = (row: any) => {
    history.push(`/teams/${currTeam.id}/monitors/${row.original.id}`);
  };

  const handleCreateMonitorClick = () => {
    history.push(`/teams/${currTeam.id}/monitors/create/step_1`);
  };

  const tableData =
    listMonitorsQuery.data?.data?.rows?.map((monitor) => {
      return {
        id: monitor.id,
        name: capitalize(monitor.name),
        description: monitor.description,
        kind: monitor.kind,
        updated_at: relativeDate(monitor.updated_at),
        created_at: relativeDate(monitor.created_at),
        cron_schedule: monitor.cron_schedule,
      };
    }) || [];

  const renderMonitors = () => {
    if (listMonitorsQuery.isLoading) {
      return (
        <Placeholder>
          <Spinner />
        </Placeholder>
      );
    }

    return (
      <Table
        rowHeight={"3.5em"}
        columns={columns}
        data={tableData}
        onRowClick={handleResourceClick}
      />
    );
  };

  return (
    <>
      <H1>Monitors</H1>
      <HorizontalSpacer spacepixels={12} />
      <P>
        This dashboard displays all monitors that are set up for your modules.
      </P>
      <FlexRowRight>
        <StandardButton
          label="Create monitor"
          material_icon="add"
          on_click={handleCreateMonitorClick}
          margin="0"
        />
      </FlexRowRight>
      <HorizontalSpacer spacepixels={30} />
      {renderMonitors()}
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

export default MonitoringView;
