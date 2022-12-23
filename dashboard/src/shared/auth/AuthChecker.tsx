import React from "react";
import { useHistory } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import api from "shared/api";

const AuthChecker: React.FunctionComponent = ({ children }) => {
  let history = useHistory();

  const { status, data, error, isFetching } = useQuery({
    queryKey: ["current_user"],
    queryFn: async () => {
      const res = await api.getCurrentUser();
      return res;
    },
    // Requery every 10 seconds for the current user
    refetchInterval: 10000,
    retry: 2,
  });

  if (error) {
    history.push("/login");
  }

  return <>{children}</>;
};

export default AuthChecker;
