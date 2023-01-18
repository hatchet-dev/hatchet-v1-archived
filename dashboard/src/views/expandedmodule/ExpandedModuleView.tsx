import {
  H1,
  HorizontalSpacer,
  P,
  Span,
  Breadcrumbs,
  HeirarchyGraph,
  TabList,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import RunsList from "../../components/runslist";
import ExpandedModuleMonitors from "./components/monitors";

const TabOptions = [
  "Runs",
  "Resource Explorer",
  "Policies",
  "Monitors",
  "Settings",
];

const ExpandedModuleView: React.FunctionComponent = () => {
  const [selectedTab, setSelectedTab] = useState(TabOptions[0]);

  const renderTabContents = () => {
    switch (selectedTab) {
      case "Runs":
        return (
          <RunsList
            runs={[
              {
                status: "deployed",
                date: "7:09 AM on June 23rd, 2022",
              },
            ]}
          />
        );
      case "Resource Explorer":
        return <HeirarchyGraph width={100} height={100} />;
      case "Monitors":
        return <ExpandedModuleMonitors />;
      default:
        return <Span>Settings</Span>;
    }
  };

  return (
    <>
      <Breadcrumbs
        breadcrumbs={[
          {
            label: "Modules",
            link: "/modules",
          },
          {
            label: "Staging: team-1-gke",
            link: "",
          },
        ]}
      />
      <HorizontalSpacer spacepixels={12} />
      <H1>Staging: team-1-gke</H1>
      <HorizontalSpacer spacepixels={20} />
      <P>
        This page contains information about the team-1-gke workspace in the
        Staging environment.
      </P>
      <HorizontalSpacer spacepixels={20} />
      <TabList tabs={TabOptions} selectTab={setSelectedTab} />
      {renderTabContents()}
    </>
  );
};

export default ExpandedModuleView;
