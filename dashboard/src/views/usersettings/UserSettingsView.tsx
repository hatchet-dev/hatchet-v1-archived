import {
  FlexColCenter,
  H1,
  H2,
  HorizontalSpacer,
  P,
  SmallSpan,
  FlexColScroll,
} from "hatchet-components";
import React from "react";

import UserMetaForm from "./components/UserMetaForm";
import UserOrgs from "./components/UserOrgs";
import ResetPasswordForm from "./components/ResetPasswordForm";
import { css } from "styled-components";
import theme from "shared/theme";
import DeleteUserForm from "./components/DeleteUserForm";

const UserSettingsView: React.FunctionComponent = () => {
  return (
    <FlexColCenter height={"100%"}>
      <FlexColScroll width="100%" maxWidth="640px" height={"100%"}>
        <H1>Profile</H1>
        <HorizontalSpacer spacepixels={14} />
        <P>Manage your Hatchet profile.</P>
        <HorizontalSpacer
          spacepixels={80}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <H2>Display Name</H2>
        <HorizontalSpacer spacepixels={14} />
        <SmallSpan>
          Your display name is what is shown to other members in your
          organizations and teams. You cannot change your email address without
          contacting an admin.
        </SmallSpan>
        <HorizontalSpacer spacepixels={16} />
        <UserMetaForm />
        <HorizontalSpacer
          spacepixels={80}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <H2>Organizations</H2>
        <HorizontalSpacer spacepixels={14} />
        <SmallSpan>
          All organizations that you're a member of. Note that you cannot leave
          an organization if you are an owner. You must either delete the
          organization or transfer ownership to another member from organization
          settings.
        </SmallSpan>
        <HorizontalSpacer spacepixels={14} />
        <UserOrgs />
        <HorizontalSpacer
          spacepixels={80}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <H2>Reset Password</H2>
        <HorizontalSpacer spacepixels={8} />
        <SmallSpan>Reset your password.</SmallSpan>
        <HorizontalSpacer spacepixels={16} />
        <ResetPasswordForm />
        <HorizontalSpacer
          spacepixels={80}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <H2>Delete User</H2>
        <HorizontalSpacer spacepixels={16} />
        <DeleteUserForm />
      </FlexColScroll>
    </FlexColCenter>
  );
};

export default UserSettingsView;
