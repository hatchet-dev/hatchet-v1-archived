import {
  H2,
  H3,
  HorizontalSpacer,
  Placeholder,
  SmallSpan,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { CreateTeamRequest, Team } from "shared/api/generated/data-contracts";
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import CreateTeamForm from "../createteamform";
import TeamList from "../teamlist";

const defaultAddTeamHelper =
  "Add organization members by entering their email and assigning them a role.";

type Props = {
  create_team?: (team: Team) => void;
  can_remove?: boolean;
  header_level?: "H2" | "H3";
  add_team_helper?: string;
};

const TeamManager: React.FunctionComponent<Props> = ({
  can_remove = false,
  header_level = "H2",
  add_team_helper = defaultAddTeamHelper,
  create_team,
}) => {
  const [currOrg] = useAtom(currOrgAtom);
  const [err, setErr] = useState("");
  const history = useHistory();

  const { data, isLoading, refetch } = useQuery({
    queryKey: ["organization_teams", currOrg.id],
    queryFn: async () => {
      const res = await api.listTeams(currOrg.id);
      return res;
    },
    retry: false,
  });

  const mutation = useMutation({
    mutationKey: ["create_team", currOrg.id],
    mutationFn: async (team: CreateTeamRequest) => {
      const res = await api.createTeam(currOrg.id, team);
      return res;
    },
    onSuccess: (data) => {
      setErr("");
      refetch();
      create_team && create_team(data?.data);
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  if (!currOrg) {
    history.push("/");
  }

  const Header = header_level == "H2" ? H2 : H3;

  return (
    <>
      <Header>Create Team</Header>
      <HorizontalSpacer spacepixels={20} />
      <SmallSpan>{add_team_helper}</SmallSpan>
      <HorizontalSpacer spacepixels={20} />
      <CreateTeamForm
        submit={async (team, cb) => {
          try {
            await mutation.mutateAsync(team);
          } catch (e) {}

          cb();
        }}
        err={err}
      />
      <HorizontalSpacer spacepixels={24} />
      <Header>Teams</Header>
      <HorizontalSpacer spacepixels={24} />
      {isLoading && <Placeholder loading={isLoading}></Placeholder>}
      {!isLoading && (
        <TeamList
          teams={data.data?.rows}
          // remove_member={can_remove && removeMember}
        />
      )}
    </>
  );
};

export default TeamManager;
