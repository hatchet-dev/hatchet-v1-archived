import {
  H2,
  HorizontalSpacer,
  P,
  FlexColScroll,
  FlexColCenter,
  FlexCol,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { css } from "styled-components";
import theme from "shared/theme";
import { currOrgAtom, currTeamAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { CreateOrgMemberInviteRequest } from "shared/api/generated/data-contracts";
import MemberManager from "components/organization/membermanager/MemberManager";
import { useHistory } from "react-router-dom";
import TeamWithMembers from "components/team/teamwithmembers";
import TeamMetaForm from "./components/TeamMetaForm";
import DeleteTeamForm from "./components/DeleteTeamForm";

const TeamSettingsView: React.FunctionComponent = () => {
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);
  const [err, setErr] = useState("");

  // const { data, isLoading, refetch } = useQuery({
  //   queryKey: ["current_organization_members", currOrg.id],
  //   queryFn: async () => {
  //     const res = await api.listOrgMembers(currOrg.id);
  //     return res;
  //   },
  //   retry: false,
  // });

  // const mutation = useMutation({
  //   mutationKey: ["create_organization_invite", currOrg.id],
  //   mutationFn: (invite: CreateOrgMemberInviteRequest) => {
  //     return api.createOrgMemberInvite(currOrg.id, invite);
  //   },
  //   onSuccess: (data) => {
  //     setErr("");
  //     refetch();
  //   },
  //   onError: (err: any) => {
  //     if (!err.error.errors || err.error.errors.length == 0) {
  //       setErr("An unexpected error occurred. Please try again.");
  //     }

  //     setErr(err.error.errors[0].description);
  //   },
  // });

  return (
    <FlexColCenter height={"100%"}>
      <FlexColScroll width="100%" maxWidth="840px" height={"100%"}>
        <H2>{currTeam.display_name} Settings</H2>
        <HorizontalSpacer spacepixels={14} />
        <P>Manage settings for the {currTeam.display_name} team.</P>
        <HorizontalSpacer
          spacepixels={80}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <H2>Team Name</H2>
        <HorizontalSpacer spacepixels={14} />
        <TeamMetaForm />
        <HorizontalSpacer
          spacepixels={80}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <FlexCol height="600px">
          <H2>Manage Team Members</H2>
          <HorizontalSpacer spacepixels={14} />
          <TeamWithMembers team={currTeam} collapsible={false} />
        </FlexCol>
        <HorizontalSpacer
          spacepixels={80}
          overrides={css({
            borderBottom: theme.line.thick,
          }).toString()}
        />
        <H2>Delete Team</H2>
        <HorizontalSpacer spacepixels={16} />
        <DeleteTeamForm />
      </FlexColScroll>
    </FlexColCenter>
  );
};

export default TeamSettingsView;
