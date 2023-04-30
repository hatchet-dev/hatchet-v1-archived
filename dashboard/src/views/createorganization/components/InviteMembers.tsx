import {
  FlexCol,
  FlexRowRight,
  HorizontalSpacer,
  StandardButton,
  SectionArea,
} from "hatchet-components";
import React, { useState } from "react";
import { useHistory } from "react-router-dom";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { CreateOrgMemberInviteRequest } from "shared/api/generated/data-contracts";
import { currOrgAtom } from "shared/atoms/atoms";
import { useAtom } from "jotai";
import MemberManager from "components/organization/membermanager/MemberManager";

const inviteMemberHelper =
  "Add organization members by entering their email and assigning them a role. You can also add members later from organization settings.";

const InviteMembers: React.FunctionComponent = () => {
  const [currOrg] = useAtom(currOrgAtom);
  const [err, setErr] = useState("");
  const history = useHistory();

  const { data, isLoading, refetch } = useQuery({
    queryKey: ["current_organization_members"],
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

  if (!currOrg) {
    history.push("/");
  }

  return (
    <FlexCol>
      <SectionArea width="600px">
        <MemberManager add_member_helper={inviteMemberHelper} />
      </SectionArea>
      <HorizontalSpacer spacepixels={24} />
      <FlexRowRight>
        <StandardButton
          label="Next"
          material_icon="chevron_right"
          icon_side="right"
          on_click={() => {
            history.push("/organization/create/create_teams");
          }}
          margin={"0"}
          is_loading={mutation.isLoading}
        />
      </FlexRowRight>
    </FlexCol>
  );
};

export default InviteMembers;
