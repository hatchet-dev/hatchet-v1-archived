import Breadcrumbs from "components/breadcrumbs";
import {
  FlexCol,
  FlexRow,
  FlexRowRight,
  H1,
  H2,
  HorizontalSpacer,
  P,
} from "components/globals";
import Selector from "components/selector";
import React, { useState } from "react";
import { useHistory, useParams, useLocation } from "react-router-dom";
import gitRepository from "assets/git_repository.png";
import github from "assets/github.png";
import branch from "assets/branch.png";
import TextInput from "components/textinput";
import FormArea from "components/sectionarea";
import StandardButton from "components/buttons";
import { css } from "styled-components";
import theme from "shared/theme";
import SectionArea from "components/sectionarea";
import { AppWrapper } from "components/appwrapper";
import { ViewWrapper } from "components/viewwrapper";
import NameOrganization from "./components/NameOrganization";
import InviteMembers from "./components/InviteMembers";

const options = [
  {
    icon: github,
    label: "hatchet-dev/hatchet",
    value: "hatchet-dev/hatchet",
  },
  {
    icon: github,
    label: "hatchet-dev/hatchet-2",
    value: "hatchet-dev/hatchet-2",
  },
];

const branchOptions = [
  {
    icon: branch,
    label: "master",
    value: "master",
  },
  {
    icon: branch,
    label: "belanger/feat-1",
    value: "belanger/feat-1",
  },
];

const CreateOrganizationView: React.FunctionComponent = () => {
  const history = useHistory();
  const location = useLocation();
  const { step } = useParams<{ step: string }>();

  const renderForm = () => {
    switch (location.pathname) {
      case "/organizations/create":
        return <NameOrganization />;
      case "/organizations/create/invite_members":
        return <InviteMembers />;
    }
  };

  return (
    <AppWrapper>
      <FlexRow>
        <FlexCol>
          {/* <Breadcrumbs breadcrumbs={getBreadcrumbs()} /> */}
          {/* <HorizontalSpacer spacepixels={12} /> */}
          <H1>Create a New Organization</H1>
          <HorizontalSpacer spacepixels={20} />
          {renderForm()}
          <HorizontalSpacer spacepixels={20} />
          {/* <FlexRowRight>{renderButton()}</FlexRowRight> */}
        </FlexCol>
      </FlexRow>
    </AppWrapper>
  );
};

export default CreateOrganizationView;
