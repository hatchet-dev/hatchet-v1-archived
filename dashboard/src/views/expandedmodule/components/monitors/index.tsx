import StandardButton from "components/buttons";
import { FlexRowRight, H4, HorizontalSpacer } from "components/globals";
import SectionArea from "components/sectionarea";
import Table from "components/table";
import React from "react";

const ExpandedModuleMonitors: React.FC = () => {
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

  const handleResourceClick = (row: any) => {};

  return (
    <>
      <HorizontalSpacer spacepixels={12} />
      <SectionArea>
        <H4>Monitors</H4>
        <FlexRowRight>
          <StandardButton label="Add monitor" material_icon="add" />
        </FlexRowRight>
        <HorizontalSpacer spacepixels={20} />
        <Table
          rowHeight={"3.5em"}
          columns={columns}
          data={data}
          onRowClick={handleResourceClick}
        />
      </SectionArea>
    </>
  );
};

export default ExpandedModuleMonitors;
