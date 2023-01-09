import Breadcrumbs from "components/breadcrumbs";
import {
  FlexCol,
  FlexColCenter,
  FlexRowRight,
  Grid,
  H1,
  H2,
  HorizontalSpacer,
  P,
  Span,
} from "components/globals";
import { GridCard } from "components/gridcard";
import Example from "components/heirarchygraph";
import Paginator from "components/paginator";
import RunsList from "components/runslist";
import Table from "components/table";
import TabList from "components/tablist";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";

import { useHistory } from "react-router-dom";
import api from "shared/api";
import TextInput from "components/textinput";
import SectionArea from "components/sectionarea";
import StandardButton from "components/buttons";
import Spinner from "components/loaders";
import SectionCard from "components/sectioncard";
import ErrorBar from "components/errorbar";
import OrgList from "components/orglist";
import UserMetaForm from "./components/UserMetaForm";
import UserOrgs from "./components/UserOrgs";
import { CenteredContainer } from "components/viewwrapper";

const UserSettingsView: React.FunctionComponent = () => {
  return (
    <FlexColCenter>
      <FlexCol>
        <H1>Profile</H1>
        <HorizontalSpacer spacepixels={16} />
        <H2>Display Name</H2>
        <HorizontalSpacer spacepixels={16} />
        <UserMetaForm />
        <HorizontalSpacer spacepixels={20} />
        <H2>Organizations</H2>
        <HorizontalSpacer spacepixels={16} />
        <UserOrgs />
        <HorizontalSpacer spacepixels={16} />
        <H2>Reset Password</H2>
        <HorizontalSpacer spacepixels={16} />
        <H2>Delete User</H2>
        <HorizontalSpacer spacepixels={16} />
      </FlexCol>
    </FlexColCenter>
  );
};

export default UserSettingsView;
