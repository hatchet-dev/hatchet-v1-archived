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

const LoginView: React.FunctionComponent = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const history = useHistory();

  const { mutate, isLoading } = useMutation(api.loginUser, {
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

  const submit = async () => {
    mutate({
      email: email,
      password: password,
    });
  };

  return (
    <ThemeProvider theme={invertedTheme}>
      <AppWrapper>
        <FlexRow>
          <SectionArea width={400} background={"white"}>
            <H2>Login</H2>
            <HorizontalSpacer
              spacepixels={10}
              overrides={css({
                borderBottom: theme.line.thick,
              }).toString()}
            />
            <HorizontalSpacer spacepixels={18} />
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
                label="Login"
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

export default LoginView;
