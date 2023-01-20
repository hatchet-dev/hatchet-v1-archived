import React, { useEffect } from "react";
import { useHistory, useLocation, useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { useAtom } from "jotai";
import { currOrgAtom, currTeamAtom } from "shared/atoms/atoms";
import { Spinner } from "@hatchet-dev/hatchet-components";

type Props = {
  organization?: boolean;
  team?: boolean;
  children?: React.ReactNode;
};

const Populator: React.FunctionComponent<Props> = ({
  organization,
  team,
  children,
}) => {
  const history = useHistory();
  const params: any = useParams();
  const location = useLocation();

  const [currOrg, setCurrOrg] = useAtom(currOrgAtom);
  const [currTeam, setCurrTeam] = useAtom(currTeamAtom);
  const orgEnabled = !!organization;
  const teamEnabled = !!team;

  const { data, isLoading, isFetching } = useQuery({
    queryKey: ["current_user_organizations"],
    queryFn: async () => {
      const res = await api.listUserOrganizations();
      return res;
    },
    retry: false,
    enabled: orgEnabled,
  });

  const currTeamsQuery = useQuery({
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

  const matchedOrg = data?.data?.rows?.filter(
    (org) => currOrg?.id == org.id
  )[0];

  const matchedTeam = currTeamsQuery.data?.data?.rows?.filter(
    (team) => currTeam?.id == team.id
  )[0];

  const matchedParamTeam = currTeamsQuery.data?.data?.rows?.filter(
    (team) => params?.team == team.id
  )[0];

  useEffect(() => {
    if (orgEnabled) {
      // if curr org id is not set, or it is set but is not found in the current list,
      // set it to the first item in the array, or redirect to the creation screen if no orgs
      if (!isFetching) {
        if (!currOrg || !matchedOrg) {
          if (data?.data?.rows?.length > 0) {
            setCurrOrg(data?.data?.rows[0]);
          } else {
            history.push("/organization/create");
          }
        }
      }
    }
  }, [currOrg, data, isFetching, orgEnabled]);

  useEffect(() => {
    if (teamEnabled) {
      if (matchedParamTeam && matchedParamTeam.id != currTeam.id) {
        setCurrTeam(matchedParamTeam);
      } else if (!currTeamsQuery.isFetching) {
        // if curr team id is not set, or it is set but is not found in the current list,
        // set it to the first item in the array, or redirect to the creation screen if no teams

        if (!currTeam || !matchedTeam) {
          if (currTeamsQuery.data?.data?.rows?.length > 0) {
            setCurrTeam(currTeamsQuery.data?.data?.rows[0]);
          } else {
            history.push("/teams/create");
          }
        }
      }
    }
  }, [
    currTeam,
    currTeamsQuery.data,
    currTeamsQuery.isFetching,
    teamEnabled,
    matchedParamTeam,
  ]);

  if (isLoading || !currOrg || !matchedOrg) {
    return <div>Loading</div>;
  }

  return <>{children}</>;
};

export default Populator;
