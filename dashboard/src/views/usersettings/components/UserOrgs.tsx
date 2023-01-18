import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { SectionArea } from "@hatchet-dev/hatchet-components";
import OrgList from "components/organization/orglist";

const UserOrgs: React.FunctionComponent = () => {
  const [err, setErr] = useState("");

  const orgQuery = useQuery({
    queryKey: ["current_user_organizations"],
    queryFn: async () => {
      const res = await api.listUserOrganizations();
      return res;
    },
    retry: false,
  });

  const leaveOrgMutation = useMutation({
    mutationKey: ["leave_organization"],
    mutationFn: (orgId: string) => {
      return api.leaveOrg(orgId);
    },
    onSuccess: (data) => {
      orgQuery.refetch();
    },
    onError: (err: any) => {
      if (!err.error || !err.error.errors || err.error.errors.length == 0) {
        setErr("An unexpected error occurred. Please try again.");
      }

      setErr(err.error.errors[0].description);
    },
  });

  return (
    <SectionArea width="600px" loading={orgQuery.isLoading}>
      <OrgList
        orgs={orgQuery?.data?.data.rows}
        leave_org={(org) => leaveOrgMutation.mutate(org.id)}
        err={err}
      />
    </SectionArea>
  );
};

export default UserOrgs;
