import {
  FlexCol,
  FlexColCenter,
  FlexRow,
  FlexRowRight,
  H2,
  HorizontalSpacer,
  P,
} from "components/globals";
import React, { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";
import StandardButton from "components/buttons";
import theme from "shared/theme";
import { css } from "styled-components";
import { AppWrapper } from "components/appwrapper";
import ErrorBar from "components/errorbar";
import SectionAreaWithLogo from "components/sectionareawithlogo";

const VerifyEmailPromptView: React.FunctionComponent = () => {
  const [success, setSuccess] = useState(false);
  const [err, setErr] = useState("");

  const resendMutation = useMutation(api.resendVerificationEmail, {
    mutationKey: ["resend_verification_email"],
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

  const renderContents = () => {
    return (
      <SectionAreaWithLogo width={400}>
        <HorizontalSpacer spacepixels={18} />
        <FlexColCenter>
          <H2>Verify your Email</H2>
        </FlexColCenter>
        <HorizontalSpacer
          spacepixels={10}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <HorizontalSpacer spacepixels={30} />
        {success && (
          <P>
            The verification email has been resent. Remember to check your spam
            folder.
          </P>
        )}
        {!success && (
          <P>
            Your email must be verified in order to use Hatchet. Check your
            inbox for an email verification link.
          </P>
        )}
        <HorizontalSpacer spacepixels={30} />
        {err && <ErrorBar text={err} />}
        {err && <HorizontalSpacer spacepixels={12} />}
        <FlexRowRight>
          <StandardButton
            label="Resend Email"
            material_icon="refresh"
            icon_side="right"
            on_click={() => {
              setErr("");
              resendMutation.mutate({});
            }}
            margin={"0"}
            is_loading={resendMutation.isLoading}
          />
        </FlexRowRight>
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

export default VerifyEmailPromptView;
