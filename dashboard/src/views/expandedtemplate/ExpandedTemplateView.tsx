import {
  FlexRowRight,
  H1,
  HorizontalSpacer,
  P,
  Span,
  Breadcrumbs,
  Paginator,
  Table,
  TabList,
} from "hatchet-components";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import DetailedVersionList from "./components/detailedversionlist";

const TabOptions = ["Modules", "Versions", "Settings"];

const ExpandedTemplateView: React.FunctionComponent = () => {
  const [selectedTab, setSelectedTab] = useState(TabOptions[0]);

  let history = useHistory();

  const columns = [
    {
      Header: "Repo Name",
      accessor: "name",
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
      Header: "Version",
      accessor: "version",
    },
  ];

  const data = [
    {
      id: "1111",
      name: "my-company/team-1",
      last_deployed: "10 days ago",
      source: "Github",
      path: "./infra/gke",
      version: "v0.25.3",
    },
    {
      id: "2222",
      name: "my-company/team-2",
      last_deployed: "15 days ago",
      source: "Github",
      path: "./infra/gke",
      version: "v0.25.3",
    },
    {
      id: "3333",
      name: "my-company/team-3",
      last_deployed: "23 days ago",
      source: "Github",
      path: "./infra/gke",
      version: "v0.25.3",
    },
  ];

  const handleResourceClick = (row: any) => {
    history.push(`/modules/${row.original.id}`);
  };

  const renderTabContents = () => {
    switch (selectedTab) {
      case "Modules":
        return (
          <>
            <HorizontalSpacer spacepixels={16} />
            <Table
              columns={columns}
              data={data}
              onRowClick={handleResourceClick}
            />
            {/* <FlexRowRight>
              <Paginator />
            </FlexRowRight> */}
          </>
        );
      case "Versions":
        return (
          <DetailedVersionList
            versions={[
              {
                version: "v0.2.0",
                link: "https://github.com/hatchet-dev/hatchet/releases/v0.35.0",
              },
              {
                version: "v0.1.0",
                link: "https://github.com/hatchet-dev/hatchet/releases/v0.35.0",
              },
            ]}
          />
        );
      default:
        return <Span>Settings</Span>;
    }
  };

  return (
    <>
      <Breadcrumbs
        breadcrumbs={[
          {
            label: "Templates",
            link: "/templates",
          },
          {
            label: "Google Kubernetes Engine (GKE)",
            link: "",
          },
        ]}
      />
      <H1>Google Kubernetes Engine (GKE)</H1>
      <HorizontalSpacer spacepixels={20} />
      <P>
        This page contains information about which resources are using the
        Google Kubernetes Engine (GKE) template.
      </P>
      <HorizontalSpacer spacepixels={20} />
      <TabList tabs={TabOptions} selectTab={setSelectedTab} />
      {renderTabContents()}
    </>
  );
};

export default ExpandedTemplateView;
