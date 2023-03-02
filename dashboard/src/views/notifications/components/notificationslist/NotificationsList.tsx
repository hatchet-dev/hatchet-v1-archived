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
  FlexColScroll,
  SmallSpan,
} from "@hatchet-dev/hatchet-components";
import React, { useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import api from "shared/api";
import { relativeDate } from "shared/utils";
import { useAtom } from "jotai";
import { currOrgAtom } from "shared/atoms/atoms";
import { useHistory } from "react-router-dom";
import { NotificationMeta } from "shared/api/generated/data-contracts";
import {
  NotificationDate,
  NotificationListContainer,
  NotificationMetaContainer,
  NotificationTitle,
} from "./styles";
import theme from "shared/theme";
import { css } from "styled-components";

type Props = {
  notifications: NotificationMeta[];
};

const NotificationsList: React.FC<Props> = ({ notifications }) => {
  const history = useHistory();

  const selectNotification = (notif: NotificationMeta) => {
    history.push(`/notifications/teams/${notif.team_id}/${notif.id}`);
  };

  return (
    <NotificationListContainer width="460px">
      {notifications.map((notif) => {
        return (
          <NotificationMetaContainer
            key={notif.id}
            onClick={() => selectNotification(notif)}
          >
            <FlexRow>
              <NotificationTitle>{notif.title}</NotificationTitle>
              <NotificationDate>
                {relativeDate(notif.updated_at)}
              </NotificationDate>
            </FlexRow>
            <HorizontalSpacer spacepixels={20} />
            <FlexRow>
              <SmallSpan>{notif.message}</SmallSpan>
            </FlexRow>
          </NotificationMetaContainer>
        );
      })}
    </NotificationListContainer>
  );
};

export default NotificationsList;
