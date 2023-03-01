import {
  FlexCol,
  FlexColCenter,
  FlexRow,
  FlexRowRight,
  H1,
  HorizontalSpacer,
  P,
  StyledDeprecatedText,
  Table,
  StandardButton,
  Spinner,
  Placeholder,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { relativeDate } from "shared/utils";
import { useAtom } from "jotai";
import { currOrgAtom } from "shared/atoms/atoms";
import { useHistory } from "react-router-dom";
import NotificationsList from "./components/notificationslist/NotificationsList";
import ExpandedNotification from "./components/expandednotification/ExpandedNotification";

const Notifications: React.FunctionComponent = () => {
  const history = useHistory();
  const [currOrg, setCurrOrg] = useAtom(currOrgAtom);
  const [err, setErr] = useState("");

  const notifsQuery = useQuery({
    queryKey: ["notifications", currOrg.id],
    queryFn: async () => {
      const res = await api.listNotifications(currOrg.id);
      return res;
    },
    retry: false,
  });

  if (notifsQuery.isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  return (
    <FlexColCenter>
      <FlexCol width="100%" height="calc(100vh - 200px)">
        <H1>Notifications</H1>
        <HorizontalSpacer spacepixels={30} />
        <FlexRow height="100%" width="100%">
          <NotificationsList notifications={notifsQuery.data?.data?.rows} />
          <ExpandedNotification />
        </FlexRow>
      </FlexCol>
    </FlexColCenter>
  );
};

export default Notifications;
