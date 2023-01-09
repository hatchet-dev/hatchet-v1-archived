import Breadcrumbs from "components/breadcrumbs";
import {
  FlexCol,
  FlexRow,
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

const UserPATsView: React.FunctionComponent = () => {
  const history = useHistory();

  const columns = [
    {
      Header: "Name",
      accessor: "name",
    },
    {
      Header: "Created",
      accessor: "created_at",
    },
  ];

  const data = [
    {
      id: "1111",
      name: "cli-token-1",
      created_at: "10 days ago",
    },
    {
      id: "2222",
      name: "cli-token-2",
      created_at: "15 days ago",
    },
    {
      id: "3333",
      name: "cli-token-3",
      created_at: "20 days ago",
    },
  ];

  return (
    <>
      <H1>Personal Access Tokens</H1>
      <HorizontalSpacer spacepixels={16} />
      <FlexRowRight>
        <StandardButton
          label="Create new PAT"
          material_icon="add"
          on_click={() => history.push("/user/settings/pats/create")}
        />
      </FlexRowRight>
      <HorizontalSpacer spacepixels={12} />
      <Table columns={columns} data={data} onRowClick={() => {}} />
      <HorizontalSpacer spacepixels={20} />
    </>
  );
};

export default UserPATsView;
