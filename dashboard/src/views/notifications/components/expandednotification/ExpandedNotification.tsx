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
  FlexRowLeft,
} from "hatchet-components";
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
import LinkButton from "components/linkbutton";
import ExpandedResultOverview from "components/monitor/expandedresultoverview";
import ExpandedRunOverview from "components/run/expandedrunoverview";

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

    const getModuleLink = () => {
      return `/teams/${params.team}/modules/${notif.module_id}`;
    };

    const renderResolvedStatus = () => {
      if (notif.resolved) {
        return <Status status_text="Resolved" material_icon="check" />;
      }

      return <Status status_text="Unresolved" material_icon="error" />;
    };

    const renderMonitorLink = () => {
      if (notif.monitor_results && notif.monitor_results.length > 0) {
        const result = notif.monitor_results[0];

        return (
          <LinkButton
            text="Monitor"
            link={`/teams/${params.team}/monitors/${result.module_monitor_id}`}
          />
        );
      }
    };

    const renderDetails = () => {
      if (notif.monitor_results && notif.monitor_results.length > 0) {
        const result = notif.monitor_results[0];

        return (
          <ExpandedResultOverview
            team_id={params.team}
            module_monitor_result={result}
          />
        );
      }

      if (notif.runs && notif.runs.length > 0) {
        const run = notif.runs[0];

        return (
          <ExpandedRunOverview
            team_id={params.team}
            module={notif.module}
            module_run={run}
          />
        );
      }
    };

    return (
      <FlexCol>
        <FlexRow>
          <H4>{notif.title}</H4>
          <FlexRowRight gap="10px">
            <Status
              status_text={relativeDate(notif.updated_at)}
              material_icon="schedule"
            />
            {renderResolvedStatus()}
          </FlexRowRight>
        </FlexRow>
        <HorizontalSpacer spacepixels={30} />
        <FlexRowLeft gap="8px">
          <LinkButton text="Module" link={getModuleLink()} />
          {renderMonitorLink()}
        </FlexRowLeft>
        <HorizontalSpacer spacepixels={30} />
        <P>{notif.message}</P>
        <HorizontalSpacer spacepixels={30} />
        {renderDetails()}
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
