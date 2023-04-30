import {
  FlexCol,
  FlexRowRight,
  Grid,
  H1,
  H2,
  HorizontalSpacer,
  P,
  SmallSpan,
  Span,
  Breadcrumbs,
  GridCard,
  TextInput,
  SectionArea,
  StandardButton,
  ErrorBar,
  CopyCodeline,
} from "hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { CreatePATResponse } from "shared/api/generated/data-contracts";

type Props = {
  post_create: () => void;
};

const CreatePATForm: React.FunctionComponent<Props> = ({ post_create }) => {
  const [displayName, setDisplayName] = useState("");
  const [createdToken, setCreatedToken] = useState<CreatePATResponse>();
  const [err, setErr] = useState("");

  const { mutate, isLoading } = useMutation(api.createPersonalAccessToken, {
    mutationKey: ["create_pat"],
    onSuccess: (data) => {
      setCreatedToken(data?.data);
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  const submit = () => {
    if (displayName != "") {
      mutate({
        display_name: displayName,
      });
    }
  };

  if (createdToken) {
    return (
      <SectionArea>
        <HorizontalSpacer spacepixels={6} />
        <SmallSpan>
          Your personal access token has been created. The token will only be
          shown once, so make sure you store it somewhere safe.
        </SmallSpan>
        <HorizontalSpacer spacepixels={20} />
        <CopyCodeline code_line={createdToken.token} />
        <HorizontalSpacer spacepixels={30} />
        <FlexRowRight>
          <StandardButton
            label="Done"
            material_icon="chevron_right"
            icon_side="right"
            on_click={() => {
              post_create();
            }}
            margin={"0"}
          />
        </FlexRowRight>
      </SectionArea>
    );
  }

  return (
    <SectionArea>
      <TextInput
        placeholder="My Token"
        label="Personal access token name"
        type="text"
        width="400px"
        on_change={(val) => {
          setDisplayName(val);
        }}
      />
      <HorizontalSpacer spacepixels={20} />
      {err && <ErrorBar text={err} />}
      <HorizontalSpacer spacepixels={30} />
      <FlexRowRight>
        <StandardButton
          label="Create"
          material_icon="chevron_right"
          icon_side="right"
          on_click={() => {
            submit();
          }}
          disabled={displayName == ""}
          margin={"0"}
          is_loading={isLoading}
        />
      </FlexRowRight>
    </SectionArea>
  );
};

export default CreatePATForm;
