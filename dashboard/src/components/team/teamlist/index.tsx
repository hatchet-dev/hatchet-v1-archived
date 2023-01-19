import {
  FlexColScroll,
  MaterialIcon,
  StandardButton,
} from "@hatchet-dev/hatchet-components";
import React from "react";
import {
  OrganizationMemberSanitized,
  Team,
} from "shared/api/generated/data-contracts";
import { capitalize } from "shared/utils";
import TeamWithMembers from "../teamwithmembers";
import { TeamContainer, TeamName } from "./styles";

export type Props = {
  teams: Team[];
  remove_team?: (team: Team) => void;
};

const TeamList: React.FC<Props> = ({ teams, remove_team }) => {
  return (
    <FlexColScroll maxHeight="400px">
      {teams.map((team, i) => {
        return <TeamWithMembers key={team.id} team={team} />;
        // return (
        //   <TeamContainer key={team.id}>
        //     <TeamName>{team.display_name}</TeamName>
        //     {remove_team && (
        //       <StandardButton
        //         label="Remove"
        //         style_kind="muted"
        //         on_click={() => remove_team(team)}
        //       />
        //     )}
        //   </TeamContainer>
        // );
      })}
    </FlexColScroll>
  );
};

export default TeamList;
