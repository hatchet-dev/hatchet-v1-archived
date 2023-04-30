import {
  FlexColScroll,
  MaterialIcon,
  StandardButton,
  Table,
} from "hatchet-components";
import React from "react";
import {
  OrganizationMemberSanitized,
  Team,
} from "shared/api/generated/data-contracts";
import { capitalize, relativeDate } from "shared/utils";
import TeamWithMembers from "../teamwithmembers";
import { TeamContainer, TeamName } from "./styles";
import { useHistory } from "react-router-dom";

export type Props = {
  teams: Team[];
  remove_team?: (team: Team) => void;
  expanded?: boolean;
};

const TeamList: React.FC<Props> = ({ teams, remove_team, expanded = true }) => {
  const history = useHistory();

  if (expanded) {
    return (
      <FlexColScroll maxHeight="400px">
        {teams.map((team, i) => {
          return <TeamWithMembers key={team.id} team={team} />;
        })}
      </FlexColScroll>
    );
  }

  const columns = [
    {
      Header: "Name",
      accessor: "name",
    },
    {
      Header: "Created",
      accessor: "created_at",
    },
  ];

  const tableData = teams.map((row) => {
    return {
      id: row.id,
      name: row.display_name,
      created_at: relativeDate(row.created_at),
    };
  });

  return (
    <Table
      columns={columns}
      data={tableData}
      dataName="teams"
      onRowClick={(row: any) => {
        history.push(`/teams/${row.original.id}/modules`);
      }}
    />
  );
};

export default TeamList;
