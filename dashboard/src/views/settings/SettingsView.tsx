import {
  FlexColCenter,
  H1,
  H2,
  HorizontalSpacer,
  P,
  SmallSpan,
  FlexColScroll,
  FlexCol,
} from "components/globals";
import React, { useEffect, useState } from "react";
import { css } from "styled-components";
import theme from "shared/theme";
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import MemberList from "components/organization/memberlist";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import InviteMemberForm from "components/organization/invitememberform";
import { CreateOrgMemberInviteRequest } from "shared/api/generated/data-contracts";
import InviteMembers from "views/createorganization/components/InviteMembers";
import MemberManager from "components/organization/membermanager/MemberManager";
import TabList from "components/tablist";
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
