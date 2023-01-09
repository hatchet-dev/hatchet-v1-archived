import Breadcrumbs from "components/breadcrumbs";
import {
  FlexCol,
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

const TabOptions = ["Runs", "Resource Explorer", "Configuration", "Settings"];

const UserOrgs: React.FunctionComponent = () => {
  const [selectedTab, setSelectedTab] = useState(TabOptions[0]);
  const [displayName, setDisplayName] = useState("");
  const [err, setErr] = useState("");

  let history = useHistory();

  const orgQuery = useQuery({
    queryKey: ["current_user_organizations"],
    queryFn: async () => {
      const res = await api.listUserOrganizations();
      return res;
    },
    retry: false,
  });

  const leaveOrgMutation = useMutation({
    mutationFn: (orgId: string) => {
      return api.leaveOrg(orgId);
    },
    onSuccess: (data) => {
      console.log("GOT ON SUCCESS", data);

      orgQuery.refetch();
    },
    onError: (err: any) => {
      console.log("GOT AN ERROR", err);

      if (!err.error || !err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      console.log("SETTING ERROR", err.error.errors[0].description);

      setErr(err.error.errors[0].description);
    },
  });

  if (orgQuery.isLoading) {
    return <Spinner></Spinner>;
  }

  return (
    <SectionArea width={600}>
      <OrgList
        orgs={orgQuery?.data?.data.rows}
        leave_org={(org) => leaveOrgMutation.mutate(org.id)}
        err={err}
      />
    </SectionArea>
  );
};

export default UserOrgs;
