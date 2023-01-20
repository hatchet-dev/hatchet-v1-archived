import {
  FlexColScroll,
  FlexRow,
  H3,
  HorizontalSpacer,
  MaterialIcon,
  P,
  Placeholder,
  Spinner,
  StandardButton,
} from "@hatchet-dev/hatchet-components";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useAtom } from "jotai";
import React, { useState } from "react";
import api from "shared/api";
import {
  AddTeamMemberRequest,
  OrganizationMemberSanitized,
  Team,
} from "shared/api/generated/data-contracts";
import { currOrgAtom } from "shared/atoms/atoms";
import theme from "shared/theme";
import { capitalize } from "shared/utils";
import { css } from "styled-components";
import AddTeamMemberForm from "../addteammemberform";
import TeamMemberList from "../teammemberlist";
import { TeamHeader, ExpandIcon, TeamNameWithIcon } from "./styles";

export type Props = {
  team: Team;
  remove_team?: (team: Team) => void;
  collapsible?: boolean;
};

export const lookupOrgMember = (
  org_members: OrganizationMemberSanitized[],
  org_member_id: string
) => {
  return (
    org_members.filter((org_member) => org_member.id == org_member_id)[0] ||
    null
  );
};

const TeamWithMembers: React.FC<Props> = ({
  team,
  remove_team,
  collapsible = true,
}) => {
  const [currOrg] = useAtom(currOrgAtom);

  const [isExpanded, setIsExpanded] = useState(!collapsible);
  const [addMember, setAddMember] = useState(false);
  const [err, setErr] = useState("");
  const [addMemberErr, setAddMemberErr] = useState("");

  const { data, isLoading, refetch } = useQuery({
    queryKey: ["current_team_members", team.id],
    queryFn: async () => {
      const res = await api.listTeamMembers(team.id);
      return res;
    },
    retry: false,
  });

  const currentOrgMembersQuery = useQuery({
    queryKey: ["current_organization_members"],
    queryFn: async () => {
      const res = await api.listOrgMembers(currOrg.id);
      return res;
    },
    retry: false,
  });

  const addTeamMemberMutation = useMutation({
    mutationKey: ["add_team_member", team.id],
    mutationFn: (member: AddTeamMemberRequest) => {
      return api.addTeamMember(team.id, member);
    },
    onSuccess: (data) => {
      setAddMemberErr("");
      refetch();
    },
    onError: (err: any) => {
      if (!err.error.errors || err.error.errors.length == 0) {
        setAddMemberErr("An unexpected error occurred. Please try again.");
      }

      setAddMemberErr(err.error.errors[0].description);
    },
  });

  const deleteTeamMemberMutation = useMutation({
    mutationKey: ["delete_team_member", team.id],
    mutationFn: (member_id: string) => {
      return api.deleteTeamMember(team.id, member_id);
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

  const renderTeamMembers = () => {
    if (isLoading || currentOrgMembersQuery.isLoading) {
      return (
        <Placeholder>
          <Spinner />
        </Placeholder>
      );
    }

    return (
      <>
        <TeamMemberList
          members={data.data.rows}
          org_members={currentOrgMembersQuery.data.data.rows}
          remove_member={(member) => deleteTeamMemberMutation.mutate(member.id)}
          err={err}
        />
        {addMember && (
          <AddTeamMemberForm
            submit={(member, cb) => {
              addTeamMemberMutation.mutate(member);

              cb && cb();
            }}
            current_team_members={data.data.rows}
            org_members={currentOrgMembersQuery.data.data.rows}
            err={addMemberErr}
          />
        )}
      </>
    );
  };

  return (
    <>
      <TeamHeader>
        <TeamNameWithIcon>
          <MaterialIcon className="material-icons">group</MaterialIcon>
          <H3>{team.display_name}</H3>
        </TeamNameWithIcon>
        <TeamNameWithIcon>
          <StandardButton
            size="small"
            label="Add Member"
            on_click={() => {
              setIsExpanded(true);
              setAddMember(true);
            }}
          />
          {collapsible && (
            <ExpandIcon
              className="material-icons"
              onClick={() => setIsExpanded(!isExpanded)}
            >
              {isExpanded ? "expand_more" : "expand_less"}
            </ExpandIcon>
          )}
        </TeamNameWithIcon>
      </TeamHeader>
      <HorizontalSpacer
        spacepixels={8}
        overrides={css({
          borderTop: theme.line.thick,
        }).toString()}
      />
      {isExpanded && renderTeamMembers()}
      <HorizontalSpacer spacepixels={20} />
    </>
  );
};

export default TeamWithMembers;
