import { FlexRowRight, HorizontalSpacer } from "components/globals";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import TextInput from "components/textinput";
import SectionArea from "components/sectionarea";
import StandardButton from "components/buttons";
import ErrorBar from "components/errorbar";

const UserMetaForm: React.FunctionComponent = () => {
  const [displayName, setDisplayName] = useState("");
  const [err, setErr] = useState("");

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

  return (
    <SectionArea width={600} loading={isUserLoading}>
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
  );
};

export default UserMetaForm;
