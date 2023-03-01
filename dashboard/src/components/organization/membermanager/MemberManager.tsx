import {
  H2,
  H3,
  HorizontalSpacer,
  SmallSpan,
  Placeholder,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import {
  CreateOrgMemberInviteRequest,
  OrganizationMemberSanitized,
} from "shared/api/generated/data-contracts";
import InviteMemberForm from "components/organization/invitememberform";
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import MemberList from "components/organization/memberlist";

const defaultAddMemberHelper =
  "Add organization members by entering their email and assigning them a role.";

type Props = {
  can_remove?: boolean;
  header_level?: "H2" | "H3";
  add_member_helper?: string;
};

const MemberManager: React.FunctionComponent<Props> = ({
  can_remove = false,
  header_level = "H2",
  add_member_helper = defaultAddMemberHelper,
}) => {
  const [currOrg] = useAtom(currOrgAtom);
  const [err, setErr] = useState("");
  const history = useHistory();

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
    mutationFn: async (invite: CreateOrgMemberInviteRequest) => {
      const res = await api.createOrgMemberInvite(currOrg.id, invite);
      return res;
    },
    onSuccess: (data) => {
      setErr("");
      refetch();
    },
    onError: (err: any) => {
      if (!err?.error?.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
        return;
      }

      setErr(err.error.errors[0].description);
    },
  });

  const removeMemberMutation = useMutation({
    mutationKey: ["create_organization_invite", currOrg.id],
    mutationFn: async (orgMemberId: string) => {
      const res = await api.deleteOrgMember(currOrg.id, orgMemberId);
      return res;
    },
    onSuccess: (data) => {
      setErr("");
      refetch();
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

  const removeMember = (member: OrganizationMemberSanitized) => {
    removeMemberMutation.mutate(member.id);
  };

  const Header = header_level == "H2" ? H2 : H3;

  return (
    <>
      <Header>Current Members</Header>
      <HorizontalSpacer spacepixels={24} />
      {isLoading && <Placeholder loading={isLoading}></Placeholder>}
      {!isLoading && (
        <MemberList
          members={data.data?.rows}
          remove_member={can_remove && removeMember}
        />
      )}
      <HorizontalSpacer spacepixels={24} />
      <Header>Add Members</Header>
      <HorizontalSpacer spacepixels={20} />
      <SmallSpan>{add_member_helper}</SmallSpan>
      <HorizontalSpacer spacepixels={20} />
      <InviteMemberForm
        submit={async (invite, cb) => {
          try {
            await mutation.mutateAsync(invite);
          } catch (e) {}

          cb();
        }}
        err={err}
      />
    </>
  );
};

export default MemberManager;
