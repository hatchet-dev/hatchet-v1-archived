import {
  FlexRow,
  HorizontalSpacer,
  P,
  StyledClickableP,
} from "components/globals";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";

import { useHistory } from "react-router-dom";
import api from "shared/api";
import TextInput from "components/textinput";
import SectionArea from "components/sectionarea";
import StandardButton from "components/buttons";
import ErrorBar from "components/errorbar";
import { ResetPasswordEmailRequest } from "shared/api/generated/data-contracts";

const UserMetaForm: React.FunctionComponent = () => {
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [successText, setSuccessText] = useState("");
  const [err, setErr] = useState("");

  const currentUserQuery = useQuery({
    queryKey: ["current_user"],
    queryFn: async () => {
      const res = await api.getCurrentUser();
      return res;
    },
    retry: false,
  });

  const resetPasswordMutation = useMutation(api.resetPasswordManual, {
    mutationKey: ["reset_password_manual"],
    onSuccess: (data) => {
      setSuccessText("You successfully changed your password.");
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  const sendEmailMutation = useMutation({
    mutationFn: (resetReq: ResetPasswordEmailRequest) => {
      return api.resetPasswordEmail(resetReq, {
        credentials: "omit",
      });
    },
    onSuccess: () => {
      setSuccessText(
        "Reset password email sent. Remember to check your spam folder."
      );
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  const submit = () => {
    if (oldPassword != "" && newPassword != "") {
      resetPasswordMutation.mutate({
        old_password: oldPassword,
        new_password: newPassword,
      });
    }
  };

  const sendVerificationEmail = () => {
    const email = currentUserQuery.data?.data?.email;

    if (email) {
      sendEmailMutation.mutate({
        email: email,
      });
    }
  };

  if (successText) {
    return (
      <SectionArea width={600}>
        <P>{successText}</P>
      </SectionArea>
    );
  }

  return (
    <SectionArea width={600}>
      <TextInput
        placeholder="Old Password"
        label="Your old password"
        type="password"
        width="400px"
        on_change={(val) => {
          setOldPassword(val);
        }}
      />
      <HorizontalSpacer spacepixels={20} />
      <TextInput
        placeholder="New Password"
        label="Your new password"
        type="password"
        width="400px"
        on_change={(val) => {
          setNewPassword(val);
        }}
      />
      <HorizontalSpacer spacepixels={30} />
      {err && <ErrorBar text={err} />}
      <HorizontalSpacer spacepixels={30} />
      <FlexRow>
        <StyledClickableP onClick={() => sendVerificationEmail()}>
          Send reset password email instead.
        </StyledClickableP>
        <StandardButton
          label="Reset Password"
          material_icon="chevron_right"
          icon_side="right"
          on_click={() => {
            submit();
          }}
          disabled={oldPassword == "" || newPassword == ""}
          margin={"0"}
          is_loading={
            resetPasswordMutation.isLoading || sendEmailMutation.isLoading
          }
        />
      </FlexRow>
    </SectionArea>
  );
};

export default UserMetaForm;
