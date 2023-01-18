import {
  FlexCol,
  FlexColCenter,
  FlexRow,
  FlexRowRight,
  H2,
  HorizontalSpacer,
  P,
  StandardButton,
  AppWrapper,
  ErrorBar,
  SectionAreaWithLogo,
  Spinner,
} from "@hatchet-dev/hatchet-components";
import React, { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import api from "shared/api";
import theme from "shared/theme";
import { css } from "styled-components";
import useQueryParam from "shared/hooks/usequeryparam";

const VerifyEmailView: React.FunctionComponent = () => {
  const [success, setSuccess] = useState(false);
  const [err, setErr] = useState("");
  const history = useHistory();
  const query = useQueryParam();

  const tokenId = query.get("token_id");
  const token = query.get("token");

  useEffect(() => {
    if (!tokenId || !token) {
      history.push("/");
    }
  }, [tokenId, token]);

  const verifyMutation = useMutation(api.verifyEmail, {
    mutationKey: ["verify_email"],
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

  useEffect(() => {
    if (!success) {
      verifyMutation.mutate({
        token: token,
        token_id: tokenId,
      });
    }
  }, []);

  const renderContents = () => {
    return (
      <SectionAreaWithLogo width="400px">
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
            Your email was successfully verified. You will be automatically
            redirected to the Hatchet dashboard.
          </P>
        )}
        {!success && <Spinner />}
        <HorizontalSpacer spacepixels={30} />
        {err && <ErrorBar text={err} />}
        {err && <HorizontalSpacer spacepixels={12} />}
        {success && (
          <FlexRowRight>
            <StandardButton
              label="Go to Dashboard"
              material_icon="chevron_right"
              icon_side="right"
              on_click={() => {
                history.push("/");
              }}
              margin={"0"}
            />
          </FlexRowRight>
        )}
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

export default VerifyEmailView;
