import React, { useEffect } from "react";
import { useHistory } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { useAtom } from "jotai";
import { currOrgAtom } from "shared/atoms/atoms";
import Spinner from "components/loaders";

type Props = {
  organization?: boolean;
};

const Populator: React.FunctionComponent<Props> = ({
  organization,
  children,
}) => {
  const history = useHistory();

  const [currOrgId, setCurrOrgId] = useAtom(currOrgAtom);
  const orgEnabled = !!organization;

  const { data, isLoading, isFetching } = useQuery({
    queryKey: ["current_user_organizations"],
    queryFn: async () => {
      const res = await api.listUserOrganizations();
      return res;
    },
    retry: false,
    enabled: orgEnabled,
  });

  const matchedOrg = data?.data?.rows?.filter((org) => currOrgId == org.id)[0];

  useEffect(() => {
    if (orgEnabled) {
      // if curr org id is not set, or it is set but is not found in the current list,
      // set it to the first item in the array, or redirect to the creation screen if no orgs
      if (!isFetching) {
        if (currOrgId == "" || !matchedOrg) {
          if (data?.data?.rows?.length > 0) {
            setCurrOrgId(data?.data?.rows[0]?.id);
          } else {
            history.push("/organizations/create");
          }
        }
      }
    }
  }, [currOrgId, data, isFetching, orgEnabled]);

  if (isLoading || currOrgId == "" || !matchedOrg) {
    return <Spinner />;
  }

  return <>{children}</>;
};

export default Populator;
