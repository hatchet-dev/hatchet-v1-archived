import {
  H1,
  HorizontalSpacer,
  P,
  Span,
  Breadcrumbs,
  HeirarchyGraph,
  TabList,
  StandardButton,
  FlexRowRight,
  FlexCol,
  Placeholder,
  Spinner,
  FlexRowLeft,
  FlexRow,
} from "@hatchet-dev/hatchet-components";
import { useMutation, useQuery } from "@tanstack/react-query";
import CodeBlock from "components/codeblock";
import GithubRef from "components/githubref";
import ResultsList from "components/monitor/resultslist";
import Status from "components/status";
import { useAtom } from "jotai";
import React, { useState, useEffect } from "react";
import { useHistory, useParams } from "react-router-dom";
import api from "shared/api";
import { currTeamAtom } from "shared/atoms/atoms";
import { relativeDate } from "shared/utils";

const TabOptions = ["Recent Results", "Policy", "Settings"];

const ExpandedMonitorView: React.FunctionComponent = () => {
  const history = useHistory();
  const [selectedTab, setSelectedTab] = useState(TabOptions[0]);
  const [currTeam] = useAtom(currTeamAtom);
  const params: any = useParams();

  useEffect(() => {
    if (!params?.monitor) {
      history.push(`/team/${currTeam?.id}/monitors`);
    }
  }, [params]);

  const { data, isLoading, refetch, isFetching } = useQuery({
    queryKey: ["monitor", currTeam.id, params?.monitor],
    queryFn: async () => {
      const res = await api.getMonitor(currTeam.id, params?.monitor);
      return res;
    },
    retry: false,
  });

  const renderTabContents = () => {
    switch (selectedTab) {
      case "Recent Results":
        return (
          <ResultsList
            team_id={currTeam?.id}
            module_monitor_id={data?.data.id}
          />
        );
      case "Policy":
        return (
          <CodeBlock
            value={data?.data?.policy_bytes}
            height="calc(100% - 250px)"
            readOnly={true}
          />
        );
      case "Settings":
        return <div>Policy goes here</div>;
    }
  };

  if (isLoading) {
    return (
      <Placeholder>
        <Spinner />
      </Placeholder>
    );
  }

  return (
    <>
      <Breadcrumbs
        breadcrumbs={[
          {
            label: "Monitors",
            link: `/team/${currTeam.id}/monitors`,
          },
          {
            label: data?.data.name,
            link: "",
          },
        ]}
      />
      <HorizontalSpacer spacepixels={12} />
      <FlexRow>
        <H1>{data?.data.name}</H1>
      </FlexRow>
      <HorizontalSpacer spacepixels={20} />
      <FlexCol>
        <P>
          This page contains information about the {data?.data.name} monitor.
        </P>
        <HorizontalSpacer spacepixels={8} />
        <FlexRowLeft gap="16px">
          {/* <Status
              status_text={relativeDate(data?.data.updated_at)}
              material_icon="schedule"
            />
            <Status status_text="Deployed" material_icon="check" />
            {renderLock()} */}
        </FlexRowLeft>
      </FlexCol>
      <HorizontalSpacer spacepixels={30} />
      <TabList tabs={TabOptions} selectTab={setSelectedTab} />
      {renderTabContents()}
    </>
  );
};

export default ExpandedMonitorView;
