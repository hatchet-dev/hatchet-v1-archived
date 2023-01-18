import React from "react";
import { useHistory } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import api from "shared/api";
import VerifyEmailPromptView from "views/verifyemailprompt/VerifyEmailPromptView";

type Props = {
  check_authenticated?: boolean;
  require_unverified_email?: boolean;
  allow_unverified_email?: boolean;
  children?: React.ReactNode;
};

const AuthChecker: React.FunctionComponent<Props> = ({
  check_authenticated = true,
  require_unverified_email = false,
  allow_unverified_email = false,
  children,
}) => {
  let history = useHistory();

  const currentUserQuery = useQuery({
    queryKey: ["current_user"],
    queryFn: async () => {
      const res = await api.getCurrentUser();
      return res;
    },
    // Requery every 10 seconds for the current user
    refetchInterval: 10000,
    retry: false,
  });

  const metadataQuery = useQuery({
    queryKey: ["api_metadata"],
    queryFn: async () => {
      const res = await api.getServerMetadata();
      return res;
    },
    retry: false,
  });

  if (!currentUserQuery.isLoading && !metadataQuery.isLoading) {
    if (check_authenticated && currentUserQuery.error) {
      history.push("/login");
    } else if (!check_authenticated && !currentUserQuery.error) {
      history.push("/");
    }

    if (
      !allow_unverified_email &&
      check_authenticated &&
      !currentUserQuery.error &&
      !currentUserQuery.data?.data?.email_verified &&
      metadataQuery.data?.data?.auth.require_email_verification
    ) {
      return <VerifyEmailPromptView />;
    }

    if (
      require_unverified_email &&
      check_authenticated &&
      !currentUserQuery.error &&
      currentUserQuery.data?.data?.email_verified &&
      metadataQuery.data?.data?.auth.require_email_verification
    ) {
      history.push("/");
    }
  }

  return <>{children}</>;
};

export default AuthChecker;
