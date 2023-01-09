import { FlexRow, MaterialIcon } from "components/globals";
import Selector, { Selection } from "components/selector";
import React from "react";
import github from "assets/github.png";
import { TopBarProductName, TopBarWrapper } from "./styles";
import api from "shared/api";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useHistory } from "react-router-dom";

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

const TopBar: React.FunctionComponent = () => {
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
      history.push("/user/settings");
    }
  };

  return (
    <TopBarWrapper>
      <FlexRow>
        <TopBarProductName>Hatchet</TopBarProductName>
        <Selector
          placeholder={data?.data?.display_name}
          placeholder_material_icon="person"
          options={options}
          select={onSelect}
          option_alignment="right"
          fill_selection={false}
        />
      </FlexRow>
    </TopBarWrapper>
  );
};

export default TopBar;
