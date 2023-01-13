import {
  FlexColCenter,
  H1,
  H2,
  HorizontalSpacer,
  P,
  SmallSpan,
  FlexColScroll,
  FlexCol,
} from "components/globals";
import React, { useEffect, useState } from "react";
import { css } from "styled-components";
import theme from "shared/theme";
import OrganizationMetaForm from "./components/OrganizationMetaForm";
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import MemberList from "components/organization/memberlist";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import InviteMemberForm from "components/organization/invitememberform";
import { CreateOrgMemberInviteRequest } from "shared/api/generated/data-contracts";
import InviteMembers from "views/createorganization/components/InviteMembers";
import MemberManager from "components/organization/membermanager/MemberManager";
import DeleteOrganizationForm from "./components/DeleteOrganizationForm";
import TabList from "components/tablist";
import usePrevious from "shared/hooks/useprevious";
import { useHistory } from "react-router-dom";

const OrganizationSettingsView: React.FunctionComponent = () => {
  const history = useHistory();
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
    mutationFn: (invite: CreateOrgMemberInviteRequest) => {
      return api.createOrgMemberInvite(currOrg.id, invite);
    },
    onSuccess: (data) => {
      setErr("");
      refetch();
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  return (
    <FlexColScroll width="100%" height={"100%"}>
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
  );
};

export default OrganizationSettingsView;
