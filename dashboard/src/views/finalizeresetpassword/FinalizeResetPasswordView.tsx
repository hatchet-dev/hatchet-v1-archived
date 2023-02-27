import {
  FlexCol,
  FlexColCenter,
  FlexRow,
  FlexRowRight,
  H2,
  HorizontalSpacer,
  P,
  StyledLink,
  TextInput,
  StandardButton,
  AppWrapper,
  ErrorBar,
  SectionAreaWithLogo,
} from "@hatchet-dev/hatchet-components";
import React, { useCallback, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import theme from "shared/theme";
import { css } from "styled-components";
import useQueryParam from "shared/hooks/usequeryparam";

const FinalizeResetPasswordView: React.FunctionComponent = () => {
  const [newPassword, setNewPassword] = useState("");
  const [success, setSuccess] = useState(false);
  const [err, setErr] = useState("");
  const history = useHistory();

  const query = useQueryParam();

  const tokenId = query.get("token_id");
  const token = query.get("token");
  const email = query.get("email");

  useEffect(() => {
    if (!tokenId || !token || !email) {
      history.push("/login");
    }
  }, [tokenId, token, email]);

  const { isLoading, isError } = useQuery({
    queryKey: ["reset_password_email_verify"],
    queryFn: async () => {
      const res = await api.resetPasswordEmailVerify({
        email: email,
        token: token,
        token_id: tokenId,
      });

      return res;
    },
    retry: false,
  });

  const resetMutation = useMutation(api.resetPasswordEmailFinalize, {
    mutationKey: ["reset_password_email"],
    onSuccess: (data) => {
      setSuccess(true);
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  const handleKeyPress = useCallback(
    (e: any) => {
      e.key === "Enter" && submit();
    },
    [newPassword]
  );

  useEffect(() => {
    window.addEventListener("keydown", handleKeyPress);

    return () => {
      window.removeEventListener("keydown", handleKeyPress);
    };
  }, [handleKeyPress]);

  const submit = () => {
    setErr("");

    if (newPassword != "") {
      resetMutation.mutate({
        token: token,
        token_id: tokenId,
        email: email,
        new_password: newPassword,
      });
    }
  };

  const renderInnerContents = () => {
    if (success) {
      return (
        <>
          <P>Successfully reset your password!</P>
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
          placeholder="New Password"
          label="Your new password"
          type="password"
          width="100%"
          on_change={(val) => {
            setNewPassword(val);
          }}
        />
        {err && <HorizontalSpacer spacepixels={20} />}
        {err && <ErrorBar text={err} />}
        <HorizontalSpacer spacepixels={30} />
        <FlexRowRight>
          <StandardButton
            label="Reset Password"
            material_icon="chevron_right"
            icon_side="right"
            on_click={() => {
              submit();
            }}
            margin={"0"}
            disabled={newPassword == ""}
            is_loading={resetMutation.isLoading}
          />
        </FlexRowRight>
      </>
    );
  };

  const renderContents = () => {
    if (isError) {
      return (
        <SectionAreaWithLogo width="400px" loading={isLoading}>
          <HorizontalSpacer spacepixels={18} />
          <P>
            This password reset link is no longer valid.{" "}
            <StyledLink to="/reset_password/initiate">
              Generate new link?
            </StyledLink>
          </P>
        </SectionAreaWithLogo>
      );
    }

    return (
      <SectionAreaWithLogo width="400px" loading={isLoading}>
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
        <HorizontalSpacer spacepixels={18} />
        {renderInnerContents()}
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

export default FinalizeResetPasswordView;
