import {
  H2,
  HorizontalSpacer,
  P,
  FlexColScroll,
  FlexColCenter,
  FlexCol,
} from "hatchet-components";
import React, { useState } from "react";
import { css } from "styled-components";
import theme from "shared/theme";
import OrganizationMetaForm from "./components/OrganizationMetaForm";
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { CreateOrgMemberInviteRequest } from "shared/api/generated/data-contracts";
import MemberManager from "components/organization/membermanager/MemberManager";
import DeleteOrganizationForm from "./components/DeleteOrganizationForm";
import { useHistory } from "react-router-dom";

const OrganizationSettingsView: React.FunctionComponent = () => {
  const [currOrg, setCurrOrg] = useAtom(currOrgAtom);
  const [err, setErr] = useState("");

  const { data, isLoading, refetch } = useQuery({
    queryKey: ["current_organization_members", currOrg.id],
    queryFn: async () => {
      const res = await api.listOrgMembers(currOrg.id);
      return res;
    },
    retry: false,
  });

  const mutation = useMutation({
    mutationKey: ["create_organization_invite", currOrg.id],
    mutationFn: async (invite: CreateOrgMemberInviteRequest) => {
      const res = await api.createOrgMemberInvite(currOrg.id, invite);
      return res;
    },
    onSuccess: (data) => {
      setErr("");
      refetch();
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  return (
    <FlexColCenter height={"100%"}>
      <FlexColScroll width="100%" maxWidth="840px" height={"100%"}>
        <H2>{currOrg.display_name} Settings</H2>
        <HorizontalSpacer spacepixels={14} />
        <P>Manage settings for the {currOrg.display_name} organization.</P>
        <HorizontalSpacer
          spacepixels={80}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <H2>Organization Name</H2>
        <HorizontalSpacer spacepixels={14} />
        <OrganizationMetaForm />
        <HorizontalSpacer
          spacepixels={80}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <FlexCol height="600px">
          <H2>Manage Organization Members</H2>
          <HorizontalSpacer spacepixels={14} />
          <MemberManager can_remove={true} header_level="H3" />
        </FlexCol>
        <HorizontalSpacer spacepixels={14} />
        {/* <UserOrgs /> */}
        <HorizontalSpacer
          spacepixels={80}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <H2>Delete Organization</H2>
        <HorizontalSpacer spacepixels={16} />
        <DeleteOrganizationForm />
      </FlexColScroll>
    </FlexColCenter>
  );
};

export default OrganizationSettingsView;
