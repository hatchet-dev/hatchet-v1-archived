import {
  FlexRow,
  FlexRowRight,
  H2,
  HorizontalSpacer,
} from "components/globals";
import TextInput from "components/textinput";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";
import StandardButton from "components/buttons";
import SectionArea from "components/sectionarea";
import theme, { invertedTheme } from "shared/theme";
import styled, { css, ThemeProvider } from "styled-components";
import { AppWrapper } from "components/appwrapper";

const RegisterView: React.FunctionComponent = () => {
  const [displayName, setDisplayName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const history = useHistory();

  const { mutate, isLoading } = useMutation(api.createUser, {
    onSuccess: (data) => {
      console.log(data);
      history.push("/");
    },
    onError: (err: any) => {
      console.log("ERR", err.error);
    },
    onSettled: () => {
      //   queryClient.invalidateQueries('login_user');
    },
  });

  const submit = () => {
    mutate({
      display_name: displayName,
      email: email,
      password: password,
    });
  };

  return (
    <ThemeProvider theme={invertedTheme}>
      <AppWrapper>
        <FlexRow>
          <SectionArea width={400} background={"white"}>
            <H2>Create an Account</H2>
            <HorizontalSpacer
              spacepixels={10}
              overrides={css({
                borderBottom: theme.line.thick,
              }).toString()}
            />
            <HorizontalSpacer spacepixels={18} />
            <TextInput
              placeholder="Hatchet User"
              label="Your name"
              type="text"
              width="100%"
              on_change={(val) => {
                setDisplayName(val);
              }}
            />
            <HorizontalSpacer spacepixels={20} />
            <TextInput
              placeholder="you@example.com"
              label="Your email"
              type="text"
              width="100%"
              on_change={(val) => {
                setEmail(val);
              }}
            />
            <HorizontalSpacer spacepixels={20} />
            <TextInput
              placeholder=""
              label="Your password"
              type="password"
              width="100%"
              on_change={(val) => {
                setPassword(val);
              }}
            />
            <HorizontalSpacer spacepixels={30} />
            <FlexRowRight>
              <StandardButton
                label="Create Account"
                material_icon="chevron_right"
                icon_side="right"
                on_click={() => {
                  submit();
                }}
                margin={"0"}
              />
            </FlexRowRight>
          </SectionArea>
        </FlexRow>
      </AppWrapper>
    </ThemeProvider>
  );
};

export default RegisterView;
