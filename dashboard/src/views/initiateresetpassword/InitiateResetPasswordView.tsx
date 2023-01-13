import {
  FlexCol,
  FlexColCenter,
  FlexRow,
  FlexRowRight,
  H2,
  HorizontalSpacer,
  P,
} from "components/globals";
import TextInput from "components/textinput";
import React, { useCallback, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";
import StandardButton from "components/buttons";
import theme from "shared/theme";
import { css } from "styled-components";
import { AppWrapper } from "components/appwrapper";
import ErrorBar from "components/errorbar";
import SectionAreaWithLogo from "components/sectionareawithlogo";

const InitiateResetPasswordView: React.FunctionComponent = () => {
  const [email, setEmail] = useState("");
  const [success, setSuccess] = useState(false);
  const [err, setErr] = useState("");
  const history = useHistory();

  const { mutate, isLoading } = useMutation(api.resetPasswordEmail, {
    mutationKey: ["reset_password_email"],
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

  const handleKeyPress = useCallback(
    (e: any) => {
      e.key === "Enter" && submit();
    },
    [email]
  );

  useEffect(() => {
    window.addEventListener("keydown", handleKeyPress);

    return () => {
      window.removeEventListener("keydown", handleKeyPress);
    };
  }, [handleKeyPress]);

  const submit = () => {
    setErr("");

    if (email != "") {
      mutate({
        email: email,
      });
    }
  };

  const renderInnerForm = () => {
    if (success) {
      return (
        <>
          <P>
            If your email address exists in our database, you will receive a
            password reset link at your email address in a few minutes.
          </P>
          <HorizontalSpacer spacepixels={30} />
          <FlexRowRight>
            <StandardButton
              label="Back to Login"
              material_icon="reply"
              icon_side="right"
              on_click={() => {
                history.push("/login");
              }}
              margin={"0"}
            />
          </FlexRowRight>
        </>
      );
    }

    return (
      <>
        <TextInput
          placeholder="you@example.com"
          label="Your email"
          type="text"
          width="100%"
          on_change={(val) => {
            setEmail(val);
          }}
        />
        {err && <HorizontalSpacer spacepixels={20} />}
        {err && <ErrorBar text={err} />}
        <HorizontalSpacer spacepixels={30} />
        <FlexRowRight>
          <StandardButton
            label="Send Reset Email"
            material_icon="chevron_right"
            icon_side="right"
            on_click={() => {
              submit();
            }}
            margin={"0"}
            disabled={email == ""}
            is_loading={isLoading}
          />
        </FlexRowRight>
      </>
    );
  };

  const renderContents = () => {
    return (
      <SectionAreaWithLogo width={400}>
        <HorizontalSpacer spacepixels={18} />
        <FlexColCenter>
          <H2>Reset your Password</H2>
        </FlexColCenter>
        <HorizontalSpacer
          spacepixels={10}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={30} />
        {renderInnerForm()}
      </SectionAreaWithLogo>
    );
  };

  return (
    <AppWrapper>
      <FlexRow>
        <FlexCol>{renderContents()}</FlexCol>
      </FlexRow>
    </AppWrapper>
  );
};

export default InitiateResetPasswordView;
