import {
  FlexCol,
  FlexRow,
  H1,
  HorizontalSpacer,
  AppWrapper,
} from "@hatchet-dev/hatchet-components";
import React from "react";
import { useLocation } from "react-router-dom";
import NameOrganization from "./components/NameOrganization";
import InviteMembers from "./components/InviteMembers";
import CreateTeams from "./components/CreateTeams";

const CreateOrganizationView: React.FunctionComponent = () => {
  const location = useLocation();

  const renderForm = () => {
    switch (location.pathname) {
      case "/organization/create":
        return <NameOrganization />;
      case "/organization/create/invite_members":
        return <InviteMembers />;
      case "/organization/create/create_teams":
        return <CreateTeams />;
    }
  };

  return (
    <AppWrapper>
      <FlexRow>
        <FlexCol>
          <H1>Create a New Organization</H1>
          <HorizontalSpacer spacepixels={20} />
          {renderForm()}
          <HorizontalSpacer spacepixels={20} />
        </FlexCol>
      </FlexRow>
    </AppWrapper>
  );
};

export default CreateOrganizationView;
