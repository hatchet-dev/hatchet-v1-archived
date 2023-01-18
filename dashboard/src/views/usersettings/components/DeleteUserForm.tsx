import {
  FlexRowRight,
  HorizontalSpacer,
  P,
  SmallSpan,
  SectionArea,
  StandardButton,
  ErrorBar,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";

import { useHistory } from "react-router-dom";
import api from "shared/api";

const DeleteUserForm: React.FunctionComponent = () => {
  const [success, setSuccess] = useState(false);
  const [err, setErr] = useState("");

  const orgQuery = useQuery({
    queryKey: ["current_user_organizations"],
    queryFn: async () => {
      const res = await api.listUserOrganizations();
      return res;
    },
    retry: false,
  });

  const currentUserQuery = useQuery({
    queryKey: ["current_user"],
    queryFn: async () => {
      const res = await api.getCurrentUser();
      return res;
    },
    retry: false,
  });

  const { mutate, isLoading } = useMutation(api.deleteCurrentUser, {
    mutationKey: ["delete_current_user"],
    onSuccess: (data) => {
      setSuccess(true);
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  const submit = () => {
    mutate({});
  };

  if (orgQuery.isLoading || currentUserQuery.isLoading) {
    return <SectionArea width="600px" loading={true} />;
  }

  if (success) {
    return (
      <SectionArea width="600px">
        <P>
          You successfully deleted your user. You will be automatically logged
          out within 10 seconds.
        </P>
      </SectionArea>
    );
  }

  const userEmail = currentUserQuery.data?.data.email;

  const ownerOrgs = orgQuery.data?.data?.rows.filter((row) => {
    return userEmail == row.owner.email;
  });

  const renderMessage = () => {
    if (ownerOrgs.length > 0) {
      return (
        <SmallSpan>
          You are currently the owner of the following organizations:{" "}
          {ownerOrgs.map((row) => row.display_name).join(", ")}. Before deleting
          your user, you must delete these organizations or transfer ownership
          to another member.
        </SmallSpan>
      );
    }

    return <SmallSpan>This operation cannot be undone.</SmallSpan>;
  };

  return (
    <SectionArea width="600px">
      {renderMessage()}
      {err && <HorizontalSpacer spacepixels={12} />}
      {err && <ErrorBar text={err} />}
      <HorizontalSpacer spacepixels={12} />
      <FlexRowRight>
        <StandardButton
          label="Delete User"
          material_icon="delete"
          icon_side="right"
          on_click={() => {
            submit();
          }}
          margin={"0"}
          is_loading={isLoading}
          disabled={ownerOrgs.length > 0}
        />
      </FlexRowRight>
    </SectionArea>
  );
};

export default DeleteUserForm;
