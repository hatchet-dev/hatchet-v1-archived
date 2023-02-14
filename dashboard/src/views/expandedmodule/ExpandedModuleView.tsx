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
import GithubRef from "components/githubref";
import Status from "components/status";
import { useAtom } from "jotai";
import React, { useState, useEffect } from "react";
import { useHistory, useParams } from "react-router-dom";
import api from "shared/api";
import { currTeamAtom } from "shared/atoms/atoms";
import { relativeDate } from "shared/utils";
import RunsList from "../../components/runslist";
import ModuleRunsList from "./components/modulerunslist";
import ExpandedModuleMonitors from "./components/monitors";
import RunsTab from "./components/runstab";

const TabOptions = [
  "Runs",
  "Resource Explorer",
  "Policies",
  "Monitors",
  "Settings",
];

const ExpandedModuleView: React.FunctionComponent = () => {
  const history = useHistory();
  const [selectedTab, setSelectedTab] = useState(TabOptions[0]);
  const [currTeam] = useAtom(currTeamAtom);
  const params: any = useParams();

  useEffect(() => {
    if (!params?.module) {
      history.push(`/team/${currTeam?.id}/modules`);
    }
  }, [params]);

  const { data, isLoading, refetch, isFetching } = useQuery({
    queryKey: ["module", currTeam.id, params?.module],
    queryFn: async () => {
      const res = await api.getModule(currTeam.id, params?.module);
      return res;
    },
    retry: false,
  });

  const renderTabContents = () => {
    switch (selectedTab) {
      case "Runs":
        return <RunsTab team_id={currTeam.id} module={data?.data} />;
      case "Resource Explorer":
        return <HeirarchyGraph width={100} height={100} />;
      case "Monitors":
        return <ExpandedModuleMonitors />;
      default:
        return <Span>Settings</Span>;
    }
  };

  const renderGithubRepoRef = () => {
    let depl = data?.data.deployment;
    if (depl.github_repo_name != "") {
      return (
        <GithubRef
          text={`${depl.github_repo_owner}/${depl.github_repo_name}`}
          link={`https://github.com/${depl.github_repo_owner}/${depl.github_repo_name}`}
        />
      );
    }
  };

  const renderLock = () => {
    if (data?.data.lock_id != "") {
      return (
        <Status
          status_text={`Lock: ${data?.data.lock_id}`}
          material_icon="lock"
        />
      );
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
            label: "Modules",
            link: `/team/${currTeam.id}/modules`,
          },
          {
            label: `${data.data.name}`,
            link: "",
          },
        ]}
      />
      <HorizontalSpacer spacepixels={12} />
      <FlexRow>
        <H1>{data.data.name}</H1>
        {renderGithubRepoRef()}
      </FlexRow>
      <HorizontalSpacer spacepixels={20} />
      <FlexCol>
        <P>This page contains information about the {data.data.name} module.</P>
        <HorizontalSpacer spacepixels={8} />
        <FlexRowLeft gap="16px">
          <Status
            status_text={relativeDate(data?.data.updated_at)}
            material_icon="schedule"
          />
          <Status status_text="Deployed" material_icon="check" />
          {renderLock()}
        </FlexRowLeft>
      </FlexCol>
      <HorizontalSpacer spacepixels={30} />
      <TabList tabs={TabOptions} selectTab={setSelectedTab} />
      {renderTabContents()}
    </>
  );
};

export default ExpandedModuleView;
