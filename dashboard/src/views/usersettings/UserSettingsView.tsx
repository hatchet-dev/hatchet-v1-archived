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

const TabOptions = ["Runs", "Resource Explorer", "Configuration", "Settings"];

const UserSettingsView: React.FunctionComponent = () => {
  const [selectedTab, setSelectedTab] = useState(TabOptions[0]);
  const [displayName, setDisplayName] = useState("");
  const [err, setErr] = useState("");

  let history = useHistory();

  const query = useQuery({
    queryKey: ["current_user"],
    queryFn: async () => {
      const res = await api.getCurrentUser();
      return res;
    },
    retry: false,
  });

  const isUserLoading = query.isLoading;
  const data = query.data;

  const { mutate, isLoading } = useMutation(api.updateCurrentUser, {
    onSuccess: (data) => {
      query.refetch();
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  const submit = () => {
    if (displayName != "" && displayName != data?.data?.display_name) {
      mutate({
        display_name: displayName,
      });
    }
  };

  if (isUserLoading) {
    return <Spinner></Spinner>;
  }

  return (
    <>
      <H1>Profile</H1>
      <HorizontalSpacer spacepixels={16} />
      <H2>Display Name</H2>
      <HorizontalSpacer spacepixels={16} />
      <SectionArea width={600}>
        <TextInput
          placeholder="Hatchet User"
          initial_value={data?.data?.display_name}
          label="Your name"
          type="text"
          width="400px"
          on_change={(val) => {
            setDisplayName(val);
          }}
        />
        <HorizontalSpacer spacepixels={20} />
        <TextInput
          placeholder="you@example.com"
          initial_value={data?.data?.email}
          label="Your email"
          type="text"
          width="400px"
          disabled={true}
        />
        <HorizontalSpacer spacepixels={30} />
        {err && <ErrorBar text={err} />}
        <HorizontalSpacer spacepixels={30} />
        <FlexRowRight>
          <StandardButton
            label="Update"
            material_icon="chevron_right"
            icon_side="right"
            on_click={() => {
              submit();
            }}
            disabled={
              displayName == "" || displayName == data?.data?.display_name
            }
            margin={"0"}
            is_loading={isLoading}
          />
        </FlexRowRight>
      </SectionArea>
      <HorizontalSpacer spacepixels={20} />
      <H2>Organizations</H2>
      <HorizontalSpacer spacepixels={16} />
      <H2>Reset Password</H2>
      <HorizontalSpacer spacepixels={16} />
      <H2>Delete User</H2>
      <HorizontalSpacer spacepixels={16} />
    </>
  );
};

export default UserSettingsView;
