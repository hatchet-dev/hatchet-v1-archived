import React from "react";
import api from "shared/api";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useHistory } from "react-router-dom";
import { TopBar, Selector, Selection } from "@hatchet-dev/hatchet-components";

const options = [
  {
    material_icon: "settings",
    label: "Settings",
    value: "settings",
  },
  {
    material_icon: "logout",
    label: "Logout",
    value: "logout",
  },
];

type Props = {
  children?: React.ReactNode;
};

const TopBarWithProfile: React.FunctionComponent<Props> = ({ children }) => {
  const history = useHistory();

  const { data, refetch } = useQuery({
    queryKey: ["current_user"],
    queryFn: async () => {
      const res = await api.getCurrentUser();
      return res;
    },
    retry: false,
  });

  const { mutate, isLoading } = useMutation(api.logoutUser, {
    mutationKey: ["logout"],
    onSuccess: (data) => {
      refetch();
      history.push("/login");
    },
    onError: (err: any) => {
      console.log("ERR", err.error);
    },
  });

  const onSelect = (selection: Selection) => {
    if (selection.value == "logout") {
      mutate({});
    } else if (selection.value == "settings") {
      history.push("/user/settings/profile");
    }
  };

  return (
    <TopBar>
      <Selector
        placeholder={data?.data?.display_name}
        placeholder_material_icon="person"
        options={options}
        select={onSelect}
        option_alignment="right"
        fill_selection={false}
      />
      {children}
    </TopBar>
  );
};

export default TopBarWithProfile;
