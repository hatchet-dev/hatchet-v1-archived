import {
  FlexCol,
  FlexColCenter,
  FlexRow,
  FlexRowRight,
  H2,
  HorizontalSpacer,
  SmallSpan,
  StyledSmallLink,
  TextInput,
  StandardButton,
  AppWrapper,
  ErrorBar,
  SectionAreaWithLogo,
} from "@hatchet-dev/hatchet-components";
import React, { useCallback, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";
import theme from "shared/theme";
import { css } from "styled-components";

const RegisterView: React.FunctionComponent = () => {
  const [displayName, setDisplayName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState("");
  const history = useHistory();

  const { mutate, isLoading } = useMutation(api.createUser, {
    mutationKey: ["create_user"],
    onSuccess: (data) => {
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
    [displayName, email, password]
  );

  useEffect(() => {
    window.addEventListener("keydown", handleKeyPress);

    return () => {
      window.removeEventListener("keydown", handleKeyPress);
    };
  }, [handleKeyPress]);

  const submit = () => {
    setErr("");

    if (displayName != "" && email != "" && password != "") {
      mutate({
        display_name: displayName,
        email: email,
        password: password,
      });
    }
  };

  return (
    <AppWrapper>
      <FlexRow>
        <FlexCol>
          <SectionAreaWithLogo width="400px">
            <HorizontalSpacer spacepixels={18} />
            <FlexColCenter>
              <H2>Create a Hatchet Account</H2>
            </FlexColCenter>
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
                label="Create Account"
                material_icon="chevron_right"
                icon_side="right"
                on_click={() => {
                  submit();
                }}
                margin={"0"}
                disabled={displayName == "" || email == "" || password == ""}
                is_loading={isLoading}
              />
            </FlexRowRight>
          </SectionAreaWithLogo>
          <HorizontalSpacer spacepixels={20} />
          <FlexRowRight>
            <SmallSpan>
              Already have an account?{" "}
              <StyledSmallLink to="/login">Sign in here.</StyledSmallLink>
            </SmallSpan>
          </FlexRowRight>
        </FlexCol>
      </FlexRow>
    </AppWrapper>
  );
};

export default RegisterView;
