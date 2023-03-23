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

const TemplatesView: React.FunctionComponent = () => {
  let history = useHistory();

  const columns = [
    {
      Header: "Name",
      accessor: "name",
    },
    {
      Header: "Last Updated",
      accessor: "last_updated",
    },
    {
      Header: "Source",
      accessor: "source",
    },
    {
      Header: "Usage Count",
      accessor: "usage_count",
    },
    {
      Header: "Latest Version",
      accessor: "latest_version",
    },
  ];

  const data = [
    {
      id: "1111",
      name: "Elastic Kubernetes Service (EKS)",
      last_updated: "10 days ago",
      source: "Github",
      usage_count: 10,
      latest_version: "v0.25.3",
    },
    {
      id: "2222",
      name: "Relational Database Service (RDS)",
      last_updated: "16 days ago",
      source: "Github",
      usage_count: 6,
      latest_version: "v0.14.1",
    },
    {
      id: "3333",
      name: "Google Kubernetes Engine (GKE)",
      last_updated: "22 days ago",
      source: "Github",
      usage_count: 13,
      latest_version: "v0.8.1",
    },
  ];

  const handleModuleClick = (row: any) => {
    history.push(`/templates/${row.original.id}`);
  };

  const handleAddTemplateClick = () => {
    history.push(`/templates/add`);
  };

  return (
    <>
      <H1>Templates</H1>
      <HorizontalSpacer spacepixels={12} />
      <P>
        Templates are reusable pieces of infrastructure, stored in either a
        remote Terraform registry or linked directly from Github.
      </P>
      <FlexRowRight>
        <StandardButton
          label="Add template"
          material_icon="add"
          on_click={handleAddTemplateClick}
        />
      </FlexRowRight>
      <HorizontalSpacer spacepixels={12} />
      <Table columns={columns} data={data} onRowClick={handleModuleClick} />
      <FlexRowRight>
        <Paginator />
      </FlexRowRight>
    </>
  );
};

export default TemplatesView;
