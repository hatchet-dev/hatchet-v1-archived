import {
  FlexCol,
  FlexColCenter,
  FlexRowRight,
  H1,
  HorizontalSpacer,
  P,
  StyledDeprecatedText,
  Table,
  StandardButton,
  Spinner,
  Placeholder,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { relativeDate } from "shared/utils";
import { useAtom } from "jotai";
import { currOrgAtom } from "shared/atoms/atoms";
import { useHistory } from "react-router-dom";

const TeamsList: React.FunctionComponent = () => {
  const history = useHistory();
  const [currOrg, setCurrOrg] = useAtom(currOrgAtom);
  const [err, setErr] = useState("");

  const userTeamsQuery = useQuery({
    queryKey: ["current_user_teams", currOrg.id],
    queryFn: async () => {
      const res = await api.listUserTeams({
        organization_id: currOrg.id,
      });
      return res;
    },
    retry: false,
    enabled: !!currOrg,
  });

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

  const tableData = userTeamsQuery.data?.data?.rows.map((row) => {
    return {
      id: row.id,
      name: row.display_name,
      created_at: relativeDate(row.created_at),
    };
  });

  if (userTeamsQuery.isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  const renderTeamData = () => {
    return (
      <>
        <FlexRowRight>
          <StandardButton
            label="Create team"
            material_icon="add"
            on_click={() => {
              // setCreatePAT(true);
              history.push("/organization/teams/create");
            }}
          />
        </FlexRowRight>
        <HorizontalSpacer spacepixels={20} />
        <Table columns={columns} data={tableData} dataName="teams" />
        <HorizontalSpacer spacepixels={20} />
      </>
    );
  };

  return (
    <FlexColCenter>
      <FlexCol width="100%" maxWidth="640px">
        <H1>Teams</H1>
        <HorizontalSpacer spacepixels={12} />
        {renderTeamData()}
      </FlexCol>
    </FlexColCenter>
  );
};

export default TeamsList;
