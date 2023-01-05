import {
  FlexCol,
  FlexRow,
  FlexRowRight,
  H2,
  HorizontalSpacer,
  SmallSpan,
  Span,
  StyledLink,
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
import Spinner from "components/loaders";

const LoginView: React.FunctionComponent = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState("");
  const history = useHistory();

  const { mutate, isLoading } = useMutation(api.loginUser, {
    onSuccess: (data) => {
      console.log(data);
      history.push("/");
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
    [email, password]
  );

  useEffect(() => {
    window.addEventListener("keydown", handleKeyPress);

    return () => {
      window.removeEventListener("keydown", handleKeyPress);
    };
  }, [handleKeyPress]);

  const submit = () => {
    setErr("");

    if (email != "" && password != "") {
      mutate({
        email: email,
        password: password,
      });
    }
  };

  return (
    <AppWrapper>
      <FlexRow>
        <FlexCol>
          <SectionArea width={400}>
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
              placeholder="Password"
              label="Your password"
              type="password"
              width="100%"
              on_change={(val) => {
                setPassword(val);
              }}
            />
            <HorizontalSpacer spacepixels={30} />
            {err && <ErrorBar text={err} />}
            <HorizontalSpacer spacepixels={30} />
            <FlexRowRight>
              <StandardButton
                label="Login"
                material_icon="chevron_right"
                icon_side="right"
                on_click={() => {
                  submit();
                }}
                disabled={email == "" || password == ""}
                margin={"0"}
                is_loading={isLoading}
              />
            </FlexRowRight>
          </SectionArea>
          <HorizontalSpacer spacepixels={20} />
          <FlexRowRight>
            <SmallSpan>
              Don't have an account?{" "}
              <StyledSmallLink to="/register">Register here.</StyledSmallLink>
            </SmallSpan>
          </FlexRowRight>
        </FlexCol>
      </FlexRow>
    </AppWrapper>
  );
};

export default LoginView;
