import { FlexCol, FlexRow, H1, HorizontalSpacer } from "components/globals";
import React from "react";
import { useHistory, useParams, useLocation } from "react-router-dom";
import github from "assets/github.png";
import branch from "assets/branch.png";
import { AppWrapper } from "components/appwrapper";
import NameOrganization from "./components/NameOrganization";
import InviteMembers from "./components/InviteMembers";

const CreateOrganizationView: React.FunctionComponent = () => {
  const location = useLocation();

  const renderForm = () => {
    switch (location.pathname) {
      case "/organization/create":
        return <NameOrganization />;
      case "/organization/create/invite_members":
        return <InviteMembers />;
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
