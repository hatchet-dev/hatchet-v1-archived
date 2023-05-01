import {
  Grid,
  H1,
  H2,
  HorizontalSpacer,
  P,
  GridCard,
  Placeholder,
  Spinner,
  FlexRowLeft,
} from "hatchet-components";
import { useQuery } from "@tanstack/react-query";
import RunsList from "components/run/runslist";
import TeamList from "components/team/teamlist";
import { useAtom } from "jotai";
import React from "react";
import api from "shared/api";
import { currOrgAtom } from "shared/atoms/atoms";
import theme from "shared/theme";
import { css } from "styled-components";

const HomeView: React.FunctionComponent = () => {
  const [currOrg, setCurrOrg] = useAtom(currOrgAtom);

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

  const renderUserTeams = () => {
    if (userTeamsQuery.isLoading) {
      return (
        <Placeholder>
          <Spinner />
        </Placeholder>
      );
    }

    return (
      <FlexRowLeft width="400px">
        <TeamList teams={userTeamsQuery.data?.data?.rows} expanded={false} />
      </FlexRowLeft>
    );
  };

  return (
    <>
      <H1>Home</H1>
      <HorizontalSpacer
        spacepixels={16}
        overrides={css({
          borderBottom: theme.line.thick,
        }).toString()}
      />
      <HorizontalSpacer spacepixels={16} />
      <H2>Your Teams</H2>
      <HorizontalSpacer spacepixels={16} />
      <P>
        Here are all your joined teams in the {currOrg.display_name}{" "}
        organization.
      </P>
      <HorizontalSpacer spacepixels={32} />
      {renderUserTeams()}
    </>
  );
};

export default HomeView;
