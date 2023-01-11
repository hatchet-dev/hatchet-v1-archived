import {
  FlexCol,
  FlexRowRight,
  H2,
  HorizontalSpacer,
} from "components/globals";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation, useQuery } from "@tanstack/react-query";
import StandardButton from "components/buttons";
import SectionArea from "components/sectionarea";
import api from "shared/api";
import { CreateOrgMemberInviteRequest } from "shared/api/generated/data-contracts";
import InviteMemberForm from "components/invitememberform";
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import MemberList from "components/memberlist";
import Spinner from "components/loaders";

const InviteMembers: React.FunctionComponent = () => {
  const [currOrgId] = useAtom(currOrgAtom);
  const [err, setErr] = useState("");
  const history = useHistory();

  const { data, isLoading, refetch } = useQuery({
    queryKey: ["current_organization_members"],
    queryFn: async () => {
      const res = await api.listOrgMembers(currOrgId);
      return res;
    },
    retry: false,
  });

  const mutation = useMutation({
    mutationFn: (invite: CreateOrgMemberInviteRequest) => {
      return api.createOrgMemberInvite(currOrgId, invite);
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

  if (!currOrgId) {
    history.push("/");
  }

  return (
    <FlexCol>
      <SectionArea width={600}>
        <H2>Current Members</H2>
        <HorizontalSpacer spacepixels={24} />
        {isLoading && <Spinner />}
        {!isLoading && <MemberList members={data.data?.rows} />}
        <HorizontalSpacer spacepixels={24} />
        <H2>Add Members</H2>
        <HorizontalSpacer spacepixels={24} />
        <InviteMemberForm
          submit={async (invite, cb) => {
            try {
              await mutation.mutateAsync(invite);
            } catch (e) {}

            cb();
          }}
          err={err}
        />
      </SectionArea>
      <HorizontalSpacer spacepixels={24} />
      <FlexRowRight>
        <StandardButton
          label="Next"
          material_icon="chevron_right"
          icon_side="right"
          on_click={() => {
            history.push("/");
          }}
          margin={"0"}
          is_loading={mutation.isLoading}
        />
      </FlexRowRight>
    </FlexCol>
  );
};

export default InviteMembers;
