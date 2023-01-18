import {
  FlexRowRight,
  H1,
  HorizontalSpacer,
  P,
  StandardButton,
  Filter,
  Paginator,
  Table,
} from "@hatchet-dev/hatchet-components";
import React from "react";
import { useHistory } from "react-router-dom";

const ModulesView: React.FunctionComponent = () => {
  let history = useHistory();

  const columns = [
    {
      Header: "Repo Name",
      accessor: "name",
    },
    {
      Header: "Environment",
      accessor: "environment",
    },
    {
      Header: "Last Deployed",
      accessor: "last_deployed",
    },
    {
      Header: "Source",
      accessor: "source",
    },
    {
      Header: "Module Path",
      accessor: "path",
    },
    {
      Header: "Latest Commit",
      accessor: "latest_commit",
    },
  ];

  const data = [
    {
      id: "1111",
      environment: "Staging",
      name: "my-company/team-1",
      last_deployed: "10 days ago",
      source: "Github",
      path: "./infra/gke",
      latest_commit: "7e6f221",
    },
    {
      id: "2222",
      environment: "Production",
      name: "my-company/team-2",
      last_deployed: "15 days ago",
      source: "Github",
      path: "./infra/gke",
      latest_commit: "7e6f221",
    },
    {
      id: "3333",
      environment: "Production",
      name: "my-company/team-3",
      last_deployed: "23 days ago",
      source: "Github",
      path: "./infra/gke",
      latest_commit: "7e6f221",
    },
  ];

  const handleResourceClick = (row: any) => {
    history.push(`/modules/${row.original.id}`);
  };

  const handleLinkModuleClick = () => {
    history.push(`/modules/link/step_1`);
  };

  return (
    <>
      <H1>Modules</H1>
      <HorizontalSpacer spacepixels={12} />
      <P>Modules are all Terraform modules which have a Terraform state.</P>
      <FlexRowRight>
        <Filter />
        <StandardButton
          label="Link module"
          material_icon="add"
          on_click={handleLinkModuleClick}
        />
      </FlexRowRight>
      <HorizontalSpacer spacepixels={12} />
      <Table columns={columns} data={data} onRowClick={handleResourceClick} />
      <FlexRowRight>
        <Paginator />
      </FlexRowRight>
    </>
  );
};

export default ModulesView;
