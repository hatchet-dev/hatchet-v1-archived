import React from "react";
import { useHistory } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import api from "shared/api";

type Props = {
  check_authenticated?: boolean;
};

const AuthChecker: React.FunctionComponent<Props> = ({
  check_authenticated = true,
  children,
}) => {
  let history = useHistory();

  const { error, isLoading, isInitialLoading } = useQuery({
    queryKey: ["current_user"],
    queryFn: async () => {
      const res = await api.getCurrentUser();
      return res;
    },
    // Requery every 10 seconds for the current user
    refetchInterval: 10000,
    retry: false,
  });

  // TODO(abelanger5): style loading
  // if (isInitialLoading) {
  //   return <div>Loading...</div>;
  // }

  if (!isLoading) {
    if (check_authenticated && error) {
      history.push("/login");
    } else if (!check_authenticated && !error) {
      history.push("/");
    }
  }

  return <>{children}</>;
};

export default AuthChecker;
