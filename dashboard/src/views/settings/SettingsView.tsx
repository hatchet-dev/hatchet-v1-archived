import {
  FlexColCenter,
  H1,
  HorizontalSpacer,
  FlexCol,
  TabList,
} from "@hatchet-dev/hatchet-components";
import React, { useEffect, useState } from "react";
import usePrevious from "shared/hooks/useprevious";
import { useHistory, useLocation } from "react-router-dom";
import OrganizationSettingsView from "views/organizationsettings/OrganizationSettingsView";

const TabOptions = ["Organization Settings", "Team Settings"];

type Props = {
  defaultOption?: string;
};

const SettingsView: React.FunctionComponent<Props> = ({ defaultOption }) => {
  const history = useHistory();
  const location = useLocation();
  const [selectedTab, setSelectedTab] = useState(
    defaultOption || TabOptions[0]
  );
  const prevSelectedTab = usePrevious(selectedTab);

  useEffect(() => {
    if (location.pathname == "/settings") {
      history.push("/settings/organization");
    }
  }, [location.pathname]);

  useEffect(() => {
    if (selectedTab && prevSelectedTab && selectedTab != prevSelectedTab) {
      switch (selectedTab) {
        case "Organization Settings":
          history.push("/settings/organization");
          break;
        case "Team Settings":
          history.push("/settings/team");
          break;
      }
    }
  }, [selectedTab]);

  const renderTabContents = () => {
    switch (selectedTab) {
      case "Organization Settings":
        return <OrganizationSettingsView />;
      case "Team Settings":
        return <div>Not implemented</div>;
    }
  };

  return (
    <FlexColCenter height={"100%"}>
      <FlexCol width="100%" maxWidth="840px" height={"100%"}>
        <H1>Settings</H1>
        <HorizontalSpacer spacepixels={20} />
        <TabList tabs={TabOptions} selectTab={setSelectedTab} />
        <HorizontalSpacer spacepixels={20} />
        {renderTabContents()}
      </FlexCol>
    </FlexColCenter>
  );
};

export default SettingsView;
