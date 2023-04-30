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
