import {
  FlexCol,
  FlexColCenter,
  FlexRowRight,
  FlexRow,
  H1,
  HorizontalSpacer,
  P,
  StyledDeprecatedText,
  Table,
  StandardButton,
  Spinner,
  Placeholder,
  H4,
  FlexColScroll,
  SmallSpan,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { relativeDate } from "shared/utils";
import { useAtom } from "jotai";
import { currOrgAtom } from "shared/atoms/atoms";
import { useHistory, useParams } from "react-router-dom";
import { NotificationMeta } from "shared/api/generated/data-contracts";
import { ExpandedNotificationContainer } from "./styles";
import theme from "shared/theme";
import { css } from "styled-components";
import Status from "components/status";

type Props = {};

const ExpandedNotification: React.FC<Props> = () => {
  const history = useHistory();
  const params: any = useParams();

  const isSelected = !!(params.team && params.notification);

  const notificationQuery = useQuery({
    queryKey: ["notification", params.team, params.notification],
    queryFn: async () => {
      const res = await api.getNotification(params.team, params.notification);
      return res;
    },
    retry: false,
    refetchOnReconnect: false,
    refetchOnWindowFocus: false,
    enabled: isSelected,
  });

  const renderContents = () => {
    if (!isSelected) {
      return (
        <Placeholder>
          <SmallSpan>No notification selected</SmallSpan>
        </Placeholder>
      );
    }

    const notif = notificationQuery?.data?.data;

    if (notificationQuery.isLoading || notificationQuery.isFetching || !notif) {
      return (
        <Placeholder>
          <Spinner />
        </Placeholder>
      );
    }

    return (
      <FlexCol>
        <FlexRow>
          <H4>{notif.title}</H4>
          <FlexRowRight gap="10px">
            <Status
              status_text={relativeDate(notif.updated_at)}
              material_icon="schedule"
            />
            <Status status_text="Resolved" material_icon="check" />
          </FlexRowRight>
        </FlexRow>
        <HorizontalSpacer spacepixels={30} />
        <P>{notif.message}</P>
      </FlexCol>
    );
  };

  return (
    <ExpandedNotificationContainer>
      {renderContents()}
    </ExpandedNotificationContainer>
  );
};

export default ExpandedNotification;
