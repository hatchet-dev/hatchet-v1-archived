import {
  FlexRow,
  HorizontalSpacer,
  P,
  StyledClickableP,
} from "components/globals";
import React, { useState } from "react";
import { useMutation } from "@tanstack/react-query";

import { useHistory } from "react-router-dom";
import api from "shared/api";
import TextInput from "components/textinput";
import SectionArea from "components/sectionarea";
import StandardButton from "components/buttons";
import ErrorBar from "components/errorbar";

const UserMetaForm: React.FunctionComponent = () => {
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [successText, setSuccessText] = useState("");
  const [err, setErr] = useState("");

  const resetPasswordMutation = useMutation(api.resetPasswordManual, {
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

  const sendEmailMutation = useMutation(api.resendVerificationEmail, {
    onSuccess: (data) => {
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
    sendEmailMutation.mutate({});
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
          Send verification email instead.
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
