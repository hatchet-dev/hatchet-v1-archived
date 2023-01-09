import {
  FlexCol,
  FlexRow,
  FlexRowRight,
  H2,
  HorizontalSpacer,
  SmallSpan,
  StyledSmallLink,
} from "components/globals";
import TextInput from "components/textinput";
import React, { useCallback, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";
import StandardButton from "components/buttons";
import SectionArea from "components/sectionarea";
import theme, { invertedTheme } from "shared/theme";
import styled, { css, ThemeProvider } from "styled-components";
import { AppWrapper } from "components/appwrapper";
import ErrorBar from "components/errorbar";
import { atom, useAtom } from "jotai";
import { currOrgAtom } from "shared/atoms/atoms";

const NameOrganization: React.FunctionComponent = () => {
  const [displayName, setDisplayName] = useState("");
  const [, setCurrOrgId] = useAtom(currOrgAtom);
  const [err, setErr] = useState("");
  const history = useHistory();

  const { mutate, isLoading } = useMutation(api.createOrganization, {
    onSuccess: (data) => {
      setCurrOrgId(data.data.id);

      history.push("/organizations/create/invite_members");
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  const handleKeyPress = useCallback(
    (e: any) => {
      e.key === "Enter" && submit();
    },
    [displayName]
  );

  useEffect(() => {
    window.addEventListener("keydown", handleKeyPress);

    return () => {
      window.removeEventListener("keydown", handleKeyPress);
    };
  }, [handleKeyPress]);

  const submit = () => {
    setErr("");

    if (displayName != "") {
      mutate({
        display_name: displayName,
      });
    }
  };

  return (
    <SectionArea width={600}>
      <HorizontalSpacer spacepixels={24} />
      <H2>Name your Organization</H2>
      <HorizontalSpacer spacepixels={24} />
      <TextInput
        placeholder="My Organization"
        label="Your organization"
        type="text"
        width="100%"
        on_change={(val) => {
          setDisplayName(val);
        }}
      />
      <HorizontalSpacer spacepixels={30} />
      {err && <ErrorBar text={err} />}
      <HorizontalSpacer spacepixels={30} />
      <FlexRowRight>
        <StandardButton
          label="Create Organization"
          material_icon="chevron_right"
          icon_side="right"
          on_click={() => {
            submit();
          }}
          margin={"0"}
          disabled={displayName == ""}
          is_loading={isLoading}
        />
      </FlexRowRight>
    </SectionArea>
  );
};

export default NameOrganization;
