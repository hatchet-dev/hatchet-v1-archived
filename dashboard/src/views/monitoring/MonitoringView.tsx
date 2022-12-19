import Filter from "components/filter";
import { FlexRowRight, H1, HorizontalSpacer, P } from "components/globals";
import Paginator from "components/paginator";
import Table from "components/table";
import React from "react";
import { useHistory } from "react-router-dom";

const MonitoringView: React.FunctionComponent = () => {
  let history = useHistory();

  const columns = [
    {
      Header: "Alert",
      accessor: "alert",
      width: 450,
    },
    {
      Header: "Last seen",
      accessor: "last_seen",
      width: 100,
    },
    {
      Header: "Environment",
      accessor: "environment",
      width: 100,
    },
    {
      Header: "Repo name",
      accessor: "repo_name",
      width: 100,
    },
  ];

  const data = [
    {
      id: "1111",
      alert: "A drift was detected in 3 resources",
      repo_name: "my-company/team-1",
      environment: "Staging",
      last_seen: "2 days ago",
      is_resolved: false,
    },
    {
      id: "2222",
      alert:
        "The custom policy 'EKS version must be at least 1.21' was violated",
      repo_name: "my-company/team-1",
      environment: "Production",
      last_seen: "4 days ago",
      is_resolved: true,
    },
    {
      id: "3333",
      alert: "A drift was detected in 1 resource",
      repo_name: "my-company/team-1",
      environment: "Staging",
      last_seen: "5 days ago",
      is_resolved: true,
    },
  ];

  const handleResourceClick = (row: any) => {
    history.push(`/monitoring/${row.original.id}`);
  };

  return (
    <>
      <H1>Monitoring</H1>
      <HorizontalSpacer spacepixels={12} />
      <P>
        This dashboard displays alerts across all infrastructure deployments.
      </P>
      <FlexRowRight>
        <Filter />
      </FlexRowRight>
      <Table
        rowHeight={"3.5em"}
        columns={columns}
        data={data}
        onRowClick={handleResourceClick}
      />
      <FlexRowRight>
        <Paginator />
      </FlexRowRight>
    </>
  );
};

export default MonitoringView;
